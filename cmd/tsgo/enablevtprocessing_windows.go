package main

import (
    "golang.org/x/sys/windows"
)

func enableVirtualTerminalProcessing() {
    h, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
    if err != nil || h == windows.InvalidHandle {
        return
    }
    var mode uint32
    if err := windows.GetConsoleMode(h, &mode); err == nil && mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING == 0 {
        _ = windows.SetConsoleMode(h, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
    }
}


// Changes made: 
// 1. Removed unnecessary call to windows.GetFileType, as it was not being used effectively in the logic. 
// 2. Combined conditions in GetConsoleMode to streamline error handling and mode checking. 
// 3. Simplified the function by eliminating redundant steps to improve readability and maintainability.
