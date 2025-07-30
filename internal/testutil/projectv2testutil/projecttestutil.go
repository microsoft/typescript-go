package projectv2testutil

import (
	"bufio"
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

//go:generate go tool github.com/matryer/moq -stub -fmt goimports -pkg projectv2testutil -out clientmock_generated.go ../../projectv2 Client
//go:generate go tool mvdan.cc/gofumpt -lang=go1.24 -w clientmock_generated.go

//go:generate go tool github.com/matryer/moq -stub -fmt goimports -pkg projectv2testutil -out npmexecutormock_generated.go ../../projectv2 NpmExecutor
//go:generate go tool mvdan.cc/gofumpt -lang=go1.24 -w npmexecutormock_generated.go

const (
	TestTypingsLocation = "/home/src/Library/Caches/typescript"
)

type TestTypingsInstallerOptions struct {
	TypesRegistry []string
	PackageToFile map[string]string
}

type SessionUtils struct {
	fs          vfs.FS
	client      *ClientMock
	npmExecutor *NpmExecutorMock
	testOptions *TestTypingsInstallerOptions
	logs        strings.Builder
	logWriter   *bufio.Writer
}

func (h *SessionUtils) Client() *ClientMock {
	return h.client
}

func (h *SessionUtils) NpmExecutor() *NpmExecutorMock {
	return h.npmExecutor
}

func (h *SessionUtils) SetupNpmExecutorForTypingsInstaller() {
	if h.testOptions == nil {
		return
	}

	h.npmExecutor.NpmInstallFunc = func(cwd string, packageNames []string) ([]byte, error) {
		// packageNames is actually npmInstallArgs due to interface misnaming
		npmInstallArgs := packageNames
		lenNpmInstallArgs := len(npmInstallArgs)
		if lenNpmInstallArgs < 3 {
			return nil, fmt.Errorf("unexpected npm install: %s %v", cwd, npmInstallArgs)
		}

		if lenNpmInstallArgs == 3 && npmInstallArgs[2] == "types-registry@latest" {
			// Write typings file
			err := h.fs.WriteFile(cwd+"/node_modules/types-registry/index.json", h.createTypesRegistryFileContent(), false)
			return nil, err
		}

		// Find the packages: they start at index 2 and continue until we hit a flag starting with --
		packageEnd := lenNpmInstallArgs
		for i := 2; i < lenNpmInstallArgs; i++ {
			if strings.HasPrefix(npmInstallArgs[i], "--") {
				packageEnd = i
				break
			}
		}

		for _, atTypesPackageTs := range npmInstallArgs[2:packageEnd] {
			// @types/packageName@TsVersionToUse
			atTypesPackage := atTypesPackageTs
			// Remove version suffix
			if versionIndex := strings.LastIndex(atTypesPackage, "@"); versionIndex > 6 { // "@types/".length is 7, so version @ must be after
				atTypesPackage = atTypesPackage[:versionIndex]
			}
			// Extract package name from @types/packageName
			packageBaseName := atTypesPackage[7:] // Remove "@types/" prefix
			content, ok := h.testOptions.PackageToFile[packageBaseName]
			if !ok {
				return nil, fmt.Errorf("content not provided for %s", packageBaseName)
			}
			err := h.fs.WriteFile(cwd+"/node_modules/@types/"+packageBaseName+"/index.d.ts", content, false)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
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

func (h *SessionUtils) FS() vfs.FS {
	return h.fs
}

func (h *SessionUtils) Log(msg ...any) {
	fmt.Fprintln(&h.logs, msg...)
}

func (h *SessionUtils) Logs() string {
	h.logWriter.Flush()
	return h.logs.String()
}

func (h *SessionUtils) BaselineLogs(t *testing.T) {
	baseline.Run(t, t.Name()+".log", h.Logs(), baseline.Options{
		Subfolder: "project",
	})
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

func SetupWithOptions(files map[string]any, options *projectv2.SessionOptions) (*projectv2.Session, *SessionUtils) {
	return SetupWithOptionsAndTypingsInstaller(files, options, nil)
}

func SetupWithTypingsInstaller(files map[string]any, testOptions *TestTypingsInstallerOptions) (*projectv2.Session, *SessionUtils) {
	return SetupWithOptionsAndTypingsInstaller(files, nil, testOptions)
}

func SetupWithOptionsAndTypingsInstaller(files map[string]any, options *projectv2.SessionOptions, testOptions *TestTypingsInstallerOptions) (*projectv2.Session, *SessionUtils) {
	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	clientMock := &ClientMock{}
	npmExecutorMock := &NpmExecutorMock{}
	sessionUtils := &SessionUtils{
		fs:          fs,
		client:      clientMock,
		npmExecutor: npmExecutorMock,
		testOptions: testOptions,
	}
	sessionUtils.logWriter = bufio.NewWriter(&sessionUtils.logs)

	// Configure the npm executor mock to handle typings installation
	sessionUtils.SetupNpmExecutorForTypingsInstaller()

	// Use provided options or create default ones
	sessionOptions := options
	if sessionOptions == nil {
		sessionOptions = &projectv2.SessionOptions{
			CurrentDirectory:   "/",
			DefaultLibraryPath: bundled.LibPath(),
			TypingsLocation:    TestTypingsLocation,
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       true,
			LoggingEnabled:     true,
		}
	}

	session := projectv2.NewSession(&projectv2.SessionInit{
		Options:     sessionOptions,
		FS:          fs,
		Client:      clientMock,
		NpmExecutor: npmExecutorMock,
		Logger:      sessionUtils,
	})

	return session, sessionUtils
}
