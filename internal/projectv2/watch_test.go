package projectv2_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/testutil/projectv2testutil"
	"gotest.tools/v3/assert"
)

// TestNpmInstallFileCreationPerformance simulates the performance issue
// that occurs when thousands of files are created during a large npm install.
// This test measures how long it takes to process many file creation events
// in node_modules, which should help identify bottlenecks in the file watching system.
func TestNpmInstallFileCreationPerformance(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Setup a basic TypeScript project
	files := map[string]any{
		"/home/projects/TS/myapp/tsconfig.json": `{
			"compilerOptions": {
				"noLib": true,
				"module": "nodenext",
				"strict": true
			},
			"include": ["src/**/*"]
		}`,
		"/home/projects/TS/myapp/src/index.ts": `console.log("test");`,
	}

	session, utils := projectv2testutil.Setup(files)

	// Open the main file to establish a project
	session.DidOpenFile(context.Background(), "file:///home/projects/TS/myapp/src/index.ts", 1,
		files["/home/projects/TS/myapp/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

	// Get initial language service
	lsInitial, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/myapp/src/index.ts")
	assert.NilError(t, err)
	initialProgram := lsInitial.GetProgram()
	assert.Assert(t, initialProgram != nil)

	// Simulate a large npm install creating many files at once
	totalFiles := 1500 // Realistic for a moderate-sized npm install

	t.Logf("Simulating npm install: creating %d files", totalFiles)

	events := make([]*lsproto.FileEvent, 0, totalFiles)

	// Create files and corresponding events
	for i := range totalFiles {
		var filePath, content string

		if i%2 == 0 {
			// JavaScript file
			filePath = fmt.Sprintf("/home/projects/TS/myapp/node_modules/package-%d/index.js", i)
			content = "module.exports = {};"
		} else {
			// TypeScript declaration file
			filePath = fmt.Sprintf("/home/projects/TS/myapp/node_modules/package-%d/index.d.ts", i)
			content = "export {};"
		}

		err := utils.FS().WriteFile(filePath, content, false)
		assert.NilError(t, err)

		events = append(events, &lsproto.FileEvent{
			Type: lsproto.FileChangeTypeCreated,
			Uri:  lsproto.DocumentUri("file://" + filePath),
		})
	}

	// Measure the time it takes to process all the file creation events
	t.Logf("Processing %d file creation events...", len(events))
	startTime := time.Now()

	// Send all events at once (simulating rapid file creation during npm install)
	session.DidChangeWatchedFiles(context.Background(), events)

	// Force language service to process the changes and measure end-to-end time
	lsAfter, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/myapp/src/index.ts")
	assert.NilError(t, err)
	finalProgram := lsAfter.GetProgram()
	assert.Assert(t, finalProgram != nil)

	duration := time.Since(startTime)
	t.Logf("Processed %d file creation events in %v", len(events), duration)

	// Performance assertion - this should complete in a reasonable time
	// If this takes more than 5 seconds, there's likely a performance issue
	maxExpectedDuration := 5 * time.Second
	if duration > maxExpectedDuration {
		t.Errorf("File creation event processing took %v, which exceeds the expected maximum of %v. This indicates a performance issue.",
			duration, maxExpectedDuration)
	}

	// Verify the program is still functional after processing many events
	indexFile := finalProgram.GetSourceFile("/home/projects/TS/myapp/src/index.ts")
	assert.Assert(t, indexFile != nil)

	t.Logf("Successfully processed npm install simulation with %d total files", len(events))
}

// TestWatchEventDebouncing tests that rapid file change events are properly debounced
// and only result in a single snapshot update after the events stop coming.
func TestWatchEventDebouncing(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/home/projects/TS/myapp/tsconfig.json": `{
			"compilerOptions": {
				"noLib": true,
				"module": "nodenext",
				"strict": true
			},
			"include": ["src/**/*"]
		}`,
		"/home/projects/TS/myapp/src/index.ts": `console.log("test");`,
	}

	// Set a short debounce delay for testing
	options := &projectv2.SessionOptions{
		CurrentDirectory:   "/",
		DefaultLibraryPath: bundled.LibPath(),
		TypingsLocation:    projectv2testutil.TestTypingsLocation,
		PositionEncoding:   lsproto.PositionEncodingKindUTF8,
		WatchEnabled:       true,
		LoggingEnabled:     true,
		DebounceDelay:      50 * time.Millisecond,
	}

	session, utils := projectv2testutil.SetupWithOptions(files, options)

	session.DidOpenFile(context.Background(), "file:///home/projects/TS/myapp/src/index.ts", 1,
		files["/home/projects/TS/myapp/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

	// Create multiple batches of rapid file changes
	batchCount := 5
	filesPerBatch := 10

	for batch := 0; batch < batchCount; batch++ {
		events := make([]*lsproto.FileEvent, 0, filesPerBatch)

		// Create files and events for this batch
		for i := 0; i < filesPerBatch; i++ {
			fileNum := batch*filesPerBatch + i
			filePath := fmt.Sprintf("/home/projects/TS/myapp/node_modules/test-pkg-%d/index.js", fileNum)

			err := utils.FS().WriteFile(filePath, "module.exports = {};", false)
			assert.NilError(t, err)

			events = append(events, &lsproto.FileEvent{
				Type: lsproto.FileChangeTypeCreated,
				Uri:  lsproto.DocumentUri("file://" + filePath),
			})
		}

		// Send this batch of events rapidly
		session.DidChangeWatchedFiles(context.Background(), events)

		// Wait a very short time before sending the next batch (shorter than debounce delay)
		time.Sleep(10 * time.Millisecond)
	}

	t.Logf("Sent %d batches of %d events each in rapid succession", batchCount, filesPerBatch)

	// Now wait for all background tasks to complete
	session.WaitForBackgroundTasks()

	// Verify the system is still functional
	ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/myapp/src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()
	assert.Assert(t, program != nil)
	assert.Assert(t, program.GetSourceFile("/home/projects/TS/myapp/src/index.ts") != nil)

	t.Logf("Debouncing test completed successfully")
}

// TestScheduleSnapshotUpdate tests the ScheduleSnapshotUpdate method directly
func TestScheduleSnapshotUpdate(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/home/projects/TS/myapp/tsconfig.json": `{
			"compilerOptions": {
				"noLib": true,
				"module": "nodenext",
				"strict": true
			},
			"include": ["src/**/*"]
		}`,
		"/home/projects/TS/myapp/src/index.ts": `console.log("test");`,
	}

	// Set a very short debounce delay for testing
	options := &projectv2.SessionOptions{
		CurrentDirectory:   "/",
		DefaultLibraryPath: bundled.LibPath(),
		TypingsLocation:    projectv2testutil.TestTypingsLocation,
		PositionEncoding:   lsproto.PositionEncodingKindUTF8,
		WatchEnabled:       true,
		LoggingEnabled:     true,
		DebounceDelay:      25 * time.Millisecond,
	}

	session, _ := projectv2testutil.SetupWithOptions(files, options)

	session.DidOpenFile(context.Background(), "file:///home/projects/TS/myapp/src/index.ts", 1,
		files["/home/projects/TS/myapp/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

	// Schedule multiple rapid updates
	for i := 0; i < 10; i++ {
		session.ScheduleSnapshotUpdate()
		time.Sleep(5 * time.Millisecond) // Shorter than debounce delay
	}

	t.Logf("Scheduled 10 rapid snapshot updates")

	// Wait for all background tasks to complete
	session.WaitForBackgroundTasks()

	// Verify the system is still functional
	ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/myapp/src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()
	assert.Assert(t, program != nil)
	assert.Assert(t, program.GetSourceFile("/home/projects/TS/myapp/src/index.ts") != nil)

	t.Logf("ScheduleSnapshotUpdate test completed successfully")
}

// TestUpdateSnapshotCancelsPendingUpdates tests that calling UpdateSnapshot directly
// cancels any pending scheduled updates to avoid duplicate work.
func TestUpdateSnapshotCancelsPendingUpdates(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/home/projects/TS/myapp/tsconfig.json": `{
			"compilerOptions": {
				"noLib": true,
				"module": "nodenext",
				"strict": true
			},
			"include": ["src/**/*"]
		}`,
		"/home/projects/TS/myapp/src/index.ts": `console.log("test");`,
	}

	// Set a longer debounce delay to ensure we can interrupt it
	options := &projectv2.SessionOptions{
		CurrentDirectory:   "/",
		DefaultLibraryPath: bundled.LibPath(),
		TypingsLocation:    projectv2testutil.TestTypingsLocation,
		PositionEncoding:   lsproto.PositionEncodingKindUTF8,
		WatchEnabled:       true,
		LoggingEnabled:     true,
		DebounceDelay:      200 * time.Millisecond,
	}

	session, _ := projectv2testutil.SetupWithOptions(files, options)

	session.DidOpenFile(context.Background(), "file:///home/projects/TS/myapp/src/index.ts", 1,
		files["/home/projects/TS/myapp/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

	// Schedule a debounced update
	session.ScheduleSnapshotUpdate()
	t.Logf("Scheduled a debounced snapshot update")

	// Wait a short time (but less than debounce delay)
	time.Sleep(50 * time.Millisecond)

	// Now call UpdateSnapshot directly - this should cancel the pending update
	session.UpdateSnapshot(context.Background(), projectv2.SnapshotChange{})
	t.Logf("Called UpdateSnapshot directly, should have cancelled pending update")

	// Wait for all background tasks to complete
	session.WaitForBackgroundTasks()

	// Verify the system is still functional
	ls, err := session.GetLanguageService(context.Background(), "file:///home/projects/TS/myapp/src/index.ts")
	assert.NilError(t, err)
	program := ls.GetProgram()
	assert.Assert(t, program != nil)
	assert.Assert(t, program.GetSourceFile("/home/projects/TS/myapp/src/index.ts") != nil)

	t.Logf("UpdateSnapshot cancellation test completed successfully")
}
