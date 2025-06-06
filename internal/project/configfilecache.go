package project

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type configFileCacheEntry struct {
	commandLineMu   sync.RWMutex
	commandLine     *tsoptions.ParsedCommandLine
	projects        collections.SyncMap[tspath.Path, struct{}]
	pendingReload   PendingReload
	rootFilesWatch  *watchedFiles[[]string]
	extendedConfigs []tspath.Path
}

var _ watchFileHost = (*configWatchFileHost)(nil)

type configFileCache struct {
	service                *Service
	configFileCache        collections.SyncMap[tspath.Path, *configFileCacheEntry]
	extendedConfigCache    collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry]
	extendedConfigToConfig collections.SyncMap[tspath.Path, *collections.SyncMap[tspath.Path, struct{}]]
}

type configWatchFileHost struct {
	fileName string
	service  *Service
}

func (h *configWatchFileHost) Name() string {
	return h.fileName
}

func (c *configWatchFileHost) Client() Client {
	return c.service.host.Client()
}

func (c *configWatchFileHost) Log(message string) {
	c.service.Log(message)
}

func (c *configFileCache) getResolvedProjectReference(fileName string, path tspath.Path, forProject tspath.Path) *tsoptions.ParsedCommandLine {
	entry, ok := c.configFileCache.Load(path)
	if !ok {
		// Create parsed command line
		config, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c.service.host, &c.extendedConfigCache)

		var rootFilesWatch *watchedFiles[[]string]
		client := c.service.host.Client()
		if c.service.IsWatchEnabled() && client != nil {
			rootFilesWatch = newWatchedFiles(&configWatchFileHost{fileName: fileName, service: c.service}, lsproto.WatchKindChange|lsproto.WatchKindCreate|lsproto.WatchKindDelete, core.Identity, "root files")
		}
		entry, _ = c.configFileCache.LoadOrStore(path, &configFileCacheEntry{
			commandLine:    config,
			rootFilesWatch: rootFilesWatch,
			pendingReload:  PendingReloadFull,
		})
	}
	entry.projects.Store(forProject, struct{}{})
	commandLine := c.ensureConfigUpToDate(fileName, path, entry)
	return commandLine
}

func (c *configFileCache) ensureConfigUpToDate(fileName string, path tspath.Path, entry *configFileCacheEntry) *tsoptions.ParsedCommandLine {
	entry.commandLineMu.RLock()
	if entry.pendingReload == PendingReloadNone {
		entry.commandLineMu.RUnlock()
		return entry.commandLine
	}
	entry.commandLineMu.RUnlock()
	entry.commandLineMu.Lock()
	defer entry.commandLineMu.Unlock()
	switch entry.pendingReload {
	case PendingReloadFileNames:
		entry.commandLine = tsoptions.ReloadFileNamesOfParsedCommandLine(entry.commandLine, c.service.host.FS())
	case PendingReloadFull:
		entry.commandLine, _ = tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c.service.host, &c.extendedConfigCache)
		c.updateRootFilesWatch(fileName, path, entry)
	}
	entry.pendingReload = PendingReloadNone
	return entry.commandLine
}

func (c *configFileCache) updateRootFilesWatch(fileName string, path tspath.Path, entry *configFileCacheEntry) {
	var extendedConfigs []string
	if entry.commandLine != nil {
		extendedConfigs = entry.commandLine.ConfigFile.ExtendedSourceFiles
	}
	c.updateExtendedConfigUse(path, entry, extendedConfigs)
	if entry.rootFilesWatch == nil {
		return
	}

	wildcardGlobs := entry.commandLine.WildcardDirectories()
	rootFileGlobs := make([]string, 0, len(wildcardGlobs)+1+len(extendedConfigs))
	rootFileGlobs = append(rootFileGlobs, fileName)
	for _, extendedConfig := range extendedConfigs {
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

func (c *configFileCache) updateExtendedConfigUse(path tspath.Path, entry *configFileCacheEntry, extendedConfigs []string) {
	newConfigs := make([]tspath.Path, 0, len(extendedConfigs))
	for _, extendedConfig := range extendedConfigs {
		extendedPath := c.service.toPath(extendedConfig)
		newConfigs = append(newConfigs, extendedPath)
		extendedEntry, _ := c.extendedConfigToConfig.LoadOrStore(extendedPath, &collections.SyncMap[tspath.Path, struct{}]{})
		extendedEntry.Store(path, struct{}{})
	}
	for _, extendedPath := range entry.extendedConfigs {
		if !slices.Contains(newConfigs, extendedPath) {
			extendedEntry, _ := c.extendedConfigToConfig.Load(extendedPath)
			extendedEntry.Delete(path)
			if extendedEntry.Size() == 0 {
				c.extendedConfigToConfig.Delete(extendedPath)
				c.extendedConfigCache.Delete(extendedPath)
			}
		}
	}
	entry.extendedConfigs = newConfigs
}

func (c *configFileCache) releaseResolvedProjectReference(path tspath.Path, forProject tspath.Path) {
	entry, ok := c.configFileCache.Load(path)
	if !ok {
		return
	}
	entry.projects.Delete(forProject)

	// !!! sheetal todo proper handling of config file deletion
	entry.commandLineMu.Lock()
	defer entry.commandLineMu.Unlock()
	if entry.projects.Size() == 0 {
		c.configFileCache.Delete(path)
		c.updateExtendedConfigUse(path, entry, nil)
		if entry.rootFilesWatch != nil {
			entry.rootFilesWatch.update(context.Background(), nil)
			entry.rootFilesWatch = nil
		}
		entry.commandLine = nil
	}
}

func (c *configFileCache) onWatchedFilesChanged(path tspath.Path, changeKind lsproto.FileChangeType) (err error, handled bool) {
	if c.onConfigChange(path, changeKind) {
		handled = true
	}

	// !!! todo wild cards
	if entry, loaded := c.extendedConfigToConfig.Load(path); loaded {
		entry.Range(func(configPath tspath.Path, _ struct{}) bool {
			c.onConfigChange(configPath, changeKind)
			return true
		})
		handled = true
	}
	return err, handled
}

func (c *configFileCache) onConfigChange(path tspath.Path, changeKind lsproto.FileChangeType) bool {
	entry, ok := c.configFileCache.Load(path)
	if !ok {
		return false
	}

	entry.commandLineMu.Lock()
	defer entry.commandLineMu.Unlock()
	entry.pendingReload = PendingReloadFull
	entry.commandLine = nil
	entry.projects.Range(func(projectPath tspath.Path, _ struct{}) bool {
		project := c.service.ConfiguredProject(projectPath)
		if project == nil {
			return true
		}
		if projectPath == path {
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
		return true
	})
	return true
}

func (c *configFileCache) tryInvokeWildCardDirectories(fileName string, path tspath.Path) {
	// !!! sheetal concurrency
	c.configFileCache.Range(func(configPath tspath.Path, entry *configFileCacheEntry) bool {
		if entry.commandLine != nil && entry.commandLine.MatchesFileName(fileName) {
			entry.pendingReload = PendingReloadFileNames
			entry.projects.Range(func(projectPath tspath.Path, _ struct{}) bool {
				project := c.service.ConfiguredProject(projectPath)
				if project != nil {
					if projectPath == configPath {
						project.SetPendingReload(PendingReloadFileNames)
					} else {
						project.markAsDirty()
					}
				}
				return true
			})
		}
		return true
	})
}
