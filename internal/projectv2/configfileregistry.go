package projectv2

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type configFileEntry struct {
	mu sync.RWMutex

	pendingReload PendingReload
	commandLine   *tsoptions.ParsedCommandLine
	// retainingProjects is the set of projects that have called acquireConfig
	// without releasing it. A config file entry may be acquired by a project
	// either because it is the config for that project or because it is the
	// config for a referenced project.
	retainingProjects collections.Set[tspath.Path]
}

// setPendingReload sets the reload level if it is higher than the current level.
// Returns whether the level was changed.
func (e *configFileEntry) setPendingReload(level PendingReload) bool {
	if e.pendingReload < level {
		e.pendingReload = level
		return true
	}
	return false
}

type configFileRegistry struct {
	snapshot *Snapshot

	mu      sync.RWMutex
	configs map[tspath.Path]*configFileEntry
}

var _ tsoptions.ParseConfigHost = (*configFileRegistry)(nil)

// acquireConfig loads a config file entry from the cache, or parses it if not already
// cached, then adds the project (if provided) to `retainingProjects` to keep it alive
// in the cache. Each `acquireConfig` call that passes a `project` should be accompanied
// by an eventual `releaseConfig` call with the same project.
func (c *configFileRegistry) acquireConfig(fileName string, path tspath.Path, project *Project) *tsoptions.ParsedCommandLine {
	c.mu.Lock()
	entry, ok := c.configs[path]
	if !ok {
		entry = &configFileEntry{
			commandLine:   nil,
			pendingReload: PendingReloadFull,
		}
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	c.mu.Unlock()

	if project != nil {
		entry.retainingProjects.Add(project.configFilePath)
	}

	switch entry.pendingReload {
	case PendingReloadFileNames:
		entry.commandLine = tsoptions.ReloadFileNamesOfParsedCommandLine(entry.commandLine, c.snapshot.compilerFS)
	case PendingReloadFull:
		// oldCommandLine := entry.commandLine
		// !!! extended config cache
		entry.commandLine, _ = tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c, nil)
		// release oldCommandLine extended configs
	}

	return entry.commandLine
}

// releaseConfig removes the project from the config entry. Once no projects are
// associated with the config entry, it will be removed on the next call to `cleanup`.
func (c *configFileRegistry) releaseConfig(path tspath.Path, project *Project) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, ok := c.configs[path]; ok {
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry.retainingProjects.Delete(project.configFilePath)
	}
}

func (c *configFileRegistry) getConfig(path tspath.Path) *tsoptions.ParsedCommandLine {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, ok := c.configs[path]; ok {
		entry.mu.RLock()
		defer entry.mu.RUnlock()
		return entry.commandLine
	}
	return nil
}

// FS implements tsoptions.ParseConfigHost.
func (c *configFileRegistry) FS() vfs.FS {
	return c.snapshot.compilerFS
}

// GetCurrentDirectory implements tsoptions.ParseConfigHost.
func (c *configFileRegistry) GetCurrentDirectory() string {
	return c.snapshot.sessionOptions.CurrentDirectory
}
