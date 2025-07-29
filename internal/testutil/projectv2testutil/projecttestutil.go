package projectv2testutil

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

//go:generate go tool github.com/matryer/moq -stub -fmt goimports -pkg projectv2testutil -out clientmock_generated.go ../../projectv2 Client
//go:generate go tool mvdan.cc/gofumpt -lang=go1.24 -w clientmock_generated.go

const (
	TestTypingsLocation = "/home/src/Library/Caches/typescript"
)

type TestTypingsInstallerOptions struct {
	TypesRegistry []string
	PackageToFile map[string]string
}

type SessionUtils struct {
	fs            vfs.FS
	client        *ClientMock
	testOptions   *TestTypingsInstallerOptions
	preNpmInstall func(cwd string, npmInstallArgs []string)
}

type NpmInstallRequest struct {
	Cwd            string
	NpmInstallArgs []string
}

func (h *SessionUtils) Client() *ClientMock {
	return h.client
}

func (h *SessionUtils) ExpectWatchFilesCalls(count int) func(t *testing.T) {
	var actualCalls atomic.Int32
	var wg sync.WaitGroup
	wg.Add(count)
	saveFunc := h.client.WatchFilesFunc
	h.client.WatchFilesFunc = func(_ context.Context, id projectv2.WatcherID, _ []*lsproto.FileSystemWatcher) error {
		actualCalls.Add(1)
		wg.Done()
		return nil
	}
	return func(t *testing.T) {
		t.Helper()
		wg.Wait()
		assert.Equal(t, actualCalls.Load(), int32(count))
		h.client.WatchFilesFunc = saveFunc
	}
}

func (h *SessionUtils) ExpectUnwatchFilesCalls(count int) func(t *testing.T) {
	var actualCalls atomic.Int32
	var wg sync.WaitGroup
	wg.Add(count)
	saveFunc := h.client.UnwatchFilesFunc
	h.client.UnwatchFilesFunc = func(_ context.Context, id projectv2.WatcherID) error {
		actualCalls.Add(1)
		wg.Done()
		return nil
	}
	return func(t *testing.T) {
		t.Helper()
		wg.Wait()
		assert.Equal(t, actualCalls.Load(), int32(count))
		h.client.UnwatchFilesFunc = saveFunc
	}
}

func (h *SessionUtils) ExpectNpmInstallCalls(count int) func() []NpmInstallRequest {
	var calls []NpmInstallRequest
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(count)

	if h.preNpmInstall != nil {
		panic("cannot call ExpectNpmInstallCalls without invoking the return of the previous call")
	}

	h.preNpmInstall = func(cwd string, npmInstallArgs []string) {
		mu.Lock()
		defer mu.Unlock()
		calls = append(calls, NpmInstallRequest{Cwd: cwd, NpmInstallArgs: npmInstallArgs})
		wg.Done()
	}
	return func() []NpmInstallRequest {
		wg.Wait()
		mu.Lock()
		defer mu.Unlock()
		h.preNpmInstall = nil
		return calls
	}
}

func (h *SessionUtils) NpmInstall(cwd string, npmInstallArgs []string) ([]byte, error) {
	if h.testOptions == nil {
		return nil, nil
	}

	if h.preNpmInstall == nil {
		panic(fmt.Sprintf("unexpected npm install command invoked: %v", npmInstallArgs))
	}

	// Always call preNpmInstall to decrement the wait group
	h.preNpmInstall(cwd, npmInstallArgs)

	lenNpmInstallArgs := len(npmInstallArgs)
	if lenNpmInstallArgs < 3 {
		return nil, fmt.Errorf("unexpected npm install: %s %v", cwd, npmInstallArgs)
	}

	if lenNpmInstallArgs == 3 && npmInstallArgs[2] == "types-registry@latest" {
		// Write typings file
		err := h.fs.WriteFile(tspath.CombinePaths(cwd, "node_modules/types-registry/index.json"), h.createTypesRegistryFileContent(), false)
		return nil, err
	}

	for _, atTypesPackageTs := range npmInstallArgs[2 : lenNpmInstallArgs-2] {
		// @types/packageName@TsVersionToUse
		packageName := atTypesPackageTs[7 : len(atTypesPackageTs)-len(project.TsVersionToUse)-1]
		content, ok := h.testOptions.PackageToFile[packageName]
		if !ok {
			return nil, fmt.Errorf("content not provided for %s", packageName)
		}
		err := h.fs.WriteFile(tspath.CombinePaths(cwd, "node_modules/@types/"+packageName+"/index.d.ts"), content, false)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (h *SessionUtils) FS() vfs.FS {
	return h.fs
}

var (
	typesRegistryConfigTextOnce sync.Once
	typesRegistryConfigText     string
)

func TypesRegistryConfigText() string {
	typesRegistryConfigTextOnce.Do(func() {
		var result strings.Builder
		for key, value := range TypesRegistryConfig() {
			if result.Len() != 0 {
				result.WriteString(",")
			}
			result.WriteString(fmt.Sprintf("\n      \"%s\": \"%s\"", key, value))

		}
		typesRegistryConfigText = result.String()
	})
	return typesRegistryConfigText
}

var (
	typesRegistryConfigOnce sync.Once
	typesRegistryConfig     map[string]string
)

func TypesRegistryConfig() map[string]string {
	typesRegistryConfigOnce.Do(func() {
		typesRegistryConfig = map[string]string{
			"latest": "1.3.0",
			"ts2.0":  "1.0.0",
			"ts2.1":  "1.0.0",
			"ts2.2":  "1.2.0",
			"ts2.3":  "1.3.0",
			"ts2.4":  "1.3.0",
			"ts2.5":  "1.3.0",
			"ts2.6":  "1.3.0",
			"ts2.7":  "1.3.0",
		}
	})
	return typesRegistryConfig
}

func (h *SessionUtils) createTypesRegistryFileContent() string {
	var builder strings.Builder
	builder.WriteString("{\n  \"entries\": {")
	for index, entry := range h.testOptions.TypesRegistry {
		h.appendTypesRegistryConfig(&builder, index, entry)
	}
	index := len(h.testOptions.TypesRegistry)
	for key := range h.testOptions.PackageToFile {
		if !slices.Contains(h.testOptions.TypesRegistry, key) {
			h.appendTypesRegistryConfig(&builder, index, key)
			index++
		}
	}
	builder.WriteString("\n  }\n}")
	return builder.String()
}

func (h *SessionUtils) appendTypesRegistryConfig(builder *strings.Builder, index int, entry string) {
	if index > 0 {
		builder.WriteString(",")
	}
	builder.WriteString(fmt.Sprintf("\n    \"%s\": {%s\n    }", entry, TypesRegistryConfigText()))
}

func Setup(files map[string]any) (*projectv2.Session, *SessionUtils) {
	return SetupWithTypingsInstaller(files, nil)
}

func SetupWithTypingsInstaller(files map[string]any, testOptions *TestTypingsInstallerOptions) (*projectv2.Session, *SessionUtils) {
	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	clientMock := &ClientMock{}
	sessionUtils := &SessionUtils{
		fs:          fs,
		client:      clientMock,
		testOptions: testOptions,
	}

	session := projectv2.NewSession(&projectv2.SessionInit{
		Options: &projectv2.SessionOptions{
			CurrentDirectory:   "/",
			DefaultLibraryPath: bundled.LibPath(),
			TypingsLocation:    TestTypingsLocation,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       true,
			LoggingEnabled:     false,
		},
		FS:         fs,
		Client:     clientMock,
		NpmInstall: sessionUtils.NpmInstall,
	})

	return session, sessionUtils
}
