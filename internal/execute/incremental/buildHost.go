package incremental

import (
	"time"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

type BuildHost interface {
	GetMTime(fileName string) time.Time
	SetMTime(fileName string, mTime time.Time) error
	OnBuildInfoEmit(config *tsoptions.ParsedCommandLine, buildInfo *BuildInfo, hasChangedDtsFile bool)
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

func (b *buildHost) OnBuildInfoEmit(config *tsoptions.ParsedCommandLine, buildInfo *BuildInfo, hasChangedDtsFile bool) {
	// no-op
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
