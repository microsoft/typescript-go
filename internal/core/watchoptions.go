package core

type WatchOptions struct {
	FileKind        WatchFileKind      `json:"watchFile"`
	DirectoryKind   WatchDirectoryKind `json:"watchDirectory"`
	FallbackPolling PollingKind        `json:"fallbackPolling"`
	SyncWatchDir    Tristate           `json:"synchronousWatchDirectory"`
	ExcludeDir      []string           `json:"excludeDirectories"`
	ExcludeFiles    []string           `json:"excludeFiles"`
}

type WatchFileKind int32

const (
	WatchFileKind_None                         WatchFileKind = 0
	WatchFileKind_FixedPollingInterval         WatchFileKind = 1
	WatchFileKind_PriorityPollingInterval      WatchFileKind = 2
	WatchFileKind_DynamicPriorityPolling       WatchFileKind = 3
	WatchFileKind_FixedChunkSizePolling        WatchFileKind = 4
	WatchFileKind_UseFsEvents                  WatchFileKind = 5
	WatchFileKind_UseFsEventsOnParentDirectory WatchFileKind = 6
)

type WatchDirectoryKind int32

const (
	WatchDirectoryKind_None                   WatchDirectoryKind = 0
	WatchDirectoryKind_UseFsEvents            WatchDirectoryKind = 1
	WatchDirectoryKind_FixedPollingInterval   WatchDirectoryKind = 2
	WatchDirectoryKind_DynamicPriorityPolling WatchDirectoryKind = 3
	WatchDirectoryKind_FixedChunkSizePolling  WatchDirectoryKind = 4
)

type PollingKind int32

const (
	PollingKind_None             PollingKind = 0
	PollingKind_FixedInterval    PollingKind = 1
	PollingKind_PriorityInterval PollingKind = 2
	PollingKind_DynamicPriority  PollingKind = 3
	PollingKind_FixedChunkSize   PollingKind = 4
)
