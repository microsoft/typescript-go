package tsctests

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func symlinkedProjectReferenceInput(withPaths bool) *tscInput {
	baseConfig := `{
  "compilerOptions": {
    "target": "es2022",
    "module": "node16",
    "moduleResolution": "node16",
    "strict": true,
    "declaration": true,
    "composite": true,
    "skipLibCheck": true
  }
}`
	if withPaths {
		baseConfig = `{
  "compilerOptions": {
    "target": "es2022",
    "module": "node16",
    "moduleResolution": "node16",
    "strict": true,
    "declaration": true,
    "composite": true,
    "skipLibCheck": true,
    "paths": {
      "@lab/feature-gating": ["./packages/feature-gating/src/index.ts"]
    }
  }
}`
	}

	return &tscInput{
		cwd: "/home/src/workspaces/project",
		files: FileMap{
			"/home/src/workspaces/project/tsconfig.base.json": baseConfig,
			"/home/src/workspaces/project/packages/feature-gating/package.json": `{
  "name": "@lab/feature-gating",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "main": "./dist/index.js",
  "types": "./dist/index.d.ts",
  "exports": {
    ".": {
      "types": "./dist/index.d.ts",
      "default": "./dist/index.js"
    }
  }
}`,
			"/home/src/workspaces/project/packages/feature-gating/tsconfig.json": `{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "rootDir": "src",
    "outDir": "dist",
    "tsBuildInfoFile": "dist/.tsbuildinfo"
  },
  "include": ["src/**/*.ts"]
}`,
			"/home/src/workspaces/project/packages/feature-gating/src/index.ts": `import type { State } from "./types.js";

export type Selector<TState, TResult> = (state: TState) => TResult;

export const createFeatureGateSelector =
  (featureGate: string) =>
  (state: State): boolean =>
    state.featureGates[0] === featureGate;`,
			"/home/src/workspaces/project/packages/feature-gating/src/types.ts": `export interface State {
  featureGates: string[];
}`,
			"/home/src/workspaces/project/packages/app/package.json": `{
  "name": "@lab/app",
  "version": "0.0.0",
  "private": true,
  "type": "module",
  "main": "./dist/index.js",
  "types": "./dist/index.d.ts",
  "dependencies": {
    "@lab/feature-gating": "workspace:*"
  }
}`,
			"/home/src/workspaces/project/packages/app/tsconfig.json": `{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "rootDir": "src",
    "outDir": "dist",
    "tsBuildInfoFile": "dist/.tsbuildinfo"
  },
  "references": [
    { "path": "../feature-gating" }
  ],
  "include": ["src/**/*.ts"]
}`,
			"/home/src/workspaces/project/packages/app/src/index.ts": `import { createFeatureGateSelector } from "@lab/feature-gating";

export const isFooEnabled = createFeatureGateSelector("foo");`,
			"/home/src/workspaces/project/packages/app/node_modules/@lab/feature-gating": vfstest.Symlink("/home/src/workspaces/project/packages/feature-gating"),
		},
	}
}

func TestBuildModeSymlinkedProjectReferencePrefersProjectPathWhenPathsExist(t *testing.T) {
	t.Parallel()

	sys := newTestSys(symlinkedProjectReferenceInput(true), false)
	result := execute.CommandLine(sys, []string{"-b", "packages/app/tsconfig.json"}, sys)
	if result.Status != tsc.ExitStatusSuccess {
		t.Fatalf("expected build to succeed, got %v\noutput:\n%s", result.Status, sys.getOutput(true))
	}

	dts, ok := sys.fsFromFileMap().ReadFile("/home/src/workspaces/project/packages/app/dist/index.d.ts")
	if !ok {
		t.Fatal("expected app declaration output to be written")
	}
	if !strings.Contains(dts, `import("../../feature-gating/src/types.js").State`) {
		t.Fatalf("expected declaration emit to reference the project path, got:\n%s", dts)
	}
	if strings.Contains(dts, "/node_modules/") {
		t.Fatalf("expected declaration emit to avoid node_modules symlink paths, got:\n%s", dts)
	}
}

func TestBuildModeSymlinkedProjectReferenceStillErrorsWithoutPaths(t *testing.T) {
	t.Parallel()

	sys := newTestSys(symlinkedProjectReferenceInput(false), false)
	result := execute.CommandLine(sys, []string{"-b", "packages/app/tsconfig.json"}, sys)
	if result.Status == tsc.ExitStatusSuccess {
		t.Fatalf("expected build without paths mapping to fail, got success\noutput:\n%s", sys.getOutput(true))
	}
	output := sys.getOutput(true)
	if !strings.Contains(output, "TS2883") {
		t.Fatalf("expected portability diagnostic without paths mapping, got:\n%s", output)
	}
}
