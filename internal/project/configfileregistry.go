package project

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ConfigFileEntry struct {
	mu             sync.RWMutex
	commandLine    *tsoptions.ParsedCommandLine
	projects       collections.Set[*Project]
	infos          collections.Set[*ScriptInfo]
	pendingReload  PendingReload
	rootFilesWatch *watchedFiles[[]string]
}

type ExtendedConfigFileEntry struct {
	mu          sync.Mutex
	configFiles collections.Set[tspath.Path]
}

type parseStatus struct {
	wg *sync.WaitGroup
}

type ConfigFileRegistry struct {
	Host                  ProjectHost
	defaultProjectFinder  *defaultProjectFinder
	ConfigFiles           collections.SyncMap[tspath.Path, *ConfigFileEntry]
	ExtendedConfigCache   collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry]
	ExtendedConfigsUsedBy collections.SyncMap[tspath.Path, *ExtendedConfigFileEntry]

	parseMu       sync.Mutex
	activeParsing map[tspath.Path]*parseStatus
}

func NewConfigFileRegistry(registry *ConfigFileRegistry) *ConfigFileRegistry {
	registry.activeParsing = make(map[tspath.Path]*parseStatus)
	return registry
}

func (e *ConfigFileEntry) SetPendingReload(level PendingReload) bool {
	if e.pendingReload < level {
		e.pendingReload = level
		return true
	}
	return false
}

var _ watchFileHost = (*configFileWatchHost)(nil)

type configFileWatchHost struct {
	fileName string
	host     ProjectHost
}

func (h *configFileWatchHost) Name() string {
	return h.fileName
}

func (c *configFileWatchHost) Client() Client {
	return c.host.Client()
}

func (c *configFileWatchHost) Log(message string) {
	c.host.Log(message)
}

func (c *ConfigFileRegistry) releaseConfig(path tspath.Path, project *Project) {
	entry, ok := c.ConfigFiles.Load(path)
	if !ok {
		return
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	entry.projects.Delete(project)
}

func (c *ConfigFileRegistry) acquireConfig(fileName string, path tspath.Path, project *Project, info *ScriptInfo) *tsoptions.ParsedCommandLine {
	if entry, ok := c.ConfigFiles.Load(path); ok {
		entry.mu.RLock()
		needsReload := entry.pendingReload != PendingReloadNone
		commandLine := entry.commandLine
		entry.mu.RUnlock()

		if !needsReload && commandLine != nil {
			entry.mu.Lock()
			if project != nil {
				entry.projects.Add(project)
			} else if info != nil {
				entry.infos.Add(info)
			}
			entry.mu.Unlock()
			return commandLine
		}
	}

	c.parseMu.Lock()
	if status, ok := c.activeParsing[path]; ok {
		c.parseMu.Unlock()
		status.wg.Wait()

		if entry, ok := c.ConfigFiles.Load(path); ok {
			entry.mu.Lock()
			defer entry.mu.Unlock()
			if project != nil {
				entry.projects.Add(project)
			} else if info != nil {
				entry.infos.Add(info)
			}
			return entry.commandLine
		}
	}

	status := &parseStatus{wg: &sync.WaitGroup{}}
	status.wg.Add(1)
	c.activeParsing[path] = status
	c.parseMu.Unlock()

	defer func() {
		status.wg.Done()
		c.parseMu.Lock()
		delete(c.activeParsing, path)
		c.parseMu.Unlock()
	}()

	entry, loaded := c.ConfigFiles.Load(path)
	if !loaded {
		config, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c.Host, &c.ExtendedConfigCache)
		var rootFilesWatch *watchedFiles[[]string]
		client := c.Host.Client()
		if c.Host.IsWatchEnabled() && client != nil {
			rootFilesWatch = newWatchedFiles(&configFileWatchHost{fileName: fileName, host: c.Host},
				lsproto.WatchKindChange|lsproto.WatchKindCreate|lsproto.WatchKindDelete,
				core.Identity, "root files")
		}

		newEntry := &ConfigFileEntry{
			commandLine:    config,
			pendingReload:  PendingReloadNone, // Already parsed, no reload needed
			rootFilesWatch: rootFilesWatch,
		}

		entry, loaded = c.ConfigFiles.LoadOrStore(path, newEntry)
		if !loaded {
			c.updateRootFilesWatch(fileName, entry)
			c.updateExtendedConfigsUsedBy(path, entry, nil)
		}
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if project != nil {
		entry.projects.Add(project)
	} else if info != nil {
		entry.infos.Add(info)
	}

	if entry.pendingReload != PendingReloadNone {
		switch entry.pendingReload {
		case PendingReloadFileNames:
			entry.commandLine = entry.commandLine.ReloadFileNamesOfParsedCommandLine(c.Host.FS())
		case PendingReloadFull:
			oldCommandLine := entry.commandLine
			entry.commandLine, _ = tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c.Host, &c.ExtendedConfigCache)
			c.updateExtendedConfigsUsedBy(path, entry, oldCommandLine)
			c.updateRootFilesWatch(fileName, entry)
		}
		entry.pendingReload = PendingReloadNone
	}

	return entry.commandLine
}

func (c *ConfigFileRegistry) getConfig(path tspath.Path) *tsoptions.ParsedCommandLine {
	entry, ok := c.ConfigFiles.Load(path)
	if ok {
		entry.mu.RLock()
		defer entry.mu.RUnlock()
		return entry.commandLine
	}
	return nil
}

func (c *ConfigFileRegistry) releaseConfigsForInfo(info *ScriptInfo) {
	c.ConfigFiles.Range(func(path tspath.Path, entry *ConfigFileEntry) bool {
		entry.mu.Lock()
		entry.infos.Delete(info)
		entry.mu.Unlock()
		return true
	})
}

func (c *ConfigFileRegistry) updateRootFilesWatch(fileName string, entry *ConfigFileEntry) {
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
	entry.rootFilesWatch.update(context.Background(), rootFileGlobs)
}

func (c *ConfigFileRegistry) updateExtendedConfigsUsedBy(path tspath.Path, entry *ConfigFileEntry, oldCommandLine *tsoptions.ParsedCommandLine) {
	extendedConfigs := entry.commandLine.ExtendedSourceFiles()
	newConfigs := make([]tspath.Path, 0, len(extendedConfigs))
	for _, extendedConfig := range extendedConfigs {
		extendedPath := tspath.ToPath(extendedConfig, c.Host.GetCurrentDirectory(), c.Host.FS().UseCaseSensitiveFileNames())
		newConfigs = append(newConfigs, extendedPath)
		extendedEntry, _ := c.ExtendedConfigsUsedBy.LoadOrStore(extendedPath, &ExtendedConfigFileEntry{
			mu: sync.Mutex{},
		})
		extendedEntry.mu.Lock()
		extendedEntry.configFiles.Add(path)
		extendedEntry.mu.Unlock()
	}
	if oldCommandLine != nil {
		for _, extendedConfig := range oldCommandLine.ExtendedSourceFiles() {
			extendedPath := tspath.ToPath(extendedConfig, c.Host.GetCurrentDirectory(), c.Host.FS().UseCaseSensitiveFileNames())
			if !slices.Contains(newConfigs, extendedPath) {
				extendedEntry, _ := c.ExtendedConfigsUsedBy.Load(extendedPath)
				extendedEntry.mu.Lock()
				extendedEntry.configFiles.Delete(path)
				if extendedEntry.configFiles.Len() == 0 {
					c.ExtendedConfigsUsedBy.Delete(extendedPath)
					c.ExtendedConfigCache.Delete(extendedPath)
				}
				extendedEntry.mu.Unlock()
			}
		}
	}
}

func (c *ConfigFileRegistry) onWatchedFilesChanged(path tspath.Path, changeKind lsproto.FileChangeType) (err error, handled bool) {
	if c.onConfigChange(path, changeKind) {
		handled = true
	}

	if entry, loaded := c.ExtendedConfigsUsedBy.Load(path); loaded {
		entry.mu.Lock()
		for configFilePath := range entry.configFiles.Keys() {
			if c.onConfigChange(configFilePath, changeKind) {
				handled = true
			}
		}
		entry.mu.Unlock()
	}
	return err, handled
}

func (c *ConfigFileRegistry) onConfigChange(path tspath.Path, changeKind lsproto.FileChangeType) bool {
	entry, ok := c.ConfigFiles.Load(path)
	if !ok {
		return false
	}
	entry.mu.Lock()
	hasSet := entry.SetPendingReload(PendingReloadFull)
	var infos map[*ScriptInfo]struct{}
	var projects map[*Project]struct{}
	if hasSet {
		infos = maps.Clone(entry.infos.Keys())
		projects = maps.Clone(entry.projects.Keys())
	}
	entry.mu.Unlock()
	if !hasSet {
		return false
	}
	for info := range infos {
		delete(c.defaultProjectFinder.configFileForOpenFiles, info.Path())
		delete(c.defaultProjectFinder.configFilesAncestorForOpenFiles, info.Path())
	}
	for project := range projects {
		if project.configFilePath == path {
			switch changeKind {
			case lsproto.FileChangeTypeCreated:
				fallthrough
			case lsproto.FileChangeTypeChanged:
				project.deferredClose = false
				project.SetPendingReload(PendingReloadFull)
			case lsproto.FileChangeTypeDeleted:
				project.deferredClose = true
			}
		} else {
			project.markAsDirty()
		}
	}
	return true
}

func (c *ConfigFileRegistry) tryInvokeWildCardDirectories(fileName string, path tspath.Path) {
	configFiles := c.ConfigFiles.ToMap()
	for configPath, entry := range configFiles {
		entry.mu.Lock()
		hasSet := false
		if entry.commandLine != nil && entry.pendingReload == PendingReloadNone && entry.commandLine.MatchesFileName(fileName) {
			hasSet = entry.SetPendingReload(PendingReloadFileNames)
		}
		var projects map[*Project]struct{}
		if hasSet {
			projects = maps.Clone(entry.projects.Keys())
		}
		entry.mu.Unlock()
		if hasSet {
			for project := range projects {
				if project.configFilePath == configPath {
					project.SetPendingReload(PendingReloadFileNames)
				} else {
					project.markAsDirty()
				}
			}
		}
	}
}

func (c *ConfigFileRegistry) cleanup(toRemoveConfigs map[tspath.Path]*ConfigFileEntry) {
	for path, entry := range toRemoveConfigs {
		entry.mu.Lock()
		if entry.projects.Len() == 0 && entry.infos.Len() == 0 {
			c.ConfigFiles.Delete(path)
			commandLine := entry.commandLine
			entry.commandLine = nil
			c.updateExtendedConfigsUsedBy(path, entry, commandLine)
			if entry.rootFilesWatch != nil {
				entry.rootFilesWatch.update(context.Background(), nil)
			}
		}
		entry.mu.Unlock()
	}
}
