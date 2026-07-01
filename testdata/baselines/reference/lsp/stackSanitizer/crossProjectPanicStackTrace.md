Test name: `TestSanitizedCrossProjectPanicStackTrace`

# Unsanitized input:

````
github.com/microsoft/typescript-go/internal/checker.(*Checker).checkExpression(0xc0045a8000, {0x10f6688, 0xc00c2871d0}, 0xc0001fe008, 0x0)
	github.com/microsoft/typescript-go/internal/checker/checker.go:5000 +0x1a0
github.com/microsoft/typescript-go/internal/ls.(*LanguageService).provideSymbolsAndEntries(0xc008329200, {0x10f6688, 0xc00c2871d0}, {0xc00b472030, 0x28}, {0x2, 0x4}, 0x0, 0x0)
	github.com/microsoft/typescript-go/internal/ls/findallreferences.go:100 +0x200
github.com/microsoft/typescript-go/internal/ls.handleCrossProject[...].func1()
	github.com/microsoft/typescript-go/internal/ls/crossproject.go:105 +0x150
github.com/microsoft/typescript-go/internal/core.(*WorkGroup).worker(0xc000120080)
	github.com/microsoft/typescript-go/internal/core/workgroup.go:50 +0x80
created by github.com/microsoft/typescript-go/internal/core.(*WorkGroup).Queue in goroutine 35
	github.com/microsoft/typescript-go/internal/core/workgroup.go:35 +0x60
````

# Sanitized output:

````
typescript-go|>internal|>checker.(*Checker).checkExpression()
	typescript-go|>internal|>checker|>checker.go:5000
typescript-go|>internal|>ls.(*LanguageService).provideSymbolsAndEntries()
	typescript-go|>internal|>ls|>findallreferences.go:100
typescript-go|>internal|>ls.handleCrossProject[...].func1()
	typescript-go|>internal|>ls|>crossproject.go:105
typescript-go|>internal|>core.(*WorkGroup).worker()
	typescript-go|>internal|>core|>workgroup.go:50
typescript-go|>internal|>core.(*WorkGroup).Queue
	typescript-go|>internal|>core|>workgroup.go:35
````
