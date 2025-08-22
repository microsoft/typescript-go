package incremental

import (
	"time"

	"github.com/microsoft/typescript-go/internal/compiler"
)

type BuildHost interface {
	GetMTime(fileName string) time.Time
	SetMTime(fileName string, mTime time.Time) error
}

type buildHost struct {
	host compiler.CompilerHost
}

var _ BuildHost = (*buildHost)(nil)

func (b *buildHost) GetMTime(fileName string) time.Time {
	return GetMTime(b.host, fileName)
}

func (b *buildHost) SetMTime(fileName string, mTime time.Time) error {
	return SetMTime(b.host, fileName, mTime)
}

func CreateBuildHost(host compiler.CompilerHost) BuildHost {
	return &buildHost{host: host}
}

func GetMTime(host compiler.CompilerHost, fileName string) time.Time {
	stat := host.FS().Stat(fileName)
	var mTime time.Time
	if stat != nil {
		mTime = stat.ModTime()
	}
	return mTime
}

func SetMTime(host compiler.CompilerHost, fileName string, mTime time.Time) error {
	return host.FS().Chtimes(fileName, time.Time{}, mTime)
}
