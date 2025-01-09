package module_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/module"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

var skip = []string{}

type vfsModuleResolutionHost struct {
	mu               sync.Mutex
	fs               vfs.FS
	currentDirectory string
	traces           []string
}

func fixRoot(path string) string {
	rootLength := tspath.GetRootLength(path)
	if rootLength == 0 {
		return tspath.CombinePaths(".src", path)
	}
	if len(path) == rootLength {
		return "."
	}
	return path[rootLength:]
}

func newVFSModuleResolutionHost(files map[string]string, currentDirectory string) *vfsModuleResolutionHost {
	fs := fstest.MapFS{}
	for name, content := range files {
		fs[fixRoot(name)] = &fstest.MapFile{
			Data: []byte(content),
		}
	}
	if currentDirectory == "" {
		currentDirectory = "/.src"
	} else if currentDirectory[0] != '/' {
		currentDirectory = "/.src/" + currentDirectory
	}
	return &vfsModuleResolutionHost{
		fs:               vfstest.FromMapFS(fs, true /*useCaseSensitiveFileNames*/),
		currentDirectory: currentDirectory,
	}
}

func (v *vfsModuleResolutionHost) FS() vfs.FS {
	return v.fs
}

// GetCurrentDirectory implements ModuleResolutionHost.
func (v *vfsModuleResolutionHost) GetCurrentDirectory() string {
	return v.currentDirectory
}

// Trace implements ModuleResolutionHost.
func (v *vfsModuleResolutionHost) Trace(msg string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.traces = append(v.traces, msg)
}

type functionCall struct {
	call        string
	args        rawArgs
	returnValue map[string]any
}
type traceTestCase struct {
	name             string
	currentDirectory string
	trace            bool
	compilerOptions  *core.CompilerOptions
	files            map[string]string
	calls            []functionCall
}
type rawFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
type rawArgs struct {
	// getPackageScopeForPath
	Directory string `json:"directory"`

	// resolveModuleName, resolveTypeReferenceDirective
	Name            string                `json:"name"`
	ContainingFile  string                `json:"containingFile"`
	CompilerOptions *core.CompilerOptions `json:"compilerOptions"`
	ResolutionMode  int                   `json:"resolutionMode"`
	RedirectedRef   *struct {
		SourceFile struct {
			FileName string `json:"fileName"`
		} `json:"sourceFile"`
		CommandLine struct {
			Options *core.CompilerOptions `json:"options"`
		} `json:"commandLine"`
	} `json:"redirectedReference"`
}
type rawTest struct {
	Test             string         `json:"test"`
	CurrentDirectory string         `json:"currentDirectory"`
	Trace            bool           `json:"trace"`
	Files            []rawFile      `json:"files"`
	Call             string         `json:"call"`
	Args             rawArgs        `json:"args"`
	Return           map[string]any `json:"return"`
}

var typesVersionsMessageRegex = regexp.MustCompile(`that matches compiler version '[^']+'`)

func sanitizeTraceOutput(trace string) string {
	return typesVersionsMessageRegex.ReplaceAllString(trace, "that matches compiler version '3.1.0-dev'")
}

func doCall(t *testing.T, resolver *module.Resolver, call functionCall, skipLocations bool) {
	switch call.call {
	case "resolveModuleName", "resolveTypeReferenceDirective":
		var redirectedReference *module.ResolvedProjectReference
		if call.args.RedirectedRef != nil {
			redirectedReference = &module.ResolvedProjectReference{
				SourceFile: (&ast.NodeFactory{}).NewSourceFile("", call.args.RedirectedRef.SourceFile.FileName, nil).AsSourceFile(),
				CommandLine: core.ParsedOptions{
					Options: call.args.RedirectedRef.CommandLine.Options,
				},
			}
		}

		var locations *module.LookupLocations
		errorMessageArgs := []any{call.args.Name, call.args.ContainingFile}
		if call.call == "resolveModuleName" {
			resolved := resolver.ResolveModuleName(call.args.Name, call.args.ContainingFile, core.ModuleKind(call.args.ResolutionMode), redirectedReference)
			assert.Check(t, resolved != nil, "ResolveModuleName should not return nil", errorMessageArgs)
			locations = resolver.GetLookupLocationsForResolvedModule(resolved)
			if expectedResolvedModule, ok := call.returnValue["resolvedModule"].(map[string]any); ok {
				assert.Check(t, resolved.IsResolved(), errorMessageArgs)
				assert.Check(t, cmp.Equal(resolved.ResolvedFileName, expectedResolvedModule["resolvedFileName"].(string)), errorMessageArgs)
				assert.Check(t, cmp.Equal(resolved.Extension, expectedResolvedModule["extension"].(string)), errorMessageArgs)
				assert.Check(t, cmp.Equal(resolved.ResolvedUsingTsExtension, expectedResolvedModule["resolvedUsingTsExtension"].(bool)), errorMessageArgs)
				assert.Check(t, cmp.Equal(resolved.IsExternalLibraryImport, expectedResolvedModule["isExternalLibraryImport"].(bool)), errorMessageArgs)
			} else {
				assert.Check(t, !resolved.IsResolved(), errorMessageArgs)
			}
		} else {
			resolved := resolver.ResolveTypeReferenceDirective(call.args.Name, call.args.ContainingFile, core.ModuleKind(call.args.ResolutionMode), redirectedReference)
			assert.Check(t, resolved != nil, "ResolveTypeReferenceDirective should not return nil", errorMessageArgs)
			locations = resolver.GetLookupLocationsForResolvedTypeReferenceDirective(resolved)
			if expectedResolvedTypeReferenceDirective, ok := call.returnValue["resolvedTypeReferenceDirective"].(map[string]any); ok {
				assert.Check(t, resolved.IsResolved(), errorMessageArgs)
				assert.Check(t, cmp.Equal(resolved.ResolvedFileName, expectedResolvedTypeReferenceDirective["resolvedFileName"].(string)), errorMessageArgs)
				assert.Check(t, cmp.Equal(resolved.Primary, expectedResolvedTypeReferenceDirective["primary"].(bool)), errorMessageArgs)
				assert.Check(t, cmp.Equal(resolved.IsExternalLibraryImport, expectedResolvedTypeReferenceDirective["isExternalLibraryImport"].(bool)), errorMessageArgs)
			} else {
				assert.Check(t, !resolved.IsResolved(), errorMessageArgs)
			}
		}
		if skipLocations {
			break
		}
		if expectedFailedLookupLocations, ok := call.returnValue["failedLookupLocations"].([]any); ok {
			assert.Check(t, cmp.DeepEqual(locations.FailedLookupLocations, core.Map(expectedFailedLookupLocations, func(i any) string { return i.(string) })), errorMessageArgs)
		} else {
			assert.Check(t, cmp.Equal(len(locations.FailedLookupLocations), 0), errorMessageArgs)
		}
		if expectedAffectingLocations, ok := call.returnValue["affectingLocations"].([]any); ok {
			assert.Check(t, cmp.DeepEqual(locations.AffectingLocations, core.Map(expectedAffectingLocations, func(i any) string { return i.(string) })), errorMessageArgs)
		} else {
			assert.Check(t, cmp.Equal(len(locations.AffectingLocations), 0), errorMessageArgs)
		}
	case "getPackageScopeForPath":
		resolver.GetPackageScopeForPath(call.args.Directory)
	default:
		t.Errorf("Unexpected call: %s", call.call)
	}
}

func runTraceBaseline(t *testing.T, test traceTestCase) {
	t.Run(test.name, func(t *testing.T) {
		t.Parallel()

		host := newVFSModuleResolutionHost(test.files, test.currentDirectory)
		resolver := module.NewResolver(host, test.compilerOptions)

		for _, call := range test.calls {
			doCall(t, resolver, call, false /*skipLocations*/)
			if t.Failed() {
				t.FailNow()
			}
		}

		t.Run("concurrent", func(t *testing.T) {
			host := newVFSModuleResolutionHost(test.files, test.currentDirectory)
			resolver := module.NewResolver(host, test.compilerOptions)

			var wg sync.WaitGroup
			for _, call := range test.calls {
				wg.Add(1)
				go func() {
					defer wg.Done()
					doCall(t, resolver, call, true /*skipLocations*/)
				}()
			}

			wg.Wait()
		})

		if test.trace {
			t.Run("trace", func(t *testing.T) {
				var buf bytes.Buffer
				encoder := json.NewEncoder(&buf)
				encoder.SetIndent("", "    ")
				encoder.SetEscapeHTML(false)
				if err := encoder.Encode(host.traces); err != nil {
					t.Fatal(err)
				}
				baseline.Run(
					t,
					tspath.RemoveFileExtension(test.name)+".trace.json",
					sanitizeTraceOutput(buf.String()),
					baseline.Options{Subfolder: "module/resolver"},
				)
			})
		}
	})
}

func TestModuleResolver(t *testing.T) {
	t.Parallel()
	testsFilePath := filepath.Join(repo.TestDataPath, "fixtures", "module", "resolvertests.json")
	// Read file one line at a time
	file, err := os.Open(testsFilePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		file.Close()
	})
	decoder := json.NewDecoder(file)
	var currentTestCase traceTestCase
	for {
		var json rawTest
		if err := decoder.Decode(&json); err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatal(err)
		}
		if json.Files != nil {
			if currentTestCase.name != "" && !slices.Contains(skip, currentTestCase.name) {
				runTraceBaseline(t, currentTestCase)
			}
			currentTestCase = traceTestCase{
				name:             json.Test,
				currentDirectory: json.CurrentDirectory,
				// !!! no traces are passing yet because of missing cache implementation
				trace: false,
				files: make(map[string]string, len(json.Files)),
			}
			for _, file := range json.Files {
				currentTestCase.files[file.Name] = file.Content
			}
		} else if json.Call != "" {
			currentTestCase.calls = append(currentTestCase.calls, functionCall{
				call:        json.Call,
				args:        json.Args,
				returnValue: json.Return,
			})
			if currentTestCase.compilerOptions == nil && json.Args.CompilerOptions != nil {
				currentTestCase.compilerOptions = json.Args.CompilerOptions
			}
		} else {
			t.Fatalf("Unexpected JSON: %v", json)
		}
	}
}
