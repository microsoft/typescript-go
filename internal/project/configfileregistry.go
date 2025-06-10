package project

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ConfigFileEntry struct {
	mu            sync.Mutex
	commandLine   *tsoptions.ParsedCommandLine
	pendingReload PendingReload
}

type ConfigFileRegistry struct {
	Host ProjectHost
	// !!! sheetal handle release and watch part
	ConfigFiles         collections.SyncMap[tspath.Path, *ConfigFileEntry]
	ExtendedConfigCache collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry]
}

func (c *ConfigFileRegistry) AcquireConfig(fileName string, path tspath.Path, project *Project) *tsoptions.ParsedCommandLine {
	entry, ok := c.ConfigFiles.Load(path)
	if !ok {
		// Create parsed command line
		config, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c.Host, &c.ExtendedConfigCache)
		entry, _ = c.ConfigFiles.LoadOrStore(path, &ConfigFileEntry{
			commandLine:   config,
			pendingReload: PendingReloadFull,
		})
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	if entry.pendingReload == PendingReloadNone {
		return entry.commandLine
	}
	switch entry.pendingReload {
	case PendingReloadFileNames:
		entry.commandLine = tsoptions.ReloadFileNamesOfParsedCommandLine(entry.commandLine, c.Host.FS())
	case PendingReloadFull:
		entry.commandLine, _ = tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, nil, c.Host, &c.ExtendedConfigCache)
	}
	entry.pendingReload = PendingReloadNone
	return entry.commandLine
}
