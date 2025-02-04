package tsoptions

import (
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
)

var optionsForWatch = []*CommandLineOption{
	{
		Name: "watchFile",
		Kind: CommandLineOptionTypeEnum,
		// new Map(Object.entries({
		//     fixedpollinginterval: WatchFileKind.FixedPollingInterval,
		//     prioritypollinginterval: WatchFileKind.PriorityPollingInterval,
		//     dynamicprioritypolling: WatchFileKind.DynamicPriorityPolling,
		//     fixedchunksizepolling: WatchFileKind.FixedChunkSizePolling,
		//     usefsevents: WatchFileKind.UseFsEvents,
		//     usefseventsonparentdirectory: WatchFileKind.UseFsEventsOnParentDirectory,
		// })),
		category:                diagnostics.Watch_and_Build_Modes,
		description:             diagnostics.Specify_how_the_TypeScript_watch_mode_works,
		defaultValueDescription: core.WatchFileKind_UseFsEvents,
	},
	{
		Name: "watchDirectory",
		Kind: CommandLineOptionTypeEnum,
		// new Map(Object.entries({
		//     usefsevents: WatchDirectoryKind.UseFsEvents,
		//     fixedpollinginterval: WatchDirectoryKind.FixedPollingInterval,
		//     dynamicprioritypolling: WatchDirectoryKind.DynamicPriorityPolling,
		//     fixedchunksizepolling: WatchDirectoryKind.FixedChunkSizePolling,
		// })),
		category:                diagnostics.Watch_and_Build_Modes,
		description:             diagnostics.Specify_how_directories_are_watched_on_systems_that_lack_recursive_file_watching_functionality,
		defaultValueDescription: core.WatchDirectoryKind_UseFsEvents,
	},
	{
		Name: "fallbackPolling",
		Kind: CommandLineOptionTypeEnum,
		// new Map(Object.entries({
		//     fixedinterval: PollingWatchKind.FixedInterval,
		//     priorityinterval: PollingWatchKind.PriorityInterval,
		//     dynamicpriority: PollingWatchKind.DynamicPriority,
		//     fixedchunksize: PollingWatchKind.FixedChunkSize,
		// })),
		category:                diagnostics.Watch_and_Build_Modes,
		description:             diagnostics.Specify_what_approach_the_watcher_should_use_if_the_system_runs_out_of_native_file_watchers,
		defaultValueDescription: core.PollingKind_PriorityInterval,
	},
	{
		Name:                    "synchronousWatchDirectory",
		Kind:                    CommandLineOptionTypeBoolean,
		category:                diagnostics.Watch_and_Build_Modes,
		description:             diagnostics.Synchronously_call_callbacks_and_update_the_state_of_directory_watchers_on_platforms_that_don_t_support_recursive_watching_natively,
		defaultValueDescription: false,
	},
	{
		Name: "excludeDirectories",
		Kind: CommandLineOptionTypeList,
		// element: {
		//     Name: "excludeDirectory",
		//     Kind: "string",
		//     isFilePath: true,
		//     extraValidation: specToDiagnostic,
		// },
		allowConfigDirTemplateSubstitution: true,
		category:                           diagnostics.Watch_and_Build_Modes,
		description:                        diagnostics.Remove_a_list_of_directories_from_the_watch_process,
	},
	{
		Name: "excludeFiles",
		Kind: CommandLineOptionTypeList,
		// element: {
		//     Name: "excludeFile",
		//     Kind: "string",
		//     isFilePath: true,
		//     extraValidation: specToDiagnostic,
		// },
		allowConfigDirTemplateSubstitution: true,
		category:                           diagnostics.Watch_and_Build_Modes,
		description:                        diagnostics.Remove_a_list_of_files_from_the_watch_mode_s_processing,
	},
}
