package projectv2

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/dirty"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
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
	fs                  *snapshotFSBuilder
	extendedConfigCache *extendedConfigCache
	sessionOptions      *SessionOptions
	logger              *logCollector

	base            *ConfigFileRegistry
	configs         *dirty.SyncMap[tspath.Path, *configFileEntry]
	configFileNames *dirty.Map[tspath.Path, *configFileNames]
}

func newConfigFileRegistryBuilder(
	fs *snapshotFSBuilder,
	oldConfigFileRegistry *ConfigFileRegistry,
	extendedConfigCache *extendedConfigCache,
	sessionOptions *SessionOptions,
	logger *logCollector,
) *configFileRegistryBuilder {
	return &configFileRegistryBuilder{
		fs:                  fs,
		base:                oldConfigFileRegistry,
		sessionOptions:      sessionOptions,
		extendedConfigCache: extendedConfigCache,
		logger:              logger,

		configs:         dirty.NewSyncMap(oldConfigFileRegistry.configs, nil),
		configFileNames: dirty.NewMap(oldConfigFileRegistry.configFileNames),
	}
}

// Finalize creates a new configFileRegistry based on the changes made in the builder.
// If no changes were made, it returns the original base registry.
func (c *configFileRegistryBuilder) Finalize() *ConfigFileRegistry {
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
			return entry.Value().commandLine
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
		if c.logger != nil {
			c.logger.Log(fmt.Sprintf("Reloading file names for config: %s", fileName))
		}
		entry.commandLine = tsoptions.ReloadFileNamesOfParsedCommandLine(entry.commandLine, c.fs.fs)
	case PendingReloadFull:
		if c.logger != nil {
			c.logger.Log(fmt.Sprintf("Loading config file: %s", fileName))
		}
		entry.commandLine, _ = tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c, c)
		c.updateExtendingConfigs(path, entry.commandLine, entry.commandLine)
		c.updateRootFilesWatch(fileName, entry)
	default:
		return
	}
	entry.pendingReload = PendingReloadNone
}

func (c *configFileRegistryBuilder) updateExtendingConfigs(extendingConfigPath tspath.Path, newCommandLine *tsoptions.ParsedCommandLine, oldCommandLine *tsoptions.ParsedCommandLine) {
	var newExtendedConfigPaths collections.Set[tspath.Path]
	if newCommandLine != nil {
		for _, extendedConfig := range newCommandLine.ExtendedSourceFiles() {
			extendedConfigPath := c.fs.toPath(extendedConfig)
			newExtendedConfigPaths.Add(extendedConfigPath)
			entry, loaded := c.configs.LoadOrStore(extendedConfigPath, newExtendedConfigFileEntry(extendingConfigPath))
			if loaded {
				entry.ChangeIf(
					func(config *configFileEntry) bool {
						_, alreadyRetaining := config.retainingConfigs[extendingConfigPath]
						return !alreadyRetaining
					},
					func(config *configFileEntry) {
						if config.retainingConfigs == nil {
							config.retainingConfigs = make(map[tspath.Path]struct{})
						}
						config.retainingConfigs[extendingConfigPath] = struct{}{}
					},
				)
			}
		}
	}
	if oldCommandLine != nil {
		for _, extendedConfig := range oldCommandLine.ExtendedSourceFiles() {
			extendedConfigPath := c.fs.toPath(extendedConfig)
			if newExtendedConfigPaths.Has(extendedConfigPath) {
				continue
			}
			if entry, ok := c.configs.Load(extendedConfigPath); ok {
				entry.ChangeIf(
					func(config *configFileEntry) bool {
						_, ok := config.retainingConfigs[extendingConfigPath]
						return ok
					},
					func(config *configFileEntry) {
						delete(config.retainingConfigs, extendingConfigPath)
					},
				)
			}
		}
	}
}

func (c *configFileRegistryBuilder) updateRootFilesWatch(fileName string, entry *configFileEntry) {
	if entry.rootFilesWatch == nil {
		return
	}

	wildcardGlobs := entry.commandLine.WildcardDirectories()
	rootFileGlobs := make([]string, 0, len(wildcardGlobs)+1+len(entry.commandLine.ExtendedSourceFiles()))
	rootFileGlobs = append(rootFileGlobs, fileName)
	for _, extendedConfig := range entry.commandLine.ExtendedSourceFiles() {
		rootFileGlobs = append(rootFileGlobs, extendedConfig)
	}
	for dir, recursive := range wildcardGlobs {
		rootFileGlobs = append(rootFileGlobs, fmt.Sprintf("%s/%s", tspath.NormalizePath(dir), core.IfElse(recursive, recursiveFileGlobPattern, fileGlobPattern)))
	}
	for _, fileName := range entry.commandLine.LiteralFileNames() {
		rootFileGlobs = append(rootFileGlobs, fileName)
	}

	slices.Sort(rootFileGlobs)
	entry.rootFilesWatch = entry.rootFilesWatch.Clone(rootFileGlobs)
}

// acquireConfigForProject loads a config file entry from the cache, or parses it if not already
// cached, then adds the project (if provided) to `retainingProjects` to keep it alive
// in the cache. Each `acquireConfigForProject` call that passes a `project` should be accompanied
// by an eventual `releaseConfigForProject` call with the same project.
func (c *configFileRegistryBuilder) acquireConfigForProject(fileName string, path tspath.Path, project *Project) *tsoptions.ParsedCommandLine {
	entry, _ := c.configs.LoadOrStore(path, newConfigFileEntry(fileName))
	var needsRetainProject bool
	entry.ChangeIf(
		func(config *configFileEntry) bool {
			_, alreadyRetaining := config.retainingProjects[project.configFilePath]
			needsRetainProject = !alreadyRetaining
			return needsRetainProject || config.pendingReload != PendingReloadNone
		},
		func(config *configFileEntry) {
			if needsRetainProject {
				if config.retainingProjects == nil {
					config.retainingProjects = make(map[tspath.Path]struct{})
				}
				config.retainingProjects[project.configFilePath] = struct{}{}
			}
			c.reloadIfNeeded(config, fileName, path)
		},
	)
	return entry.Value().commandLine
}

// acquireConfigForOpenFile loads a config file entry from the cache, or parses it if not already
// cached, then adds the open file to `retainingOpenFiles` to keep it alive in the cache.
// Each `acquireConfigForOpenFile` call that passes an `openFilePath`
// should be accompanied by an eventual `releaseConfigForOpenFile` call with the same open file.
func (c *configFileRegistryBuilder) acquireConfigForOpenFile(configFileName string, configFilePath tspath.Path, openFilePath tspath.Path) *tsoptions.ParsedCommandLine {
	entry, _ := c.configs.LoadOrStore(configFilePath, newConfigFileEntry(configFileName))
	var needsRetainOpenFile bool
	entry.ChangeIf(
		func(config *configFileEntry) bool {
			_, alreadyRetaining := config.retainingOpenFiles[openFilePath]
			needsRetainOpenFile = !alreadyRetaining
			return needsRetainOpenFile || config.pendingReload != PendingReloadNone
		},
		func(config *configFileEntry) {
			if needsRetainOpenFile {
				if config.retainingOpenFiles == nil {
					config.retainingOpenFiles = make(map[tspath.Path]struct{})
				}
				config.retainingOpenFiles[openFilePath] = struct{}{}
			}
			c.reloadIfNeeded(config, configFileName, configFilePath)
		},
	)
	return entry.Value().commandLine
}

// releaseConfigForProject removes the project from the config entry. Once no projects
// or files are associated with the config entry, it will be removed on the next call to `cleanup`.
func (c *configFileRegistryBuilder) releaseConfigForProject(configFilePath tspath.Path, projectPath tspath.Path) {
	if entry, ok := c.configs.Load(configFilePath); ok {
		entry.ChangeIf(
			func(config *configFileEntry) bool {
				_, ok := config.retainingProjects[projectPath]
				return ok
			},
			func(config *configFileEntry) {
				delete(config.retainingProjects, projectPath)
			},
		)
	}
}

// DidCloseFile removes the open file from the config entry. Once no projects
// or files are associated with the config entry, it will be removed on the next call to `cleanup`.
func (c *configFileRegistryBuilder) DidCloseFile(path tspath.Path) {
	c.configFileNames.Delete(path)
	c.configs.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *configFileEntry]) bool {
		entry.ChangeIf(
			func(config *configFileEntry) bool {
				_, ok := config.retainingOpenFiles[path]
				return ok
			},
			func(config *configFileEntry) {
				delete(config.retainingOpenFiles, path)
			},
		)
		return true
	})
}

type changeFileResult struct {
	affectedProjects map[tspath.Path]struct{}
	affectedFiles    map[tspath.Path]struct{}
}

func (r changeFileResult) IsEmpty() bool {
	return len(r.affectedProjects) == 0 && len(r.affectedFiles) == 0
}

func (c *configFileRegistryBuilder) DidChangeFile(path tspath.Path) changeFileResult {
	return c.handlePossibleConfigChange(path, lsproto.FileChangeTypeChanged)
}

func (c *configFileRegistryBuilder) DidCreateFile(fileName string, path tspath.Path) changeFileResult {
	result := c.handlePossibleConfigChange(path, lsproto.FileChangeTypeCreated)
	if result.IsEmpty() {
		affectedProjects := c.handlePossibleRootFileCreation(fileName, path)
		if affectedProjects != nil {
			if result.affectedProjects == nil {
				result.affectedProjects = make(map[tspath.Path]struct{})
			}
			maps.Copy(result.affectedProjects, affectedProjects)
		}
	}
	return result
}

func (c *configFileRegistryBuilder) DidDeleteFile(path tspath.Path) changeFileResult {
	return c.handlePossibleConfigChange(path, lsproto.FileChangeTypeDeleted)
}

func (c *configFileRegistryBuilder) handlePossibleConfigChange(path tspath.Path, changeKind lsproto.FileChangeType) changeFileResult {
	var affectedProjects map[tspath.Path]struct{}
	if entry, ok := c.configs.Load(path); ok {
		entry.Locked(func(entry dirty.Value[*configFileEntry]) {
			affectedProjects = c.handleConfigChange(entry)
			for extendingConfigPath := range entry.Value().retainingConfigs {
				if extendingConfigEntry, ok := c.configs.Load(extendingConfigPath); ok {
					if affectedProjects == nil {
						affectedProjects = make(map[tspath.Path]struct{})
					}
					maps.Copy(affectedProjects, c.handleConfigChange(extendingConfigEntry))
				}
			}
		})
	}

	var affectedFiles map[tspath.Path]struct{}
	if changeKind != lsproto.FileChangeTypeChanged {
		baseName := tspath.GetBaseFileName(string(path))
		if baseName == "tsconfig.json" || baseName == "jsconfig.json" {
			directoryPath := path.GetDirectoryPath()
			c.configFileNames.Range(func(entry *dirty.MapEntry[tspath.Path, *configFileNames]) bool {
				if directoryPath.ContainsPath(entry.Key()) {
					if affectedFiles == nil {
						affectedFiles = make(map[tspath.Path]struct{})
					}
					affectedFiles[entry.Key()] = struct{}{}
					entry.Delete()
				}
				return true
			})
		}
	}

	return changeFileResult{
		affectedProjects: affectedProjects,
		affectedFiles:    affectedFiles,
	}
}

func (c *configFileRegistryBuilder) handleConfigChange(entry dirty.Value[*configFileEntry]) map[tspath.Path]struct{} {
	var affectedProjects map[tspath.Path]struct{}
	changed := entry.ChangeIf(
		func(config *configFileEntry) bool { return config.pendingReload != PendingReloadFull },
		func(config *configFileEntry) { config.pendingReload = PendingReloadFull },
	)
	if changed {
		affectedProjects = maps.Clone(entry.Value().retainingProjects)
	}

	return affectedProjects
}

func (c *configFileRegistryBuilder) handlePossibleRootFileCreation(fileName string, path tspath.Path) map[tspath.Path]struct{} {
	var affectedProjects map[tspath.Path]struct{}
	c.configs.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *configFileEntry]) bool {
		entry.ChangeIf(
			func(config *configFileEntry) bool {
				return config.commandLine != nil && config.pendingReload == PendingReloadNone && config.commandLine.MatchesFileName(fileName)
			},
			func(config *configFileEntry) {
				config.pendingReload = PendingReloadFileNames
				if affectedProjects == nil {
					affectedProjects = make(map[tspath.Path]struct{})
				}
				maps.Copy(affectedProjects, config.retainingProjects)
			},
		)
		return true
	})
	return affectedProjects
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
		return entry.Value().nearestConfigFileName
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
	if ancestorConfigName, found := entry.Value().ancestors[configFileName]; found {
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
	fh := c.fs.GetFileByPath(fileName, path)
	return c.extendedConfigCache.Acquire(fh, path, parse)
}

func (c *configFileRegistryBuilder) Cleanup() {
	c.configs.Range(func(entry *dirty.SyncMapEntry[tspath.Path, *configFileEntry]) bool {
		entry.DeleteIf(func(value *configFileEntry) bool {
			return len(value.retainingProjects) == 0 && len(value.retainingOpenFiles) == 0 && len(value.retainingConfigs) == 0
		})
		return true
	})
}
