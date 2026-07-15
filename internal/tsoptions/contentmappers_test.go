package tsoptions

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type resolveContentMapperHost struct {
	fs vfs.FS
}

func (h resolveContentMapperHost) FS() vfs.FS                  { return h.fs }
func (h resolveContentMapperHost) GetCurrentDirectory() string { return "/home/project" }

func TestResolveContentMapperManifest(t *testing.T) {
	t.Parallel()

	host := resolveContentMapperHost{fs: vfstest.FromMap(map[string]string{
		"/home/project/node_modules/vue-ts-mapper/package.json": `{
			"name": "vue-ts-mapper",
			"version": "1.2.3",
			"tsContentMapper": { "exec": ["node", "./dist/mapper.js"], "compilerOptions": ["target", "jsx"] }
		}`,
		"/home/node_modules/@scope/noversion/package.json": `{
			"name": "@scope/noversion",
			"tsContentMapper": { "exec": ["run"] }
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
	}, true /*useCaseSensitiveFileNames*/)}

	// Name, version, and the verbatim exec argv are preserved.
	manifest, packageDirectory, diagnostic := resolveContentMapperManifest(host, "/home/project/tsconfig.json", "vue-ts-mapper")
	assert.Assert(t, diagnostic == nil)
	assert.Equal(t, manifest.Name, "vue-ts-mapper")
	assert.Equal(t, manifest.Version, "1.2.3")
	assert.Equal(t, packageDirectory, "/home/project/node_modules/vue-ts-mapper")
	assert.DeepEqual(t, manifest.Exec, []string{"node", "./dist/mapper.js"})
	assert.DeepEqual(t, manifest.CompilerOptions, []string{"target", "jsx"})

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
}
