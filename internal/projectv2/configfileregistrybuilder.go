package projectv2

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

var _ tsoptions.ParseConfigHost = (*configFileRegistryBuilder)(nil)
var _ tsoptions.ExtendedConfigCache = (*configFileRegistryBuilder)(nil)

// configFileRegistryBuilder tracks changes made on top of a previous
// configFileRegistry, producing a new clone with `finalize()` after
// all changes have been made.
type configFileRegistryBuilder struct {
	fs                  *overlayFS
	extendedConfigCache *extendedConfigCache
	sessionOptions      *SessionOptions

	base            *ConfigFileRegistry
	configs         *dirtySyncMap[tspath.Path, *configFileEntry]
	configFileNames *dirtyMap[tspath.Path, *configFileNames]
}

func newConfigFileRegistryBuilder(
	fs *overlayFS,
	oldConfigFileRegistry *ConfigFileRegistry,
	extendedConfigCache *extendedConfigCache,
	sessionOptions *SessionOptions,
) *configFileRegistryBuilder {
	return &configFileRegistryBuilder{
		fs:                  fs,
		base:                oldConfigFileRegistry,
		sessionOptions:      sessionOptions,
		extendedConfigCache: extendedConfigCache,

		configs: newDirtySyncMap(oldConfigFileRegistry.configs, func(dirty *configFileEntry, original *configFileEntry) *configFileEntry {
			if dirty.retainingProjects == nil && original != nil {
				dirty.retainingProjects = original.retainingProjects
			}
			if dirty.retainingOpenFiles == nil && original != nil {
				dirty.retainingOpenFiles = original.retainingOpenFiles
			}
			return dirty
		}),
		configFileNames: newDirtyMap(oldConfigFileRegistry.configFileNames),
	}
}

// finalize creates a new configFileRegistry based on the changes made in the builder.
// If no changes were made, it returns the original base registry.
func (c *configFileRegistryBuilder) finalize() *ConfigFileRegistry {
	var changed bool
	newRegistry := c.base
	ensureCloned := func() {
		if !changed {
			newRegistry = newRegistry.clone()
			changed = true
		}
	}

	if configs, changedConfigs := c.configs.Finalize(); changedConfigs {
		ensureCloned()
		newRegistry.configs = configs
	}

	if configFileNames, changedNames := c.configFileNames.Finalize(); changedNames {
		ensureCloned()
		newRegistry.configFileNames = configFileNames
	}

	return newRegistry
}

func (c *configFileRegistryBuilder) findOrAcquireConfigForOpenFile(
	configFileName string,
	configFilePath tspath.Path,
	openFilePath tspath.Path,
	loadKind projectLoadKind,
) *tsoptions.ParsedCommandLine {
	switch loadKind {
	case projectLoadKindFind:
		if entry, ok := c.configs.Load(configFilePath); ok {
			return entry.value.commandLine
		}
		return nil
	case projectLoadKindCreate:
		return c.acquireConfigForOpenFile(configFileName, configFilePath, openFilePath)
	default:
		panic(fmt.Sprintf("unknown project load kind: %d", loadKind))
	}
}

// reloadIfNeeded updates the command line of the config file entry based on its
// pending reload state. This function should only be called from within the
// Change() method of a dirty map entry.
func (c *configFileRegistryBuilder) reloadIfNeeded(entry *configFileEntry, fileName string, path tspath.Path) {
	switch entry.pendingReload {
	case PendingReloadFileNames:
		entry.commandLine = tsoptions.ReloadFileNamesOfParsedCommandLine(entry.commandLine, c.fs.fs)
	case PendingReloadFull:
		newCommandLine, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c, c)
		entry.commandLine = newCommandLine
		// !!! release oldCommandLine extended configs on accepting new snapshot
	default:
		return
	}
	entry.pendingReload = PendingReloadNone
}

// acquireConfigForProject loads a config file entry from the cache, or parses it if not already
// cached, then adds the project (if provided) to `retainingProjects` to keep it alive
// in the cache. Each `acquireConfigForProject` call that passes a `project` should be accompanied
// by an eventual `releaseConfigForProject` call with the same project.
func (c *configFileRegistryBuilder) acquireConfigForProject(fileName string, path tspath.Path, project *Project) *tsoptions.ParsedCommandLine {
	entry, _ := c.configs.LoadOrStore(path, &configFileEntry{pendingReload: PendingReloadFull})
	var needsRetainProject bool
	entry = entry.ChangeIf(
		func(config *configFileEntry) bool {
			_, alreadyRetaining := config.retainingProjects[project.configFilePath]
			needsRetainProject = !alreadyRetaining
			return needsRetainProject || config.pendingReload != PendingReloadNone
		},
		func(config *configFileEntry) {
			if needsRetainProject {
				config.retainingProjects = cloneMapIfNil(config, entry.original, func(e *configFileEntry) map[tspath.Path]struct{} {
					return e.retainingProjects
				})
				config.retainingProjects[project.configFilePath] = struct{}{}
			}
			c.reloadIfNeeded(config, fileName, path)
		},
	)
	return entry.value.commandLine
}

// acquireConfigForOpenFile loads a config file entry from the cache, or parses it if not already
// cached, then adds the open file to `retainingOpenFiles` to keep it alive in the cache.
// Each `acquireConfigForOpenFile` call that passes an `openFilePath`
// should be accompanied by an eventual `releaseConfigForOpenFile` call with the same open file.
func (c *configFileRegistryBuilder) acquireConfigForOpenFile(configFileName string, configFilePath tspath.Path, openFilePath tspath.Path) *tsoptions.ParsedCommandLine {
	entry, _ := c.configs.LoadOrStore(configFilePath, &configFileEntry{pendingReload: PendingReloadFull})
	var needsRetainOpenFile bool
	entry = entry.ChangeIf(
		func(config *configFileEntry) bool {
			_, alreadyRetaining := config.retainingOpenFiles[openFilePath]
			needsRetainOpenFile = !alreadyRetaining
			return needsRetainOpenFile || config.pendingReload != PendingReloadNone
		},
		func(config *configFileEntry) {
			if needsRetainOpenFile {
				config.retainingOpenFiles = cloneMapIfNil(config, entry.original, func(e *configFileEntry) map[tspath.Path]struct{} {
					return e.retainingOpenFiles
				})
				config.retainingOpenFiles[openFilePath] = struct{}{}
			}
			c.reloadIfNeeded(config, configFileName, configFilePath)
		},
	)
	return entry.value.commandLine
}

// releaseConfigForProject removes the project from the config entry. Once no projects
// or files are associated with the config entry, it will be removed on the next call to `cleanup`.
// func (c *configFileRegistryBuilder) releaseConfigForProject(path tspath.Path, project *Project) {
// 	if entry, ok := c.load(path); ok {
// 		entry.mu.Lock()
// 		defer entry.mu.Unlock()
// 		entry.releaseProject(project.configFilePath)
// 	}
// }

// releaseConfigsForOpenFile removes the open file from the config entry. Once no projects
// or files are associated with the config entry, it will be removed on the next call to `cleanup`.
func (c *configFileRegistryBuilder) releaseConfigsForOpenFile(openFilePath tspath.Path) {
	c.configs.Range(func(entry *dirtySyncMapEntry[tspath.Path, *configFileEntry]) bool {
		entry.ChangeIf(
			func(config *configFileEntry) bool {
				_, ok := config.retainingOpenFiles[openFilePath]
				return ok
			},
			func(config *configFileEntry) {
				delete(config.retainingOpenFiles, openFilePath)
			},
		)
		return true
	})

	// !!! remove from configFileNames
}

func (c *configFileRegistryBuilder) computeConfigFileName(fileName string, skipSearchInDirectoryOfFile bool) string {
	searchPath := tspath.GetDirectoryPath(fileName)
	result, _ := tspath.ForEachAncestorDirectory(searchPath, func(directory string) (result string, stop bool) {
		tsconfigPath := tspath.CombinePaths(directory, "tsconfig.json")
		if !skipSearchInDirectoryOfFile && c.FS().FileExists(tsconfigPath) {
			return tsconfigPath, true
		}
		jsconfigPath := tspath.CombinePaths(directory, "jsconfig.json")
		if !skipSearchInDirectoryOfFile && c.FS().FileExists(jsconfigPath) {
			return jsconfigPath, true
		}
		if strings.HasSuffix(directory, "/node_modules") {
			return "", true
		}
		skipSearchInDirectoryOfFile = false
		return "", false
	})
	// !!! c.snapshot.Logf("computeConfigFileName:: File: %s:: Result: %s", fileName, result)
	return result
}

func (c *configFileRegistryBuilder) getConfigFileNameForFile(fileName string, path tspath.Path, loadKind projectLoadKind) string {
	if project.IsDynamicFileName(fileName) {
		return ""
	}

	if entry, ok := c.configFileNames.Get(path); ok {
		return entry.value.nearestConfigFileName
	}

	if loadKind == projectLoadKindFind {
		return ""
	}

	configName := c.computeConfigFileName(fileName, false)

	if _, ok := c.fs.overlays[path]; ok {
		c.configFileNames.Add(path, &configFileNames{
			nearestConfigFileName: configName,
		})
	}
	return configName
}

func (c *configFileRegistryBuilder) getAncestorConfigFileName(fileName string, path tspath.Path, configFileName string, loadKind projectLoadKind) string {
	if project.IsDynamicFileName(fileName) {
		return ""
	}

	entry, ok := c.configFileNames.Get(path)
	if !ok {
		return ""
	}
	if ancestorConfigName, found := entry.value.ancestors[configFileName]; found {
		return ancestorConfigName
	}

	if loadKind == projectLoadKindFind {
		return ""
	}

	// Look for config in parent folders of config file
	result := c.computeConfigFileName(configFileName, true)

	if _, ok := c.fs.overlays[path]; ok {
		entry.Change(func(value *configFileNames) {
			if value.ancestors == nil {
				value.ancestors = make(map[string]string)
			}
			value.ancestors[configFileName] = result
		})
	}
	return result
}

// FS implements tsoptions.ParseConfigHost.
func (c *configFileRegistryBuilder) FS() vfs.FS {
	return c.fs.fs
}

// GetCurrentDirectory implements tsoptions.ParseConfigHost.
func (c *configFileRegistryBuilder) GetCurrentDirectory() string {
	return c.sessionOptions.CurrentDirectory
}

// GetExtendedConfig implements tsoptions.ExtendedConfigCache.
func (c *configFileRegistryBuilder) GetExtendedConfig(fileName string, path tspath.Path, parse func() *tsoptions.ExtendedConfigCacheEntry) *tsoptions.ExtendedConfigCacheEntry {
	fh := c.fs.getFile(fileName)
	return c.extendedConfigCache.acquire(fh, path, parse)
}
