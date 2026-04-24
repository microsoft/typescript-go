package project

import (
	"encoding/json"
	"testing"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestCollectProjectInfoTelemetryStrictArrayVariance(t *testing.T) {
	t.Parallel()

	cmp := tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: true,
		CurrentDirectory:          "/tmp",
	}

	t.Run("true is serialized in compilerOptions", func(t *testing.T) {
		t.Parallel()
		co := &core.CompilerOptions{StrictArrayVariance: core.TSTrue}
		cmd := tsoptions.NewParsedCommandLine(co, nil, cmp)
		proj := &Project{
			Kind:           KindConfigured,
			configFileName: "/tmp/tsconfig.json",
			CommandLine:    cmd,
			Program:        &compiler.Program{},
		}
		var s Session
		ev := s.collectProjectInfoTelemetry(proj)
		assert.Assert(t, ev.ProjectInfoTelemetryEvent != nil)
		raw := ev.ProjectInfoTelemetryEvent.Properties["compilerOptions"]
		var m map[string]any
		assert.NilError(t, json.Unmarshal([]byte(raw), &m))
		assert.Equal(t, true, m["strictArrayVariance"])
	})

	t.Run("false is serialized in compilerOptions", func(t *testing.T) {
		t.Parallel()
		co := &core.CompilerOptions{StrictArrayVariance: core.TSFalse}
		cmd := tsoptions.NewParsedCommandLine(co, nil, cmp)
		proj := &Project{
			Kind:           KindConfigured,
			configFileName: "/tmp/tsconfig.json",
			CommandLine:    cmd,
			Program:        &compiler.Program{},
		}
		var s Session
		ev := s.collectProjectInfoTelemetry(proj)
		assert.Assert(t, ev.ProjectInfoTelemetryEvent != nil)
		raw := ev.ProjectInfoTelemetryEvent.Properties["compilerOptions"]
		var m map[string]any
		assert.NilError(t, json.Unmarshal([]byte(raw), &m))
		assert.Equal(t, false, m["strictArrayVariance"])
	})

	t.Run("unknown is omitted from compilerOptions", func(t *testing.T) {
		t.Parallel()
		co := &core.CompilerOptions{Strict: core.TSTrue}
		cmd := tsoptions.NewParsedCommandLine(co, nil, cmp)
		proj := &Project{
			Kind:           KindConfigured,
			configFileName: "/tmp/tsconfig.json",
			CommandLine:    cmd,
			Program:        &compiler.Program{},
		}
		var s Session
		ev := s.collectProjectInfoTelemetry(proj)
		assert.Assert(t, ev.ProjectInfoTelemetryEvent != nil)
		raw := ev.ProjectInfoTelemetryEvent.Properties["compilerOptions"]
		var m map[string]any
		assert.NilError(t, json.Unmarshal([]byte(raw), &m))
		_, ok := m["strictArrayVariance"]
		assert.Assert(t, !ok, "strictArrayVariance should be absent when Tristate is unknown")
	})

	t.Run("jsconfig.json sets configFileName telemetry bucket", func(t *testing.T) {
		t.Parallel()
		co := &core.CompilerOptions{StrictArrayVariance: core.TSTrue}
		cmd := tsoptions.NewParsedCommandLine(co, nil, cmp)
		proj := &Project{
			Kind:           KindConfigured,
			configFileName: "/tmp/jsconfig.json",
			CommandLine:    cmd,
			Program:        &compiler.Program{},
		}
		var s Session
		ev := s.collectProjectInfoTelemetry(proj)
		assert.Assert(t, ev.ProjectInfoTelemetryEvent != nil)
		assert.Equal(t, "jsconfig.json", ev.ProjectInfoTelemetryEvent.Properties["configFileName"])
		raw := ev.ProjectInfoTelemetryEvent.Properties["compilerOptions"]
		var m map[string]any
		assert.NilError(t, json.Unmarshal([]byte(raw), &m))
		assert.Equal(t, true, m["strictArrayVariance"])
	})
}
