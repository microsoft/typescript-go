package ata

// Logger is an interface for logging messages during the typings installation process.
type Logger interface {
	Log(msg ...any)
}
