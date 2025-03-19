package printer

type ResolveModuleNameResolutionHost interface {
	GetCurrentDirectory() string
	CommonSourceDirectory() string
	UseCaseSensitiveFileNames() bool
	ShouldTransformImportCall() bool
}
