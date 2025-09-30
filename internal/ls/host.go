package ls

type Host interface {
	UseCaseSensitiveFileNames() bool
	ReadFile(path string) (contents string, ok bool)
	FileExists(path string) bool
	Converters() *Converters
}
