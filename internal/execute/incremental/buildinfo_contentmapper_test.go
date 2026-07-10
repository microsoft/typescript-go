package incremental_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"gotest.tools/v3/assert"
)

func configWithMappers(mappers ...*contentmapper.Mapper) *tsoptions.ParsedCommandLine {
	return &tsoptions.ParsedCommandLine{
		ParsedConfig: &core.ParsedOptions{
			CompilerOptions: &core.CompilerOptions{},
			ContentMappers:  mappers,
		},
	}
}

func TestContentMapperIdentities(t *testing.T) {
	t.Parallel()

	assert.Assert(t, incremental.ContentMapperIdentities(configWithMappers()) == nil)

	// Identities come from the mappers' resolved name/version; a mapper with no name is omitted, and the
	// result is sorted so reordering content mappers in tsconfig does not change it.
	config := configWithMappers(
		&contentmapper.Mapper{Definition: contentmapper.Definition{Package: "vue"}, Manifest: contentmapper.Manifest{Name: "vue", Version: "2.0.0"}},
		&contentmapper.Mapper{Definition: contentmapper.Definition{Package: "svelte"}, Manifest: contentmapper.Manifest{Name: "svelte", Version: "3.0.0"}},
		&contentmapper.Mapper{Definition: contentmapper.Definition{Package: "anon"}},
	)
	assert.DeepEqual(t, incremental.ContentMapperIdentities(config), []string{"svelte@3.0.0", "vue@2.0.0"})
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
	// config now resolves to cannot be reused: the old program is discarded (nil) so the project is
	// rebuilt. The host is never touched because we bail before reconstructing the snapshot.
	buildInfo := &incremental.BuildInfo{
		Version:                 core.Version(),
		FileNames:               []string{"/src/a.ts"},
		ContentMapperIdentities: []string{"vue@1.0.0"},
	}
	config := configWithMappers(&contentmapper.Mapper{Definition: contentmapper.Definition{Package: "vue", Extensions: []string{".vue"}}, Manifest: contentmapper.Manifest{Name: "vue", Version: "2.0.0"}})

	program := incremental.ReadBuildInfoProgram(config, fakeBuildInfoReader{buildInfo}, nil)
	assert.Assert(t, program == nil, "expected the old program to be discarded when the mapper identity changed")
}
