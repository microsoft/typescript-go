package projectv2

import (
	"fmt"
	"maps"
	"sync"

	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ConfigFileRegistry struct {
	// configs is a map of config file paths to their entries.
	configs map[tspath.Path]*configFileEntry
	// configFileNames is a map of open file paths to information
	// about their ancestor config file names. It is only used as
	// a cache during
	configFileNames map[tspath.Path]configFileNames
}

type configFileEntry struct {
	// mu needs only be held by configFileRegistryBuilder methods,
	// as configFileEntries are considered immutable once they move
	// from the builder to the finalized registry.
	mu            sync.Mutex
	pendingReload PendingReload
	commandLine   *tsoptions.ParsedCommandLine
	// retainingProjects is the set of projects that have called acquireConfig
	// without releasing it. A config file entry may be acquired by a project
	// either because it is the config for that project or because it is the
	// config for a referenced project.
	retainingProjects map[tspath.Path]struct{}
	// retainingOpenFiles is the set of open files that caused this config to
	// load during project collection building. This config file may or may not
	// end up being the config for the default project for these files, but
	// determining the default project loaded this config as a candidate, so
	// subsequent calls to `projectCollectionBuilder.findDefaultConfiguredProject`
	// will use this config as part of the search, so it must be retained.
	retainingOpenFiles map[tspath.Path]struct{}
}

func (c *ConfigFileRegistry) GetConfig(path tspath.Path, project *Project) *tsoptions.ParsedCommandLine {
	if entry, ok := c.configs[path]; ok {
		if _, ok := entry.retainingProjects[project.configFilePath]; !ok {
			panic(fmt.Sprintf("project %s should have called acquireConfig for config file %s during registry building", project.Name, path))
		}
		return entry.commandLine
	}
	return nil
}

// clone creates a shallow copy of the configFileRegistry.
func (c *ConfigFileRegistry) clone() *ConfigFileRegistry {
	return &ConfigFileRegistry{
		configs:         maps.Clone(c.configs),
		configFileNames: maps.Clone(c.configFileNames),
	}
}

type configFileNames struct {
	// nearestConfigFileName is the file name of the nearest ancestor config file.
	nearestConfigFileName string
	// ancestors is a map from one ancestor config file path to the next.
	// For example, if `/a`, `/a/b`, and `/a/b/c` all contain config files,
	// the fully loaded map will look like:
	//		{
	//			"/a/b/c/tsconfig.json": "/a/b/tsconfig.json",
	//			"/a/b/tsconfig.json": "/a/tsconfig.json"
	//		}
	ancestors map[string]string
}
