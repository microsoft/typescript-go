package core

import "time"

type WatchOptions struct {
	Interval        *int               `json:"watchInterval"`
	FileKind        WatchFileKind      `json:"watchFile"`
	DirectoryKind   WatchDirectoryKind `json:"watchDirectory"`
	FallbackPolling PollingKind        `json:"fallbackPolling"`
	SyncWatchDir    Tristate           `json:"synchronousWatchDirectory"`
	ExcludeDir      []string           `json:"excludeDirectories"`
	ExcludeFiles    []string           `json:"excludeFiles"`
}

type WatchFileKind int32

const (
	WatchFileKindNone                         WatchFileKind = 0
	WatchFileKindFixedPollingInterval         WatchFileKind = 1
	WatchFileKindPriorityPollingInterval      WatchFileKind = 2
	WatchFileKindDynamicPriorityPolling       WatchFileKind = 3
	WatchFileKindFixedChunkSizePolling        WatchFileKind = 4
	WatchFileKindUseFsEvents                  WatchFileKind = 5
	WatchFileKindUseFsEventsOnParentDirectory WatchFileKind = 6
)

type WatchDirectoryKind int32

const (
	WatchDirectoryKindNone                   WatchDirectoryKind = 0
	WatchDirectoryKindUseFsEvents            WatchDirectoryKind = 1
	WatchDirectoryKindFixedPollingInterval   WatchDirectoryKind = 2
	WatchDirectoryKindDynamicPriorityPolling WatchDirectoryKind = 3
	WatchDirectoryKindFixedChunkSizePolling  WatchDirectoryKind = 4
)

type PollingKind int32

const (
	PollingKindNone             PollingKind = 0
	PollingKindFixedInterval    PollingKind = 1
	PollingKindPriorityInterval PollingKind = 2
	PollingKindDynamicPriority  PollingKind = 3
	PollingKindFixedChunkSize   PollingKind = 4
)

func (w *WatchOptions) WatchInterval() time.Duration {
	watchInterval := 1000 * time.Millisecond
	if w != nil && w.Interval != nil {
		watchInterval = time.Duration(*w.Interval) * time.Millisecond
	}
	return watchInterval
}

func (w WatchFileKind) MarshalJSON() ([]byte, error) {
	switch w {
	case WatchFileKindNone:
		return []byte("null"), nil
	case WatchFileKindFixedPollingInterval:
		return []byte(`"fixedpollinginterval"`), nil
	case WatchFileKindPriorityPollingInterval:
		return []byte(`"prioritypollinginterval"`), nil
	case WatchFileKindDynamicPriorityPolling:
		return []byte(`"dynamicprioritypolling"`), nil
	case WatchFileKindFixedChunkSizePolling:
		return []byte(`"fixedchunksizepolling"`), nil
	case WatchFileKindUseFsEvents:
		return []byte(`"usefsevents"`), nil
	case WatchFileKindUseFsEventsOnParentDirectory:
		return []byte(`"usefseventsonparentdirectory"`), nil
	default:
		return []byte("null"), nil
	}
}

func (w WatchDirectoryKind) MarshalJSON() ([]byte, error) {
	switch w {
	case WatchDirectoryKindNone:
		return []byte("null"), nil
	case WatchDirectoryKindUseFsEvents:
		return []byte(`"usefsevents"`), nil
	case WatchDirectoryKindFixedPollingInterval:
		return []byte(`"fixedpollinginterval"`), nil
	case WatchDirectoryKindDynamicPriorityPolling:
		return []byte(`"dynamicprioritypolling"`), nil
	case WatchDirectoryKindFixedChunkSizePolling:
		return []byte(`"fixedchunksizepolling"`), nil
	default:
		return []byte("null"), nil
	}
}

func (p PollingKind) MarshalJSON() ([]byte, error) {
	switch p {
	case PollingKindNone:
		return []byte("null"), nil
	case PollingKindFixedInterval:
		return []byte(`"fixedinterval"`), nil
	case PollingKindPriorityInterval:
		return []byte(`"priorityinterval"`), nil
	case PollingKindDynamicPriority:
		return []byte(`"dynamicpriority"`), nil
	case PollingKindFixedChunkSize:
		return []byte(`"fixedchunksize"`), nil
	default:
		return []byte("null"), nil
	}
}
