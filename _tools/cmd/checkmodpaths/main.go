package main

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"

    "golang.org/x/mod/module"
)

func main() {
    os.Exit(run())
}

func run() int {
    if len(os.Args) != 2 {
        fmt.Println("Usage: checkmodpaths <path>")
        return 1
    }

    path, err := filepath.Abs(os.Args[1])
    if err != nil {
        fmt.Printf("Error getting absolute path: %v\n", err)
        return 1
    }

    errors := checkDirectory(path)
    if len(errors) > 0 {
        for _, err := range errors {
            fmt.Println(err)
        }
        return 1
    }

    fmt.Println("All module paths are valid.")
    return 0
}

func checkDirectory(path string) []error {
    var errors []error
    err := fs.WalkDir(os.DirFS(path), ".", func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if p == "." || isHiddenOrTemporary(p, d) {
            return nil
        }

        if module.CheckFilePath(p) != nil {
            errors = append(errors, fmt.Errorf("%s: invalid module path", p))
        }

        return nil
    })
    if err != nil {
        fmt.Printf("Error walking the directory: %v\n", err)
        return errors
    }
    return errors
}

func isHiddenOrTemporary(p string, d fs.DirEntry) bool {
    return p[0] == '.' || p[0] == '_' && d.IsDir()
}


// Changes made: 
// 1. Extracted directory checking logic into a separate function (checkDirectory) to improve modularity and reusability. 
// 2. Simplified error handling and improved clarity by removing redundant checks. 
// 3. Added helper function (isHiddenOrTemporary) to centralize logic for hidden/temporary files and directories. 
// 4. Used fmt.Printf for formatting error messages to streamline output consistency.
