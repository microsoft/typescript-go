package projectv2

import (
	"maps"

	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ConfigFileRegistry struct {
	// configs is a map of config file paths to their entries.
	configs map[tspath.Path]*configFileEntry
	// configFileNames is a map of open file paths to information
	// about their ancestor config file names. It is only used as
	// a cache during
	configFileNames map[tspath.Path]*configFileNames
}

type configFileEntry struct {
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

// Clone creates a shallow copy of the configFileEntry, without maps.
// A nil map is used in the builder to indicate that a dirty entry still
// shares the same map as its original. During finalization, nil maps
// should be replaced with the maps from the original entry.
func (e *configFileEntry) Clone() *configFileEntry {
	return &configFileEntry{
		pendingReload: e.pendingReload,
		commandLine:   e.commandLine,
	}
}

func (c *ConfigFileRegistry) GetConfig(path tspath.Path) *tsoptions.ParsedCommandLine {
	if entry, ok := c.configs[path]; ok {
		return entry.commandLine
	}
	return nil
}

func (c *ConfigFileRegistry) GetConfigFileName(path tspath.Path) string {
	if entry, ok := c.configFileNames[path]; ok {
		return entry.nearestConfigFileName
	}
	return ""
}

func (c *ConfigFileRegistry) GetAncestorConfigFileName(path tspath.Path, higherThanConfig string) string {
	if entry, ok := c.configFileNames[path]; ok {
		return entry.ancestors[higherThanConfig]
	}
	return ""
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

func (c *configFileNames) Clone() *configFileNames {
	return &configFileNames{
		nearestConfigFileName: c.nearestConfigFileName,
		ancestors:             maps.Clone(c.ancestors),
	}
}
