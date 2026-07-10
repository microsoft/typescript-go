package incremental_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"gotest.tools/v3/assert"
)

func configWithMappers(mappers ...*core.ContentMapper) *tsoptions.ParsedCommandLine {
	return &tsoptions.ParsedCommandLine{
		ParsedConfig: &core.ParsedOptions{
			CompilerOptions: &core.CompilerOptions{},
			ContentMappers:  mappers,
		},
	}
}

// fakeContentMapperRunner reports a fixed identity for every mapper; Transform is never exercised here.
type fakeContentMapperRunner struct {
	identity string
}

func (r fakeContentMapperRunner) Identity(*core.ContentMapper) string {
	return r.identity
}

func (fakeContentMapperRunner) Transform(fileName string, content string) (compiler.ContentMapperResult, error) {
	return compiler.ContentMapperResult{}, nil
}

func TestContentMapperIdentitiesMatch(t *testing.T) {
	t.Parallel()

	buildInfo := &incremental.BuildInfo{ContentMapperIdentities: []string{"vue@1.0.0"}}
	assert.Assert(t, buildInfo.ContentMapperIdentitiesMatch([]string{"vue@1.0.0"}))
	assert.Assert(t, !buildInfo.ContentMapperIdentitiesMatch([]string{"vue@2.0.0"}))
	assert.Assert(t, !buildInfo.ContentMapperIdentitiesMatch(nil))

	empty := &incremental.BuildInfo{}
	assert.Assert(t, empty.ContentMapperIdentitiesMatch(nil))
	assert.Assert(t, !empty.ContentMapperIdentitiesMatch([]string{"vue@1.0.0"}))
}

type fakeBuildInfoReader struct {
	buildInfo *incremental.BuildInfo
}

func (r fakeBuildInfoReader) ReadBuildInfo(*tsoptions.ParsedCommandLine) *incremental.BuildInfo {
	return r.buildInfo
}

func TestReadBuildInfoProgramContentMapperIdentityMismatch(t *testing.T) {
	t.Parallel()

	// An otherwise-valid, incremental build info whose recorded mapper identity differs from the one the
	// runner now reports cannot be reused: the old program is discarded (nil) so the project is rebuilt.
	// The host is never touched because we bail before reconstructing the snapshot.
	buildInfo := &incremental.BuildInfo{
		Version:                 core.Version(),
		FileNames:               []string{"/src/a.ts"},
		ContentMapperIdentities: []string{"vue@1.0.0"},
	}
	config := configWithMappers(&core.ContentMapper{Extensions: []string{".vue"}})
	runner := fakeContentMapperRunner{identity: "vue@2.0.0"}

	program := incremental.ReadBuildInfoProgram(config, runner, fakeBuildInfoReader{buildInfo}, nil)
	assert.Assert(t, program == nil, "expected the old program to be discarded when the mapper identity changed")
}
