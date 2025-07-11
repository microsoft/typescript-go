package projectv2

import (
	"cmp"
	"slices"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ProjectCollection struct {
	toPath             func(fileName string) tspath.Path
	configFileRegistry *ConfigFileRegistry
	// fileDefaultProjects is a map of file paths to the config file path (the key
	// into `configuredProjects`) of the default project for that file. If the file
	// belongs to the inferred project, the value is "". This map contains quick
	// lookups for only the associations discovered during the latest snapshot
	// update.
	fileDefaultProjects map[tspath.Path]tspath.Path
	// configuredProjects is the set of loaded projects associated with a tsconfig
	// file, keyed by the config file path.
	configuredProjects map[tspath.Path]*Project
	// inferredProject is a fallback project that is used when no configured
	// project can be found for an open file.
	inferredProject *Project
}

func (c *ProjectCollection) ConfiguredProject(path tspath.Path) *Project {
	return c.configuredProjects[path]
}

func (c *ProjectCollection) ConfiguredProjects() []*Project {
	projects := make([]*Project, 0, len(c.configuredProjects))
	c.fillConfiguredProjects(&projects)
	return projects
}

func (c *ProjectCollection) fillConfiguredProjects(projects *[]*Project) {
	for _, p := range c.configuredProjects {
		*projects = append(*projects, p)
	}
	slices.SortFunc(*projects, func(a, b *Project) int {
		return cmp.Compare(a.Name(), b.Name())
	})
}

func (c *ProjectCollection) Projects() []*Project {
	if c.inferredProject == nil {
		return c.ConfiguredProjects()
	}
	projects := make([]*Project, 0, len(c.configuredProjects)+1)
	c.fillConfiguredProjects(&projects)
	projects = append(projects, c.inferredProject)
	return projects
}

func (c *ProjectCollection) InferredProject() *Project {
	return c.inferredProject
}

// !!! this result could be cached
func (c *ProjectCollection) GetDefaultProject(fileName string, path tspath.Path) *Project {
	if result, ok := c.fileDefaultProjects[path]; ok {
		if result == inferredProjectName {
			return c.inferredProject
		}
		return c.configuredProjects[result]
	}

	var (
		containingProjects                       []*Project
		firstConfiguredProject                   *Project
		firstNonSourceOfProjectReferenceRedirect *Project
		multipleDirectInclusions                 bool
	)
	for _, p := range c.ConfiguredProjects() {
		if p.containsFile(path) {
			containingProjects = append(containingProjects, p)
			if !multipleDirectInclusions && !p.IsSourceFromProjectReference(path) {
				if firstNonSourceOfProjectReferenceRedirect == nil {
					firstNonSourceOfProjectReferenceRedirect = p
				} else {
					multipleDirectInclusions = true
				}
			}
			if firstConfiguredProject == nil {
				firstConfiguredProject = p
			}
		}
	}
	if len(containingProjects) == 1 {
		return containingProjects[0]
	}
	if len(containingProjects) == 0 {
		if c.inferredProject != nil && c.inferredProject.containsFile(path) {
			return c.inferredProject
		}
		return nil
	}
	if !multipleDirectInclusions {
		if firstNonSourceOfProjectReferenceRedirect != nil {
			// Multiple projects include the file, but only one is a direct inclusion.
			return firstNonSourceOfProjectReferenceRedirect
		}
		// Multiple projects include the file, and none are direct inclusions.
		return firstConfiguredProject
	}
	// Multiple projects include the file directly.
	if defaultProject := c.findDefaultConfiguredProject(fileName, path); defaultProject != nil {
		return defaultProject
	}
	return firstConfiguredProject
}

func (c *ProjectCollection) findDefaultConfiguredProject(fileName string, path tspath.Path) *Project {
	if configFileName := c.configFileRegistry.GetConfigFileName(path); configFileName != "" {
		return c.findDefaultConfiguredProjectWorker(fileName, path, configFileName, nil, nil)
	}
	return nil
}

func (c *ProjectCollection) findDefaultConfiguredProjectWorker(fileName string, path tspath.Path, configFileName string, visited *collections.SyncSet[*Project], fallback *Project) *Project {
	configFilePath := c.toPath(configFileName)
	project, ok := c.configuredProjects[configFilePath]
	if !ok {
		return nil
	}
	if visited == nil {
		visited = &collections.SyncSet[*Project]{}
	}

	// Look in the config's project and its references recursively.
	search := core.BreadthFirstSearchParallelEx(
		project,
		func(project *Project) []*Project {
			if project.CommandLine == nil {
				return nil
			}
			return core.Map(project.CommandLine.ResolvedProjectReferencePaths(), func(configFileName string) *Project {
				return c.configuredProjects[c.toPath(configFileName)]
			})
		},
		func(project *Project) (isResult bool, stop bool) {
			if project.containsFile(path) {
				return true, !project.IsSourceFromProjectReference(path)
			}
			return false, false
		},
		core.BreadthFirstSearchOptions[*Project]{
			Visited: visited,
		},
	)

	if search.Stopped {
		// If we found a project that directly contains the file, return it.
		return search.Path[0]
	}
	if len(search.Path) > 0 && fallback == nil {
		// If we found a project that contains the file, but it is a source from
		// a project reference, record it as a fallback.
		fallback = search.Path[0]
	}

	// Look for tsconfig.json files higher up the directory tree and do the same. This handles
	// the common case where a higher-level "solution" tsconfig.json contains all projects in a
	// workspace.
	if config := c.configFileRegistry.GetConfig(path); config != nil && config.CompilerOptions().DisableSolutionSearching.IsTrue() {
		return fallback
	}
	if ancestorConfigName := c.configFileRegistry.GetAncestorConfigFileName(path, configFileName); ancestorConfigName != "" {
		return c.findDefaultConfiguredProjectWorker(fileName, path, ancestorConfigName, visited, fallback)
	}
	return fallback
}

// clone creates a shallow copy of the project collection.
func (c *ProjectCollection) clone() *ProjectCollection {
	return &ProjectCollection{
		toPath:              c.toPath,
		configuredProjects:  c.configuredProjects,
		inferredProject:     c.inferredProject,
		fileDefaultProjects: c.fileDefaultProjects,
	}
}
