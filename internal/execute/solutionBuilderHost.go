package execute

import (
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type configAndTime struct {
	resolved *tsoptions.ParsedCommandLine
	time     time.Duration
}

type buildInfoAndConfig struct {
	buildInfo *incremental.BuildInfo
	config    tspath.Path
}

type solutionBuilderHost struct {
	builder             *solutionBuilder
	host                compiler.CompilerHost
	extendedConfigCache collections.SyncMap[tspath.Path, *tsoptions.ExtendedConfigCacheEntry]
	sourceFiles         collections.SyncMap[ast.SourceFileParseOptions, *ast.SourceFile]
	resolvedReferences  collections.SyncMap[tspath.Path, *configAndTime]

	buildInfos            collections.SyncMap[tspath.Path, *buildInfoAndConfig]
	mTimes                collections.SyncMap[tspath.Path, time.Time]
	latestChangedDtsFiles collections.SyncMap[tspath.Path, time.Time]
}

var (
	_ vfs.FS                      = (*solutionBuilderHost)(nil)
	_ compiler.CompilerHost       = (*solutionBuilderHost)(nil)
	_ incremental.BuildInfoReader = (*solutionBuilderHost)(nil)
)

func (h *solutionBuilderHost) FS() vfs.FS {
	return h
}

func (h *solutionBuilderHost) UseCaseSensitiveFileNames() bool {
	return h.host.FS().UseCaseSensitiveFileNames()
}

func (h *solutionBuilderHost) FileExists(path string) bool {
	return h.host.FS().FileExists(path)
}

func (h *solutionBuilderHost) ReadFile(path string) (string, bool) {
	return h.host.FS().ReadFile(path)
}

func (h *solutionBuilderHost) WriteFile(path string, data string, writeByteOrderMark bool) error {
	err := h.host.FS().WriteFile(path, data, writeByteOrderMark)
	if err == nil {
		filePath := h.builder.toPath(path)
		h.buildInfos.Delete(filePath)
		h.mTimes.Delete(filePath)
	}
	return err
}

func (h *solutionBuilderHost) Remove(path string) error {
	return h.host.FS().Remove(path)
}

func (h *solutionBuilderHost) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return h.host.FS().Chtimes(path, aTime, mTime)
}

func (h *solutionBuilderHost) DirectoryExists(path string) bool {
	return h.host.FS().DirectoryExists(path)
}

func (h *solutionBuilderHost) GetAccessibleEntries(path string) vfs.Entries {
	return h.host.FS().GetAccessibleEntries(path)
}

func (h *solutionBuilderHost) Stat(path string) vfs.FileInfo {
	return h.host.FS().Stat(path)
}

func (h *solutionBuilderHost) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return h.host.FS().WalkDir(root, walkFn)
}

func (h *solutionBuilderHost) Realpath(path string) string {
	return h.host.FS().Realpath(path)
}

func (h *solutionBuilderHost) DefaultLibraryPath() string {
	return h.host.DefaultLibraryPath()
}

func (h *solutionBuilderHost) GetCurrentDirectory() string {
	return h.host.GetCurrentDirectory()
}

func (h *solutionBuilderHost) Trace(msg string) {
	//!!! TODO: implement
}

func (h *solutionBuilderHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	if existing, loaded := h.sourceFiles.Load(opts); loaded {
		return existing
	}

	file := h.host.GetSourceFile(opts)
	file, _ = h.sourceFiles.LoadOrStore(opts, file)
	return file
}

func (h *solutionBuilderHost) GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine {
	if existing, loaded := h.resolvedReferences.Load(path); loaded {
		return existing.resolved
	}
	configStart := h.builder.opts.sys.Now()
	commandLine, _ := tsoptions.GetParsedCommandLineOfConfigFilePath(fileName, path, h.builder.opts.command.CompilerOptions, h, &h.extendedConfigCache)
	configTime := h.builder.opts.sys.Now().Sub(configStart)
	configAndTime, _ := h.resolvedReferences.LoadOrStore(path, &configAndTime{resolved: commandLine, time: configTime})
	return configAndTime.resolved
}

func (h *solutionBuilderHost) ReadBuildInfo(buildInfoFileName string) *incremental.BuildInfo {
	path := h.builder.toPath(buildInfoFileName)
	if existing, loaded := h.buildInfos.Load(path); loaded {
		return existing.buildInfo
	}
	return nil
}

func (h *solutionBuilderHost) readOrStoreBuildInfo(configPath tspath.Path, buildInfoFileName string) *incremental.BuildInfo {
	if existing, loaded := h.buildInfos.Load(h.builder.toPath(buildInfoFileName)); loaded {
		return existing.buildInfo
	}

	buildInfo := incremental.NewBuildInfoReader(h).ReadBuildInfo(buildInfoFileName)
	entry := &buildInfoAndConfig{buildInfo, configPath}
	entry, _ = h.buildInfos.LoadOrStore(h.builder.toPath(buildInfoFileName), entry)
	return entry.buildInfo
}

func (h *solutionBuilderHost) hasConflictingBuildInfo(configPath tspath.Path) bool {
	if existing, loaded := h.buildInfos.Load(configPath); loaded {
		return existing.config != configPath
	}
	return false
}

func (h *solutionBuilderHost) getMTime(file string) time.Time {
	path := h.builder.toPath(file)
	if existing, loaded := h.mTimes.Load(path); loaded {
		return existing
	}
	stat := h.host.FS().Stat(file)
	var mTime time.Time
	if stat != nil {
		mTime = stat.ModTime()
	}
	mTime, _ = h.mTimes.LoadOrStore(path, mTime)
	return mTime
}

func (h *solutionBuilderHost) setMTime(file string, mTime time.Time) error {
	path := h.builder.toPath(file)
	err := h.host.FS().Chtimes(file, time.Time{}, mTime)
	if err == nil {
		h.mTimes.Store(path, mTime)
	}
	return err
}

func (h *solutionBuilderHost) getLatestChangedDtsTime(config string) time.Time {
	path := h.builder.toPath(config)
	if existing, loaded := h.latestChangedDtsFiles.Load(path); loaded {
		return existing
	}

	var changedDtsTime time.Time
	if configAndTime, loaded := h.resolvedReferences.Load(path); loaded {
		buildInfoPath := configAndTime.resolved.GetBuildInfoFileName()
		buildInfo := h.readOrStoreBuildInfo(path, buildInfoPath)
		if buildInfo != nil && buildInfo.LatestChangedDtsFile != "" {
			changedDtsTime = h.getMTime(
				tspath.GetNormalizedAbsolutePath(
					buildInfo.LatestChangedDtsFile,
					tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoPath, h.GetCurrentDirectory())),
				),
			)
		}
	}

	changedDtsTime, _ = h.mTimes.LoadOrStore(path, changedDtsTime)
	return changedDtsTime
}
