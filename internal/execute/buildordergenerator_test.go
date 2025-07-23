package execute_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestBuildOrderGenerator(t *testing.T) {
	t.Parallel()
	testCases := []*buildOrderTestCase{
		{"specify two roots", []string{"A", "G"}, []string{"D", "E", "C", "B", "A", "G"}, false},
		{"multiple parts of the same graph in various orders", []string{"A"}, []string{"D", "E", "C", "B", "A"}, false},
		{"multiple parts of the same graph in various orders", []string{"A", "C", "D"}, []string{"D", "E", "C", "B", "A"}, false},
		{"multiple parts of the same graph in various orders", []string{"D", "C", "A"}, []string{"D", "E", "C", "B", "A"}, false},
		{"other orderings", []string{"F"}, []string{"E", "F"}, false},
		{"other orderings", []string{"E"}, []string{"E"}, false},
		{"other orderings", []string{"F", "C", "A"}, []string{"E", "F", "D", "C", "B", "A"}, false},
		{"returns circular order", []string{"H"}, []string{"E", "J", "I", "H"}, true},
		{"returns circular order", []string{"A", "H"}, []string{"D", "E", "C", "B", "A", "J", "I", "H"}, true},
	}
	for _, testcase := range testCases {
		testcase.run(t)
	}
}

type buildOrderTestCase struct {
	name     string
	projects []string
	expected []string
	circular bool
}

func (b *buildOrderTestCase) configName(project string) string {
	return fmt.Sprintf("/home/src/workspaces/project/%s/tsconfig.json", project)
}

func (b *buildOrderTestCase) projectName(config string) string {
	str := strings.TrimPrefix(config, "/home/src/workspaces/project/")
	str = strings.TrimSuffix(str, "/tsconfig.json")
	return str
}

func (b *buildOrderTestCase) run(t *testing.T) {
	t.Helper()
	t.Run(b.name+" - "+strings.Join(b.projects, ","), func(t *testing.T) {
		t.Parallel()
		files := make(map[string]string)
		deps := map[string][]string{
			"A": {"B", "C"},
			"B": {"C", "D"},
			"C": {"D", "E"},
			"F": {"E"},
			"H": {"I"},
			"I": {"J"},
			"J": {"H", "E"},
		}
		reverseDeps := map[string][]string{}
		for project, deps := range deps {
			for _, dep := range deps {
				reverseDeps[dep] = append(reverseDeps[dep], project)
			}
		}
		for _, project := range []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"} {
			files[fmt.Sprintf("/home/src/workspaces/project/%s/%s.ts", project, project)] = "export {}"
			referencesStr := ""
			if deps, ok := deps[project]; ok {
				referencesStr = fmt.Sprintf(`, "references": [%s]`, strings.Join(core.Map(deps, func(dep string) string {
					return fmt.Sprintf(`{ "path": "../%s" }`, dep)
				}), ","))
			}
			files[b.configName(project)] = fmt.Sprintf(`{
                "compilerOptions": { "composite": true },
                "files": ["./%s.ts"],
                %s
            }`, project, referencesStr)
		}

		host := compiler.NewCompilerHost("/home/src/workspaces/project", vfstest.FromMap(files, true), "", nil)
		args := append([]string{"--build", "--dry"}, b.projects...)
		buildCommand := tsoptions.ParseBuildCommandLine(args, host)
		buildOrderGenerator := execute.NewBuildOrderGenerator(buildCommand, host, false)
		buildOrder := core.Map(buildOrderGenerator.Order(), b.projectName)
		assert.DeepEqual(t, buildOrder, b.expected)

		for index, project := range buildOrder {
			upstream := core.Map(buildOrderGenerator.Upstream(b.configName(project)), b.projectName)
			expectedUpstream := deps[project]
			assert.Assert(t, len(upstream) <= len(expectedUpstream), fmt.Sprintf("Expected upstream for %s to be at most %d, got %d", project, len(expectedUpstream), len(upstream)))
			for _, expected := range expectedUpstream {
				if slices.Contains(buildOrder[:index], expected) {
					assert.Assert(t, slices.Contains(upstream, expected), fmt.Sprintf("Expected upstream for %s to contain %s", project, expected))
				} else {
					assert.Assert(t, !slices.Contains(upstream, expected), fmt.Sprintf("Expected upstream for %s to not contain %s", project, expected))
				}
			}

			downstream := core.Map(buildOrderGenerator.Downstream(b.configName(project)), b.projectName)
			expectedDownstream := reverseDeps[project]
			assert.Assert(t, len(downstream) <= len(expectedDownstream), fmt.Sprintf("Expected downstream for %s to be at most %d, got %d", project, len(expectedDownstream), len(downstream)))
			for _, expected := range expectedDownstream {
				if slices.Contains(buildOrder[index+1:], expected) {
					assert.Assert(t, slices.Contains(downstream, expected), fmt.Sprintf("Expected downstream for %s to contain %s", project, expected))
				} else {
					assert.Assert(t, !slices.Contains(downstream, expected), fmt.Sprintf("Expected downstream for %s to not contain %s", project, expected))
				}
			}
		}

		if !b.circular {
			for project, projectDeps := range deps {
				child := b.configName(project)
				childIndex := slices.Index(buildOrder, child)
				if childIndex == -1 {
					continue
				}
				for _, dep := range projectDeps {
					parent := b.configName(dep)
					parentIndex := slices.Index(buildOrder, parent)

					assert.Assert(t, childIndex > parentIndex, fmt.Sprintf("Expecting child %s to be built after parent %s", project, dep))
				}
			}
		}
	})
}
