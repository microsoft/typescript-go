package projectv2

import (
	"maps"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type configFileRegistry struct {
	configs map[tspath.Path]*configFileEntry
}

func (c *configFileRegistry) getConfig(path tspath.Path) *tsoptions.ParsedCommandLine {
	if entry, ok := c.configs[path]; ok {
		return entry.commandLine
	}
	return nil
}

// clone creates a shallow copy of the configFileRegistry.
// The map is cloned, but the configFileEntry values are not.
// Use a configFileRegistryBuilder to create a clone with changes.
func (c *configFileRegistry) clone() *configFileRegistry {
	newConfigs := maps.Clone(c.configs)
	return &configFileRegistry{
		configs: newConfigs,
	}
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
}

var _ tsoptions.ParseConfigHost = (*configFileRegistryBuilder)(nil)

// configFileRegistryBuilder tracks changes made on top of a previous
// configFileRegistry, producing a new clone with `Finalize()` after
// all changes have been made. It is complicated by the fact that project
// loading (and therefore config file parsing/loading) can happen concurrently,
// so the dirty map is a SyncMap.
type configFileRegistryBuilder struct {
	snapshot *Snapshot
	base     *configFileRegistry
	dirty    collections.SyncMap[tspath.Path, *configFileEntry]
}

func newConfigFileRegistryBuilder(newSnapshot *Snapshot, oldConfigFileRegistry *configFileRegistry) *configFileRegistryBuilder {
	return &configFileRegistryBuilder{
		snapshot: newSnapshot,
		base:     oldConfigFileRegistry,
	}
}

// finalize creates a new configFileRegistry based on the changes made in the builder.
// If no changes were made, it returns the original base registry.
func (c *configFileRegistryBuilder) finalize() *configFileRegistry {
	var changed bool
	newRegistry := c.base
	c.dirty.Range(func(key tspath.Path, entry *configFileEntry) bool {
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
		if dirty, ok := c.dirty.Load(path); ok {
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
		entry, loaded := c.dirty.LoadOrStore(path, &configFileEntry{
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
	if entry, ok := c.dirty.Load(path); ok {
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

// acquireConfig loads a config file entry from the cache, or parses it if not already
// cached, then adds the project (if provided) to `retainingProjects` to keep it alive
// in the cache. Each `acquireConfig` call that passes a `project` should be accompanied
// by an eventual `releaseConfig` call with the same project.
func (c *configFileRegistryBuilder) acquireConfig(fileName string, path tspath.Path, project *Project) *tsoptions.ParsedCommandLine {
	entry, _ := c.loadOrStoreNewEntry(path)

	if project != nil {
		entry.retainProject(project.configFilePath)
	}

	// !!! move into single locked method
	switch entry.pendingReload {
	case PendingReloadFileNames:
		entry.setCommandLine(tsoptions.ReloadFileNamesOfParsedCommandLine(entry.commandLine, c.snapshot.compilerFS))
	case PendingReloadFull:
		// oldCommandLine := entry.commandLine
		// !!! extended config cache
		newCommandLine, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c, nil)
		entry.setCommandLine(newCommandLine)
		// release oldCommandLine extended configs
	}

	return entry.commandLine
}

// releaseConfig removes the project from the config entry. Once no projects are
// associated with the config entry, it will be removed on the next call to `cleanup`.
func (c *configFileRegistryBuilder) releaseConfig(path tspath.Path, project *Project) {
	if entry, ok := c.load(path); ok {
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry.releaseProject(project.configFilePath)
	}
}

// FS implements tsoptions.ParseConfigHost.
func (c *configFileRegistryBuilder) FS() vfs.FS {
	return c.snapshot.compilerFS
}

// GetCurrentDirectory implements tsoptions.ParseConfigHost.
func (c *configFileRegistryBuilder) GetCurrentDirectory() string {
	return c.snapshot.sessionOptions.CurrentDirectory
}

// configFileBuilderEntry is a wrapper around `configFileEntry` that
// stores whether the underlying entry was found in the dirty map
// (i.e., it is already a clone and can be mutated) or whether it
// came from the previous configFileRegistry (in which case it must
// be cloned into the dirty map when changes are made). Each setter
// method checks this condition and either mutates the already-dirty
// clone or adds a clone into the builder's dirty map.
type configFileBuilderEntry struct {
	b *configFileRegistryBuilder
	*configFileEntry
	key   tspath.Path
	dirty bool
}

// retainProject adds a project to the set of retaining projects.
// configFileEntries will be retained as long as the set of retaining
// projects is non-empty.
func (e *configFileBuilderEntry) retainProject(projectPath tspath.Path) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		e.retainingProjects[projectPath] = struct{}{}
	} else {
		entry := &configFileEntry{
			pendingReload: e.pendingReload,
			commandLine:   e.commandLine,
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirty.LoadOrStore(e.key, entry)
		entry.retainingProjects = maps.Clone(e.retainingProjects)
		entry.retainingProjects[projectPath] = struct{}{}
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
			pendingReload: e.pendingReload,
			commandLine:   e.commandLine,
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirty.LoadOrStore(e.key, entry)
		entry.retainingProjects = maps.Clone(e.retainingProjects)
		delete(entry.retainingProjects, projectPath)
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
			commandLine:       e.commandLine,
			retainingProjects: maps.Clone(e.retainingProjects),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirty.LoadOrStore(e.key, entry)
		entry.pendingReload = reload
		e.configFileEntry = entry
		e.dirty = true
	}
}

func (e *configFileBuilderEntry) setCommandLine(commandLine *tsoptions.ParsedCommandLine) {
	if e.dirty {
		e.mu.Lock()
		defer e.mu.Unlock()
		e.commandLine = commandLine
	} else {
		entry := &configFileEntry{
			pendingReload:     e.pendingReload,
			retainingProjects: maps.Clone(e.retainingProjects),
		}
		entry.mu.Lock()
		defer entry.mu.Unlock()
		entry, _ = e.b.dirty.LoadOrStore(e.key, entry)
		entry.commandLine = commandLine
		e.configFileEntry = entry
		e.dirty = true
	}
}
