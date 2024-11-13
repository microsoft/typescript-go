package sourcemap

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestSourceMapGeneratorEmpty(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{},
		Mappings:       "",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorSerializeEmpty(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	actual := gen.String()
	expected := `{"version":3,"file":"main.js","sourceRoot":"/","sources":[],"mappings":""}`
	assert.Equal(t, actual, expected)
}

func TestSourceMapGeneratorAddSource(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	sourceMap := gen.RawSourceMap()
	assert.Equal(t, int(sourceIndex), 0)
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorSetSourceContent(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	sourceContent := "foo"
	gen.SetSourceContent(sourceIndex, sourceContent)
	sourceMap := gen.RawSourceMap()
	assert.Equal(t, int(sourceIndex), 0)
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "",
		Names:          nil,
		SourcesContent: []*string{&sourceContent},
	})
}

func TestSourceMapGeneratorSetSourceContentForSecondSourceOnly(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	gen.AddSource("/skipped.ts")
	sourceIndex := gen.AddSource("/main.ts")
	sourceContent := "foo"
	gen.SetSourceContent(sourceIndex, sourceContent)
	sourceMap := gen.RawSourceMap()
	assert.Equal(t, int(sourceIndex), 1)
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"skipped.ts", "main.ts"},
		Mappings:       "",
		Names:          nil,
		SourcesContent: []*string{nil, &sourceContent},
	})
}

func TestSourceMapGeneratorSerializeSetSourceContentForSecondSourceOnly(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	gen.AddSource("/skipped.ts")
	sourceIndex := gen.AddSource("/main.ts")
	sourceContent := "foo"
	gen.SetSourceContent(sourceIndex, sourceContent)
	actual := gen.String()
	expected := `{"version":3,"file":"main.js","sourceRoot":"/","sources":["skipped.ts","main.ts"],"mappings":"","sourcesContent":[null,"foo"]}`
	assert.Equal(t, actual, expected)
}

func TestSourceMapGeneratorAddName(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	nameIndex := gen.AddName("foo")
	sourceMap := gen.RawSourceMap()
	assert.Equal(t, int(nameIndex), 0)
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{},
		Mappings:       "",
		Names:          []string{"foo"},
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMapping(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	gen.AddMapping(0, 0)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{},
		Mappings:       "A",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingOnSecondLineOnly(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	gen.AddMapping(1, 0)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{},
		Mappings:       ";A",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingToSource(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	gen.AddMappingSource(0, 0, sourceIndex, 0, 0)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "AAAA",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingToSourceNextGeneratedCharacter(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	gen.AddMappingSource(0, 0, sourceIndex, 0, 0)
	gen.AddMappingSource(0, 1, sourceIndex, 0, 0)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "AAAA,CAAA",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingToSourceNextGeneratedAndSourceCharacter(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	gen.AddMappingSource(0, 0, sourceIndex, 0, 0)
	gen.AddMappingSource(0, 1, sourceIndex, 0, 1)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "AAAA,CAAC",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingToSourceNextGeneratedLine(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	gen.AddMappingSource(0, 0, sourceIndex, 0, 0)
	gen.AddMappingSource(1, 0, sourceIndex, 0, 0)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "AAAA;AAAA",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingToSourcePreviousSourceCharacter(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	gen.AddMappingSource(0, 0, sourceIndex, 0, 1)
	gen.AddMappingSource(0, 1, sourceIndex, 0, 0)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "AAAC,CAAD",
		Names:          nil,
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingToSourceWithName(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	nameIndex := gen.AddName("foo")
	gen.AddMappingSourceName(0, 0, sourceIndex, 0, 0, nameIndex)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "AAAAA",
		Names:          []string{"foo"},
		SourcesContent: nil,
	})
}

func TestSourceMapGeneratorAddMappingToSourceWithPreviousName(t *testing.T) {
	t.Parallel()
	gen := NewSourceMapGenerator("main.js", "/", "/", tspath.ComparePathsOptions{})
	sourceIndex := gen.AddSource("/main.ts")
	nameIndex1 := gen.AddName("foo")
	nameIndex2 := gen.AddName("bar")
	gen.AddMappingSourceName(0, 0, sourceIndex, 0, 0, nameIndex2)
	gen.AddMappingSourceName(0, 1, sourceIndex, 0, 0, nameIndex1)
	sourceMap := gen.RawSourceMap()
	assert.DeepEqual(t, sourceMap, &RawSourceMap{
		Version:        3,
		File:           "main.js",
		SourceRoot:     "/",
		Sources:        []string{"main.ts"},
		Mappings:       "AAAAC,CAAAD",
		Names:          []string{"foo", "bar"},
		SourcesContent: nil,
	})
}
