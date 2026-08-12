package contentmappertest

import "fmt"

const PackageName = "mapper"

// PackageJSON returns a package manifest selecting the requested mapper.
func PackageJSON(mapper string) string {
	compilerOptions := ""
	dynamicConfig := ""
	if mapper == TransformingMapper {
		compilerOptions = `, "compilerOptions": ["target", "jsx"]`
	}
	if mapper == DynamicVerbatimMapper {
		dynamicConfig = `, "dynamicConfig": true`
	}
	return fmt.Sprintf(`{
	"name": %q,
	"version": "1.0.0",
	"tsContentMapper": { "exec": [%q], "extensions": { ".astro": ".ts", ".box": ".ts", ".dup": ".ts", ".lisp": ".ts", ".mdx": ".ts", ".panel": ".ts", ".svelte": ".ts", ".vue": ".ts", ".y.z": ".ts" }%s%s }
}`, PackageName, mapper, compilerOptions, dynamicConfig)
}
