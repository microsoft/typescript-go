package projectv2

import (
	"cmp"
	"maps"
	"slices"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type ProjectCollection struct {
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
		return cmp.Compare(a.Name, b.Name)
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

func (c *ProjectCollection) GetDefaultProject(fileName string, path tspath.Path) *Project {
	if result, ok := c.fileDefaultProjects[path]; ok {
		if result == "" {
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
	// !!! I'm not sure of a less hacky way to do this without repeating a lot of code.
	panic("TODO")
	// builder := newProjectCollectionBuilder(context.Background(), c.snapshot, c, c.snapshot.configFileRegistry)
	// defer func() {
	// 	c2, r2 := builder.Finalize()
	// 	if c2 != c || r2 != c.snapshot.configFileRegistry {
	// 		panic("temporary builder should have collected no changes for a find lookup")
	// 	}
	// }()

	// if entry := builder.findDefaultConfiguredProject(fileName, path); entry != nil {
	// 	return entry.project
	// }
	// return firstConfiguredProject
}

// clone creates a shallow copy of the project collection, without the
// fileDefaultProjects map.
func (c *ProjectCollection) clone() *ProjectCollection {
	return &ProjectCollection{
		configuredProjects: maps.Clone(c.configuredProjects),
		inferredProject:    c.inferredProject,
	}
}
