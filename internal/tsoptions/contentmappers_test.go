package tsoptions

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type resolveContentMapperHost struct {
	fs vfs.FS
}

func TestGetContentMapperForFileNameUsesLongestExtension(t *testing.T) {
	t.Parallel()
	zMapper := &contentmapper.Mapper{Definition: contentmapper.Definition{Package: "z", Extensions: []string{".z"}}}
	yzMapper := &contentmapper.Mapper{Definition: contentmapper.Definition{Package: "yz", Extensions: []string{".y.z"}}}
	commandLine := &ParsedCommandLine{ParsedConfig: &ParsedOptions{ContentMappers: []*contentmapper.Mapper{zMapper, yzMapper}}}

	assert.Equal(t, commandLine.GetContentMapperForFileName("/src/Component.y.z"), yzMapper)
	assert.Equal(t, commandLine.GetContentMapperForFileName("/src/Component.z"), zMapper)
}
func (h resolveContentMapperHost) FS() vfs.FS                  { return h.fs }
func (h resolveContentMapperHost) GetCurrentDirectory() string { return "/home/project" }

func TestResolveContentMapperManifest(t *testing.T) {
	t.Parallel()

	host := resolveContentMapperHost{fs: vfstest.FromMap(map[string]string{
		"/home/project/node_modules/vue-ts-mapper/package.json": `{
			"name": "vue-ts-mapper",
			"version": "1.2.3",
			"tsContentMapper": { "exec": ["node", "./dist/mapper.js"], "extensions": { ".vue": ".tsx" }, "compilerOptions": ["target", "jsx"] }
		}`,
		"/home/node_modules/@scope/noversion/package.json": `{
			"name": "@scope/noversion",
			"tsContentMapper": { "exec": ["run"], "extensions": { ".component": ".mts" } }
		}`,
		"/home/project/node_modules/no-name/package.json": `{
			"version": "1.0.0"
		}`,
		"/home/project/node_modules/no-manifest/package.json": `{
			"name": "no-manifest"
		}`,
		"/home/project/node_modules/no-exec/package.json": `{
			"name": "no-exec",
			"tsContentMapper": {}
		}`,
		"/home/project/node_modules/bad-exec/package.json": `{
			"name": "bad-exec",
			"tsContentMapper": { "exec": "node ./mapper.js" }
		}`,
		"/home/project/node_modules/no-extensions/package.json": `{
			"name": "no-extensions",
			"tsContentMapper": { "exec": ["run"] }
		}`,
		"/home/project/node_modules/bad-source-extension/package.json": `{
			"name": "bad-source-extension",
			"tsContentMapper": { "exec": ["run"], "extensions": { "vue": ".ts" } }
		}`,
		"/home/project/node_modules/bad-virtual-extension/package.json": `{
			"name": "bad-virtual-extension",
			"tsContentMapper": { "exec": ["run"], "extensions": { ".vue": ".coffee" } }
		}`,
	}, true /*useCaseSensitiveFileNames*/)}

	// Name, version, and the verbatim exec argv are preserved.
	manifest, packageDirectory, diagnostic := resolveContentMapperManifest(host, "/home/project/tsconfig.json", "vue-ts-mapper")
	assert.Assert(t, diagnostic == nil)
	assert.Equal(t, manifest.Name, "vue-ts-mapper")
	assert.Equal(t, manifest.Version, "1.2.3")
	assert.Equal(t, packageDirectory, "/home/project/node_modules/vue-ts-mapper")
	assert.DeepEqual(t, manifest.Exec, []string{"node", "./dist/mapper.js"})
	assert.DeepEqual(t, manifest.CompilerOptions, []string{"target", "jsx"})
	assert.DeepEqual(t, manifest.Extensions, map[string]string{".vue": ".tsx"})

	// Resolution walks up node_modules; a package with no version resolves to a name and empty version.
	manifest, _, diagnostic = resolveContentMapperManifest(host, "/home/project/src/tsconfig.json", "@scope/noversion")
	assert.Assert(t, diagnostic == nil)
	assert.Equal(t, manifest.Name, "@scope/noversion")
	assert.Equal(t, manifest.Version, "")

	// A package that is not installed reports a resolution diagnostic.
	_, _, diagnostic = resolveContentMapperManifest(host, "/home/project/tsconfig.json", "missing-mapper")
	assert.Assert(t, diagnostic != nil)
	assert.Equal(t, diagnostic.Code(), diagnostics.The_content_mapper_package_0_could_not_be_resolved.Code())

	// A package whose package.json has no name reports a diagnostic.
	_, _, diagnostic = resolveContentMapperManifest(host, "/home/project/tsconfig.json", "no-name")
	assert.Assert(t, diagnostic != nil)
	assert.Equal(t, diagnostic.Code(), diagnostics.The_package_json_of_the_content_mapper_package_0_does_not_specify_a_name.Code())

	// A package that does not declare a "tsContentMapper" object reports a diagnostic.
	_, _, diagnostic = resolveContentMapperManifest(host, "/home/project/tsconfig.json", "no-manifest")
	assert.Assert(t, diagnostic != nil)
	assert.Equal(t, diagnostic.Code(), diagnostics.The_package_json_of_the_content_mapper_package_0_does_not_declare_a_tsContentMapper_object.Code())

	// A "tsContentMapper" with no "exec", or an "exec" of the wrong type, reports a diagnostic.
	for _, pkg := range []string{"no-exec", "bad-exec"} {
		_, _, diagnostic = resolveContentMapperManifest(host, "/home/project/tsconfig.json", pkg)
		assert.Assert(t, diagnostic != nil, "expected a diagnostic for %s", pkg)
		assert.Equal(t, diagnostic.Code(), diagnostics.The_tsContentMapper_exec_of_the_content_mapper_package_0_must_be_a_non_empty_array_of_strings.Code())
	}

	_, _, diagnostic = resolveContentMapperManifest(host, "/home/project/tsconfig.json", "no-extensions")
	assert.Assert(t, diagnostic != nil)
	assert.Equal(t, diagnostic.Code(), diagnostics.The_tsContentMapper_extensions_of_the_content_mapper_package_0_must_be_a_non_empty_object_mapping_source_extensions_to_virtual_extensions.Code())

	_, _, diagnostic = resolveContentMapperManifest(host, "/home/project/tsconfig.json", "bad-source-extension")
	assert.Assert(t, diagnostic != nil)
	assert.Equal(t, diagnostic.Code(), diagnostics.The_source_extension_0_in_tsContentMapper_extensions_of_the_content_mapper_package_1_must_begin_with_a.Code())

	_, _, diagnostic = resolveContentMapperManifest(host, "/home/project/tsconfig.json", "bad-virtual-extension")
	assert.Assert(t, diagnostic != nil)
	assert.Equal(t, diagnostic.Code(), diagnostics.The_virtual_extension_0_for_source_extension_1_in_tsContentMapper_extensions_of_the_content_mapper_package_2_must_be_one_of_Colon_3.Code())
}
