package projectv2

import (
	"fmt"
	"maps"
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// configFileNamesEntry tracks changes to a `configFileNames` entry. When a change is requested
// on one of the underlying maps, it clones the map and adds the entry to the configFileRegistryBuilder's
// map of dirty configFileNames.
type configFileNamesEntry struct {
	configFileNames
	c     *configFileRegistryBuilder
	key   tspath.Path
	dirty bool
}

func (b *configFileNamesEntry) setConfigFileName(fileName string) {
	b.nearestConfigFileName = fileName
	if !b.dirty {
		if b.c.dirtyConfigFileNames == nil {
			b.c.dirtyConfigFileNames = make(map[tspath.Path]configFileNames)
		}
		b.c.dirtyConfigFileNames[b.key] = b.configFileNames
		b.dirty = true
	}
}

func (b *configFileNamesEntry) addAncestorConfigFileName(configFileName string, ancestorConfigFileName string) {
	if !b.dirty {
		b.ancestors = maps.Clone(b.ancestors)
		if b.c.dirtyConfigFileNames == nil {
			b.c.dirtyConfigFileNames = make(map[tspath.Path]configFileNames)
		}
		b.c.dirtyConfigFileNames[b.key] = b.configFileNames
		b.dirty = true
	}
	if b.ancestors == nil {
		b.ancestors = make(map[string]string)
	}
	b.ancestors[configFileName] = ancestorConfigFileName
}

var _ tsoptions.ParseConfigHost = (*configFileRegistryBuilder)(nil)
var _ tsoptions.ExtendedConfigCache = (*configFileRegistryBuilder)(nil)

// configFileRegistryBuilder tracks changes made on top of a previous
// configFileRegistry, producing a new clone with `finalize()` after
// all changes have been made.
type configFileRegistryBuilder struct {
	snapshot             *Snapshot
	base                 *ConfigFileRegistry
	dirtyConfigs         collections.SyncMap[tspath.Path, *configFileEntry]
	dirtyConfigFileNames map[tspath.Path]configFileNames
}

func newConfigFileRegistryBuilder(newSnapshot *Snapshot, oldConfigFileRegistry *ConfigFileRegistry) *configFileRegistryBuilder {
	return &configFileRegistryBuilder{
		snapshot: newSnapshot,
		base:     oldConfigFileRegistry,
	}
}

// finalize creates a new configFileRegistry based on the changes made in the builder.
// If no changes were made, it returns the original base registry.
func (c *configFileRegistryBuilder) finalize() *ConfigFileRegistry {
	var changed bool
	newRegistry := c.base
	c.dirtyConfigs.Range(func(key tspath.Path, entry *configFileEntry) bool {
		if !changed {
			newRegistry = newRegistry.clone()
			if newRegistry.configs == nil {
				newRegistry.configs = make(map[tspath.Path]*configFileEntry)
			}
			changed = true
		}
		newRegistry.configs[key] = entry
		return true
	})
	if len(c.dirtyConfigFileNames) > 0 {
		if !changed {
			newRegistry = newRegistry.clone()
		}
		if newRegistry.configFileNames == nil {
			newRegistry.configFileNames = make(map[tspath.Path]configFileNames)
		} else {
			newRegistry.configFileNames = maps.Clone(newRegistry.configFileNames)
		}
		for key, names := range c.dirtyConfigFileNames {
			if _, ok := newRegistry.configFileNames[key]; !ok {
				newRegistry.configFileNames[key] = names
			} else {
				// If the key already exists, we merge the names.
				existingNames := newRegistry.configFileNames[key]
				existingNames.nearestConfigFileName = names.nearestConfigFileName
				maps.Copy(existingNames.ancestors, names.ancestors)
				newRegistry.configFileNames[key] = existingNames
			}
		}
	}
	return newRegistry
}

// loadOrStoreNewEntry looks up the config file entry or creates a new one,
// returning the entry, whether it was loaded (as opposed to created),
// *and* whether the entry is in the dirty map.
func (c *configFileRegistryBuilder) loadOrStoreNewEntry(path tspath.Path) (entry *configFileBuilderEntry, loaded bool) {
	// Check for existence in the base registry first so that all SyncMap
	// access is atomic. We're trying to avoid the scenario where we
	//   1. try to load from the dirty map but find nothing,
	//   2. try to load from the base registry but find nothing, then
	//   3. have to do a subsequent Store in the dirty map for the new entry.
	if prev, ok := c.base.configs[path]; ok {
		if dirty, ok := c.dirtyConfigs.Load(path); ok {
			return &configFileBuilderEntry{
				b:               c,
				key:             path,
				configFileEntry: dirty,
				dirty:           true,
			}, true
		}
		return &configFileBuilderEntry{
			b:               c,
			key:             path,
			configFileEntry: prev,
			dirty:           false,
		}, true
	} else {
		entry, loaded := c.dirtyConfigs.LoadOrStore(path, &configFileEntry{
			pendingReload: PendingReloadFull,
		})
		return &configFileBuilderEntry{
			b:               c,
			key:             path,
			configFileEntry: entry,
			dirty:           true,
		}, loaded
	}
}

func (c *configFileRegistryBuilder) load(path tspath.Path) (*configFileBuilderEntry, bool) {
	if entry, ok := c.dirtyConfigs.Load(path); ok {
		return &configFileBuilderEntry{
			b:               c,
			key:             path,
			configFileEntry: entry,
			dirty:           true,
		}, true
	}
	if entry, ok := c.base.configs[path]; ok {
		return &configFileBuilderEntry{
			b:               c,
			key:             path,
			configFileEntry: entry,
			dirty:           false,
		}, true
	}
	return nil, false
}

func (c *configFileRegistryBuilder) getConfigFileNames(path tspath.Path) *configFileNamesEntry {
	names, inDirty := c.dirtyConfigFileNames[path]
	if !inDirty {
		names, _ = c.base.configFileNames[path]
	}
	return &configFileNamesEntry{
		c:               c,
		configFileNames: names,
		dirty:           inDirty,
	}
}

func (c *configFileRegistryBuilder) findOrAcquireConfigForOpenFile(
	configFileName string,
	configFilePath tspath.Path,
	openFilePath tspath.Path,
	loadKind projectLoadKind,
) *tsoptions.ParsedCommandLine {
	switch loadKind {
	case projectLoadKindFind:
		if config, ok := c.load(configFilePath); ok {
			return config.commandLine
		}
		return nil
	case projectLoadKindCreate:
		return c.acquireConfigForOpenFile(configFileName, configFilePath, openFilePath)
	default:
		panic(fmt.Sprintf("unknown project load kind: %d", loadKind))
	}
}

// acquireConfigForProject loads a config file entry from the cache, or parses it if not already
// cached, then adds the project (if provided) to `retainingProjects` to keep it alive
// in the cache. Each `acquireConfigForProject` call that passes a `project` should be accompanied
// by an eventual `releaseConfigForProject` call with the same project.
func (c *configFileRegistryBuilder) acquireConfigForProject(fileName string, path tspath.Path, project *Project) *tsoptions.ParsedCommandLine {
	entry, _ := c.loadOrStoreNewEntry(path)
	entry.retainProject(project.configFilePath)
	entry.reloadIfNeeded(fileName, path)
	return entry.commandLine
}

// acquireConfigForOpenFile loads a config file entry from the cache, or parses it if not already
// cached, then adds the open file to `retainingOpenFiles` to keep it alive in the cache.
// Each `acquireConfigForOpenFile` call that passes an `openFilePath`
// should be accompanied by an eventual `releaseConfigForOpenFile` call with the same open file.
func (c *configFileRegistryBuilder) acquireConfigForOpenFile(configFileName string, configFilePath tspath.Path, openFilePath tspath.Path) *tsoptions.ParsedCommandLine {
	entry, _ := c.loadOrStoreNewEntry(configFilePath)
	entry.retainOpenFile(openFilePath)
	entry.reloadIfNeeded(configFileName, configFilePath)
	return entry.commandLine
}

// releaseConfigForProject removes the project from the config entry. Once no projects are
// associated with the config entry, it will be removed on the next call to `cleanup`.
func (c *configFileRegistryBuilder) releaseConfigForProject(path tspath.Path, project *Project) {
	if entry, ok := c.load(path); ok {
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry.releaseProject(project.configFilePath)
	}
}

// releaseConfigForOpenFile removes the project from the config entry. Once no projects are
// associated with the config entry, it will be removed on the next call to `cleanup`.
func (c *configFileRegistryBuilder) releaseConfigForOpenFile(path tspath.Path, openFilePath tspath.Path) {
	if entry, ok := c.load(path); ok {
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry.releaseOpenFile(openFilePath)
	}
}

func (c *configFileRegistryBuilder) computeConfigFileName(fileName string, skipSearchInDirectoryOfFile bool) string {
	searchPath := tspath.GetDirectoryPath(fileName)
	result, _ := tspath.ForEachAncestorDirectory(searchPath, func(directory string) (result string, stop bool) {
		tsconfigPath := tspath.CombinePaths(directory, "tsconfig.json")
		if !skipSearchInDirectoryOfFile && c.snapshot.compilerFS.FileExists(tsconfigPath) {
			return tsconfigPath, true
		}
		jsconfigPath := tspath.CombinePaths(directory, "jsconfig.json")
		if !skipSearchInDirectoryOfFile && c.snapshot.compilerFS.FileExists(jsconfigPath) {
			return jsconfigPath, true
		}
		if strings.HasSuffix(directory, "/node_modules") {
			return "", true
		}
		skipSearchInDirectoryOfFile = false
		return "", false
	})
	c.snapshot.Logf("computeConfigFileName:: File: %s:: Result: %s", fileName, result)
	return result
}

func (c *configFileRegistryBuilder) getConfigFileNameForFile(fileName string, path tspath.Path, loadKind projectLoadKind) string {
	if project.IsDynamicFileName(fileName) {
		return ""
	}

	configFileNames := c.getConfigFileNames(path)
	if configFileNames.nearestConfigFileName != "" {
		return configFileNames.nearestConfigFileName
	}

	if loadKind == projectLoadKindFind {
		return ""
	}

	configName := c.computeConfigFileName(fileName, false)

	if c.snapshot.IsOpenFile(path) {
		configFileNames.setConfigFileName(configName)
	}
	return configName
}

func (c *configFileRegistryBuilder) getAncestorConfigFileName(fileName string, path tspath.Path, configFileName string, loadKind projectLoadKind) string {
	if project.IsDynamicFileName(fileName) {
		return ""
	}

	configFileNames := c.getConfigFileNames(path)
	if ancestorConfigName, found := configFileNames.ancestors[configFileName]; found {
		return ancestorConfigName
	}

	if loadKind == projectLoadKindFind {
		return ""
	}

	// Look for config in parent folders of config file
	result := c.computeConfigFileName(configFileName, true)

	if c.snapshot.IsOpenFile(path) {
		configFileNames.addAncestorConfigFileName(configFileName, result)
	}
	return result
}

// FS implements tsoptions.ParseConfigHost.
func (c *configFileRegistryBuilder) FS() vfs.FS {
	return c.snapshot.compilerFS
}

// GetCurrentDirectory implements tsoptions.ParseConfigHost.
func (c *configFileRegistryBuilder) GetCurrentDirectory() string {
	return c.snapshot.sessionOptions.CurrentDirectory
}

// GetExtendedConfig implements tsoptions.ExtendedConfigCache.
func (c *configFileRegistryBuilder) GetExtendedConfig(fileName string, path tspath.Path, parse func() *tsoptions.ExtendedConfigCacheEntry) *tsoptions.ExtendedConfigCacheEntry {
	fh := c.snapshot.GetFile(ls.FileNameToDocumentURI(fileName))
	return c.snapshot.extendedConfigCache.acquire(fh, path, parse)
}

// configFileBuilderEntry is a wrapper around `configFileEntry` that
// stores whether the underlying entry was found in the dirty map
// (i.e., it is already a clone and can be mutated) or whether it
// came from the previous configFileRegistry (in which case it must
// be cloned into the dirty map when changes are made). Each setter
// method checks this condition and either mutates the already-dirty
// clone or adds a clone into the builder's dirty map.
type configFileBuilderEntry struct {
	*configFileEntry
	b     *configFileRegistryBuilder
	key   tspath.Path
	dirty bool
}

// retainProject adds a project to the set of retaining projects.
// configFileEntries will be retained as long as the set of retaining
// projects and retaining open files are non-empty.
func (e *configFileBuilderEntry) retainProject(projectPath tspath.Path) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		if e.retainingProjects == nil {
			e.retainingProjects = make(map[tspath.Path]struct{})
		}
		e.retainingProjects[projectPath] = struct{}{}
	} else {
		entry := &configFileEntry{
			pendingReload:      e.pendingReload,
			commandLine:        e.commandLine,
			retainingOpenFiles: maps.Clone(e.retainingOpenFiles),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirtyConfigs.LoadOrStore(e.key, entry)
		entry.retainingProjects = maps.Clone(e.retainingProjects)
		entry.retainingProjects[projectPath] = struct{}{}
		e.configFileEntry = entry
		e.dirty = true
	}
}

// retainOpenFile adds an open file to the set of retaining open files.
// configFileEntries will be retained as long as the set of retaining
// projects and retaining open files are non-empty.
func (e *configFileBuilderEntry) retainOpenFile(openFilePath tspath.Path) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		if e.retainingOpenFiles == nil {
			e.retainingOpenFiles = make(map[tspath.Path]struct{})
		}
		e.retainingOpenFiles[openFilePath] = struct{}{}
	} else {
		entry := &configFileEntry{
			pendingReload:     e.pendingReload,
			commandLine:       e.commandLine,
			retainingProjects: maps.Clone(e.retainingProjects),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirtyConfigs.LoadOrStore(e.key, entry)
		entry.retainingOpenFiles = maps.Clone(e.retainingOpenFiles)
		entry.retainingOpenFiles[openFilePath] = struct{}{}
		e.configFileEntry = entry
		e.dirty = true
	}
}

// releaseProject removes a project from the set of retaining projects.
func (e *configFileBuilderEntry) releaseProject(projectPath tspath.Path) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		delete(e.retainingProjects, projectPath)
	} else {
		entry := &configFileEntry{
			pendingReload:      e.pendingReload,
			commandLine:        e.commandLine,
			retainingOpenFiles: maps.Clone(e.retainingOpenFiles),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirtyConfigs.LoadOrStore(e.key, entry)
		entry.retainingProjects = maps.Clone(e.retainingProjects)
		delete(entry.retainingProjects, projectPath)
		e.configFileEntry = entry
		e.dirty = true
	}
}

// releaseOpenFile removes an open file from the set of retaining open files.
func (e *configFileBuilderEntry) releaseOpenFile(openFilePath tspath.Path) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		delete(e.retainingOpenFiles, openFilePath)
	} else {
		entry := &configFileEntry{
			pendingReload:     e.pendingReload,
			commandLine:       e.commandLine,
			retainingProjects: maps.Clone(e.retainingProjects),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirtyConfigs.LoadOrStore(e.key, entry)
		entry.retainingOpenFiles = maps.Clone(e.retainingOpenFiles)
		delete(entry.retainingOpenFiles, openFilePath)
		e.configFileEntry = entry
		e.dirty = true
	}
}

func (e *configFileBuilderEntry) setPendingReload(reload PendingReload) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		e.pendingReload = reload
	} else {
		entry := &configFileEntry{
			commandLine:        e.commandLine,
			retainingProjects:  maps.Clone(e.retainingProjects),
			retainingOpenFiles: maps.Clone(e.retainingOpenFiles),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirtyConfigs.LoadOrStore(e.key, entry)
		entry.pendingReload = reload
		e.configFileEntry = entry
		e.dirty = true
	}
}

func (e *configFileBuilderEntry) reloadIfNeeded(fileName string, path tspath.Path) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		if e.pendingReload == PendingReloadNone {
			return
		}
	} else {
		if e.pendingReload == PendingReloadNone {
			return
		}
		entry := &configFileEntry{
			pendingReload:      e.pendingReload,
			retainingProjects:  maps.Clone(e.retainingProjects),
			retainingOpenFiles: maps.Clone(e.retainingOpenFiles),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirtyConfigs.LoadOrStore(e.key, entry)
		e.configFileEntry = entry
		e.dirty = true
	}

	switch e.pendingReload {
	case PendingReloadFileNames:
		e.commandLine = tsoptions.ReloadFileNamesOfParsedCommandLine(e.commandLine, e.b.snapshot.compilerFS)
	case PendingReloadFull:
		newCommandLine, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, e.b, e.b)
		e.commandLine = newCommandLine
		// !!! release oldCommandLine extended configs on accepting new snapshot
	}
	e.pendingReload = PendingReloadNone
}
