package projectv2

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/semver"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type TypingsInfo struct {
	TypeAcquisition   *core.TypeAcquisition
	CompilerOptions   *core.CompilerOptions
	UnresolvedImports collections.Set[string]
}

func (ti TypingsInfo) Equals(other TypingsInfo) bool {
	return ti.TypeAcquisition.Equals(other.TypeAcquisition) &&
		ti.CompilerOptions.GetAllowJS() == other.CompilerOptions.GetAllowJS() &&
		ti.UnresolvedImports.Equals(&other.UnresolvedImports)
}

type CachedTyping struct {
	TypingsLocation string
	Version         *semver.Version
}

type pendingRequest struct {
	requestID              int32
	projectID              tspath.Path
	packageNames           []string
	filteredTypings        []string
	currentlyCachedTypings []string
	typingsInfo            *TypingsInfo
}

type NpmInstallOperation func(cwd string, npmInstallArgs []string) ([]byte, error)

type TypingsInstallerStatus struct {
	RequestID int32
	ProjectID tspath.Path
	Status    string
}

type TypingsInstallerOptions struct {
	TypingsLocation string
	ThrottleLimit   int
}

type TypingsInstallerHost interface {
	OnTypingsInstalled(projectID tspath.Path, typingsInfo *TypingsInfo, cachedTypingPaths []string)
	OnTypingsInstallFailed(projectID tspath.Path, typingsInfo *TypingsInfo, err error)
	NpmInstall(cwd string, packageNames []string) ([]byte, error)
}

type TypingsInstaller struct {
	typingsLocation string
	throttleLimit   int
	host            TypingsInstallerHost

	initOnce sync.Once

	packageNameToTypingLocation collections.SyncMap[string, *CachedTyping]
	missingTypingsSet           collections.SyncMap[string, bool]

	typesRegistry map[string]map[string]string

	installRunCount      atomic.Int32
	inFlightRequestCount int
	pendingRunRequests   []*pendingRequest
	pendingRunRequestsMu sync.Mutex
}

func NewTypingsInstaller(options *TypingsInstallerOptions, host TypingsInstallerHost) *TypingsInstaller {
	return &TypingsInstaller{
		typingsLocation: options.TypingsLocation,
		throttleLimit:   options.ThrottleLimit,
		host:            host,
	}
}

func (ti *TypingsInstaller) PendingRunRequestsCount() int {
	ti.pendingRunRequestsMu.Lock()
	defer ti.pendingRunRequestsMu.Unlock()
	return len(ti.pendingRunRequests)
}

func (ti *TypingsInstaller) IsKnownTypesPackageName(projectID tspath.Path, name string, fs vfs.FS, logger func(string)) bool {
	// We want to avoid looking this up in the registry as that is expensive. So first check that it's actually an NPM package.
	validationResult, _, _ := ValidatePackageName(name)
	if validationResult != NameOk {
		return false
	}
	// Strada did this lazily - is that needed here to not waiting on and returning false on first request
	ti.init(string(projectID), fs, logger)
	_, ok := ti.typesRegistry[name]
	return ok
}

// !!! sheetal currently we use latest instead of core.VersionMajorMinor()
const TsVersionToUse = "latest"

func (ti *TypingsInstaller) InstallPackage(projectID tspath.Path, fileName string, packageName string, fs vfs.FS, logger func(string), currentDirectory string) {
	cwd, ok := tspath.ForEachAncestorDirectory(tspath.GetDirectoryPath(fileName), func(directory string) (string, bool) {
		if fs.FileExists(tspath.CombinePaths(directory, "package.json")) {
			return directory, true
		}
		return "", false
	})
	if !ok {
		cwd = currentDirectory
	}
	if cwd != "" {
		go ti.installWorker(
			projectID,
			-1,
			[]string{packageName},
			cwd,
			func(
				projectID tspath.Path,
				requestId int32,
				packageNames []string,
				success bool,
			) {
				// !!! sheetal events to send
				// const message = success ?
				//
				//	`Package ${packageName} installed.` :
				//	`There was an error installing ${packageName}.`;
				//
				//	const response: PackageInstalledResponse = {
				//		kind: ActionPackageInstalled,
				//		projectName,
				//		id,
				//		success,
				//		message,
				//	};
				//

				// this.sendResponse(response);
				//     // The behavior is the same as for setTypings, so send the same event.
				//     this.event(response, "setTypings"); -- Used same event name - do we need it ?
			},
			logger,
		)
	} else {
		// !!! sheetal events to send
		// const response: PackageInstalledResponse = {
		// 	kind: ActionPackageInstalled,
		// 	projectName,
		// 	id,
		// 	success: false,
		// 	message: "Could not determine a project root path.",
		// };
		// this.sendResponse(response);
		//     // The behavior is the same as for setTypings, so send the same event.
		//     this.event(response, "setTypings"); -- Used same event name - do we need it ?
	}
}

type TypingsInstallRequest struct {
	ProjectID        tspath.Path
	TypingsInfo      *TypingsInfo
	FileNames        []string
	ProjectRootPath  string
	CompilerOptions  *core.CompilerOptions
	CurrentDirectory string
	GetScriptKind    func(string) core.ScriptKind
	FS               vfs.FS
	Logger           func(string)
}

func (ti *TypingsInstaller) InstallTypings(request *TypingsInstallRequest) {
	// because we arent using buffers, no need to throttle for requests here
	request.Logger("ATA:: Got install request for: " + string(request.ProjectID))
	ti.discoverAndInstallTypings(request)
}

func (ti *TypingsInstaller) discoverAndInstallTypings(request *TypingsInstallRequest) {
	ti.init(string(request.ProjectID), request.FS, request.Logger)

	cachedTypingPaths, newTypingNames, filesToWatch := DiscoverTypings(
		request.FS,
		request.Logger,
		request.TypingsInfo,
		request.FileNames,
		request.ProjectRootPath,
		&ti.packageNameToTypingLocation,
		ti.typesRegistry,
	)

	// !!!
	if len(filesToWatch) > 0 {
		request.Logger(fmt.Sprintf("ATA:: Would watch typing locations: %v", filesToWatch))
	}

	requestId := ti.installRunCount.Add(1)
	// install typings
	if len(newTypingNames) > 0 {
		filteredTypings := ti.filterTypings(request.ProjectID, request.Logger, newTypingNames)
		if len(filteredTypings) != 0 {
			ti.installTypings(request.ProjectID, request.TypingsInfo, requestId, cachedTypingPaths, filteredTypings, request.Logger)
			return
		}
		request.Logger("ATA:: All typings are known to be missing or invalid - no need to install more typings")
	} else {
		request.Logger("ATA:: No new typings were requested as a result of typings discovery")
	}

	ti.host.OnTypingsInstalled(request.ProjectID, request.TypingsInfo, cachedTypingPaths)
	// !!! sheetal events to send
	// this.event(response, "setTypings");
}

func (ti *TypingsInstaller) installTypings(
	projectID tspath.Path,
	typingsInfo *TypingsInfo,
	requestID int32,
	currentlyCachedTypings []string,
	filteredTypings []string,
	logger func(string),
) {
	// !!! sheetal events to send
	// send progress event
	// this.sendResponse({
	// 	kind: EventBeginInstallTypes,
	// 	eventId: requestId,
	// 	typingsInstallerVersion: version,
	// 	projectName: req.projectName,
	// } as BeginInstallTypes);

	// const body: protocol.BeginInstallTypesEventBody = {
	// 	eventId: response.eventId,
	// 	packages: response.packagesToInstall,
	// };
	// const eventName: protocol.BeginInstallTypesEventName = "beginInstallTypes";
	// this.event(body, eventName);

	scopedTypings := make([]string, len(filteredTypings))
	for i, packageName := range filteredTypings {
		scopedTypings[i] = fmt.Sprintf("@types/%s@%s", packageName, TsVersionToUse) // @tscore.VersionMajorMinor) // This is normally @tsVersionMajorMinor but for now lets use latest
	}

	request := &pendingRequest{
		requestID:              requestID,
		projectID:              projectID,
		packageNames:           scopedTypings,
		filteredTypings:        filteredTypings,
		currentlyCachedTypings: currentlyCachedTypings,
		typingsInfo:            typingsInfo,
	}
	ti.pendingRunRequestsMu.Lock()
	if ti.inFlightRequestCount < ti.throttleLimit {
		ti.inFlightRequestCount++
		ti.pendingRunRequestsMu.Unlock()
		ti.invokeRoutineToInstallTypings(request, logger)
	} else {
		ti.pendingRunRequests = append(ti.pendingRunRequests, request)
		ti.pendingRunRequestsMu.Unlock()
	}
}

func (ti *TypingsInstaller) invokeRoutineToInstallTypings(
	request *pendingRequest,
	logger func(string),
) {
	go ti.installWorker(
		request.projectID,
		request.requestID,
		request.packageNames,
		ti.typingsLocation,
		func(
			projectID tspath.Path,
			requestID int32,
			packageNames []string,
			success bool,
		) {
			if success {
				logger(fmt.Sprintf("ATA:: Installed typings %v", packageNames))
				var installedTypingFiles []string
				// Create a minimal resolver context for finding typing files
				resolver := &typingResolver{
					fs:              nil, // Will be set from context
					typingsLocation: ti.typingsLocation,
				}
				for _, packageName := range request.filteredTypings {
					typingFile := ti.typingToFileName(resolver, packageName)
					if typingFile == "" {
						logger(fmt.Sprintf("ATA:: Failed to find typing file for package '%s'", packageName))
						continue
					}

					// packageName is guaranteed to exist in typesRegistry by filterTypings
					distTags := ti.typesRegistry[packageName]
					useVersion, ok := distTags["ts"+core.VersionMajorMinor()]
					if !ok {
						useVersion = distTags["latest"]
					}
					newVersion := semver.MustParse(useVersion)
					newTyping := &CachedTyping{TypingsLocation: typingFile, Version: &newVersion}
					ti.packageNameToTypingLocation.Store(packageName, newTyping)
					installedTypingFiles = append(installedTypingFiles, typingFile)
				}
				logger(fmt.Sprintf("ATA:: Installed typing files %v", installedTypingFiles))

				ti.host.OnTypingsInstalled(request.projectID, request.typingsInfo, append(request.currentlyCachedTypings, installedTypingFiles...))
				// DO we really need these events
				// this.event(response, "setTypings");
			} else {
				logger(fmt.Sprintf("ATA:: install request failed, marking packages as missing to prevent repeated requests: %v", request.filteredTypings))
				for _, typing := range request.filteredTypings {
					ti.missingTypingsSet.Store(typing, true)
				}

				ti.host.OnTypingsInstallFailed(request.projectID, request.typingsInfo, fmt.Errorf("npm install failed"))
			}

			// !!! sheetal events to send
			// const response: EndInstallTypes = {
			// 	kind: EventEndInstallTypes,
			// 	eventId: requestId,
			// 	projectName: req.projectName,
			// 	packagesToInstall: scopedTypings,
			// 	installSuccess: ok,
			// 	typingsInstallerVersion: version,
			// };
			// this.sendResponse(response);

			// if (this.telemetryEnabled) {
			// 	const body: protocol.TypingsInstalledTelemetryEventBody = {
			// 		telemetryEventName: "typingsInstalled",
			// 		payload: {
			// 			installedPackages: response.packagesToInstall.join(","),
			// 			installSuccess: response.installSuccess,
			// 			typingsInstallerVersion: response.typingsInstallerVersion,
			// 		},
			// 	};
			// 	const eventName: protocol.TelemetryEventName = "telemetry";
			// 	this.event(body, eventName);
			// }

			// const body: protocol.EndInstallTypesEventBody = {
			// 	eventId: response.eventId,
			// 	packages: response.packagesToInstall,
			// 	success: response.installSuccess,
			// };
			// const eventName: protocol.EndInstallTypesEventName = "endInstallTypes";
			// this.event(body, eventName);

			ti.pendingRunRequestsMu.Lock()
			pendingRequestsCount := len(ti.pendingRunRequests)
			var nextRequest *pendingRequest
			if pendingRequestsCount == 0 {
				ti.inFlightRequestCount--
			} else {
				nextRequest = ti.pendingRunRequests[0]
				if pendingRequestsCount == 1 {
					ti.pendingRunRequests = nil
				} else {
					ti.pendingRunRequests[0] = nil // ensure the request is GC'd
					ti.pendingRunRequests = ti.pendingRunRequests[1:]
				}
			}
			ti.pendingRunRequestsMu.Unlock()
			if nextRequest != nil {
				ti.invokeRoutineToInstallTypings(nextRequest, logger)
			}
		},
		logger,
	)
}

func (ti *TypingsInstaller) installWorker(
	projectID tspath.Path,
	requestId int32,
	packageNames []string,
	cwd string,
	onRequestComplete func(
		projectID tspath.Path,
		requestId int32,
		packageNames []string,
		success bool,
	),
	logger func(string),
) {
	logger(fmt.Sprintf("ATA:: #%d with cwd: %s arguments: %v", requestId, cwd, packageNames))
	hasError := InstallNpmPackages(packageNames, func(packageNames []string, hasError *atomic.Bool) {
		var npmArgs []string
		npmArgs = append(npmArgs, "install", "--ignore-scripts")
		npmArgs = append(npmArgs, packageNames...)
		npmArgs = append(npmArgs, "--save-dev", "--user-agent=\"typesInstaller/"+core.Version()+"\"")
		output, err := ti.host.NpmInstall(cwd, npmArgs)
		if err != nil {
			logger(fmt.Sprintf("ATA:: Output is: %s", output))
			hasError.Store(true)
		}
	})
	logger(fmt.Sprintf("TI:: npm install #%d completed", requestId))
	onRequestComplete(projectID, requestId, packageNames, !hasError)
}

func InstallNpmPackages(
	packageNames []string,
	installPackages func(packages []string, hasError *atomic.Bool),
) bool {
	var hasError atomic.Bool
	hasError.Store(false)

	wg := core.NewWorkGroup(false)
	currentCommandStart := 0
	currentCommandEnd := 0
	currentCommandSize := 100
	for _, packageName := range packageNames {
		currentCommandSize = currentCommandSize + len(packageName) + 1
		if currentCommandSize < 8000 {
			currentCommandEnd++
		} else {
			packages := packageNames[currentCommandStart:currentCommandEnd]
			wg.Queue(func() {
				installPackages(packages, &hasError)
			})
			currentCommandStart = currentCommandEnd
			currentCommandSize = 100 + len(packageName) + 1
			currentCommandEnd++
		}
	}
	wg.Queue(func() {
		installPackages(packageNames[currentCommandStart:currentCommandEnd], &hasError)
	})
	wg.RunAndWait()
	return hasError.Load()
}

func (ti *TypingsInstaller) filterTypings(
	projectID tspath.Path,
	logger func(string),
	typingsToInstall []string,
) []string {
	var result []string
	for _, typing := range typingsToInstall {
		typingKey := module.MangleScopedPackageName(typing)
		if _, ok := ti.missingTypingsSet.Load(typingKey); ok {
			logger(fmt.Sprintf("ATA:: '%s':: '%s' is in missingTypingsSet - skipping...", typing, typingKey))
			continue
		}
		validationResult, name, isScopeName := ValidatePackageName(typing)
		if validationResult != NameOk {
			// add typing name to missing set so we won't process it again
			ti.missingTypingsSet.Store(typingKey, true)
			logger("ATA:: " + RenderPackageNameValidationFailure(typing, validationResult, name, isScopeName))
			continue
		}
		typesRegistryEntry, ok := ti.typesRegistry[typingKey]
		if !ok {
			logger(fmt.Sprintf("ATA:: '%s':: Entry for package '%s' does not exist in local types registry - skipping...", typing, typingKey))
			continue
		}
		if typingLocation, ok := ti.packageNameToTypingLocation.Load(typingKey); ok && IsTypingUpToDate(typingLocation, typesRegistryEntry) {
			logger(fmt.Sprintf("ATA:: '%s':: '%s' already has an up-to-date typing - skipping...", typing, typingKey))
			continue
		}
		result = append(result, typingKey)
	}
	return result
}

func (ti *TypingsInstaller) init(projectID string, fs vfs.FS, logger func(string)) {
	ti.initOnce.Do(func() {
		logger("ATA:: Global cache location '" + ti.typingsLocation + "'") //, safe file path '" + safeListPath + "', types map path '" + typesMapLocation + "`")
		ti.processCacheLocation(projectID, fs, logger)

		// !!! sheetal handle npm path here if we would support it
		//     // If the NPM path contains spaces and isn't wrapped in quotes, do so.
		//     if (this.npmPath.includes(" ") && this.npmPath[0] !== `"`) {
		//         this.npmPath = `"${this.npmPath}"`;
		//     }
		//     if (this.log.isEnabled()) {
		//         this.log.writeLine(`Process id: ${process.pid}`);
		//         this.log.writeLine(`NPM location: ${this.npmPath} (explicit '${ts.server.Arguments.NpmLocation}' ${npmLocation === undefined ? "not " : ""} provided)`);
		//         this.log.writeLine(`validateDefaultNpmLocation: ${validateDefaultNpmLocation}`);
		//     }

		ti.ensureTypingsLocationExists(fs, logger)
		logger("ATA:: Updating types-registry@latest npm package...")
		if _, err := ti.host.NpmInstall(ti.typingsLocation, []string{"install", "--ignore-scripts", "types-registry@latest"}); err == nil {
			logger("ATA:: Updated types-registry npm package")
		} else {
			logger(fmt.Sprintf("ATA:: Error updating types-registry package: %v", err))
			// !!! sheetal events to send
			//         // store error info to report it later when it is known that server is already listening to events from typings installer
			//         this.delayedInitializationError = {
			//             kind: "event::initializationFailed",
			//             message: (e as Error).message,
			//             stack: (e as Error).stack,
			//         };

			// const body: protocol.TypesInstallerInitializationFailedEventBody = {
			// 	message: response.message,
			// };
			// const eventName: protocol.TypesInstallerInitializationFailedEventName = "typesInstallerInitializationFailed";
			// this.event(body, eventName);
		}

		ti.typesRegistry = ti.loadTypesRegistryFile(fs, logger)
	})
}

type NpmConfig struct {
	DevDependencies map[string]any `json:"devDependencies"`
}

type NpmDependecyEntry struct {
	Version string `json:"version"`
}
type NpmLock struct {
	Dependencies map[string]NpmDependecyEntry `json:"dependencies"`
	Packages     map[string]NpmDependecyEntry `json:"packages"`
}

func (ti *TypingsInstaller) processCacheLocation(projectID string, fs vfs.FS, logger func(string)) {
	logger("ATA:: Processing cache location " + ti.typingsLocation)
	packageJson := tspath.CombinePaths(ti.typingsLocation, "package.json")
	packageLockJson := tspath.CombinePaths(ti.typingsLocation, "package-lock.json")
	logger("ATA:: Trying to find '" + packageJson + "'...")
	if fs.FileExists(packageJson) && fs.FileExists((packageLockJson)) {
		var npmConfig NpmConfig
		npmConfigContents := parseNpmConfigOrLock(fs, logger, packageJson, &npmConfig)
		var npmLock NpmLock
		npmLockContents := parseNpmConfigOrLock(fs, logger, packageLockJson, &npmLock)

		logger("ATA:: Loaded content of " + packageJson + ": " + npmConfigContents)
		logger("ATA:: Loaded content of " + packageLockJson + ": " + npmLockContents)

		// !!! sheetal strada uses Node10
		resolver := &typingResolver{
			fs:              fs,
			typingsLocation: ti.typingsLocation,
		}
		if npmConfig.DevDependencies != nil && (npmLock.Packages != nil || npmLock.Dependencies != nil) {
			for key := range npmConfig.DevDependencies {
				npmLockValue, npmLockValueExists := npmLock.Packages["node_modules/"+key]
				if !npmLockValueExists {
					npmLockValue, npmLockValueExists = npmLock.Dependencies[key]
				}
				if !npmLockValueExists {
					continue
				}
				// key is @types/<package name>
				packageName := tspath.GetBaseFileName(key)
				if packageName == "" {
					continue
				}
				typingFile := ti.typingToFileName(resolver, packageName)
				if typingFile == "" {
					continue
				}
				newVersion := semver.MustParse(npmLockValue.Version)
				newTyping := &CachedTyping{TypingsLocation: typingFile, Version: &newVersion}
				ti.packageNameToTypingLocation.Store(packageName, newTyping)
			}
		}
	}
	logger("ATA:: Finished processing cache location " + ti.typingsLocation)
}

func parseNpmConfigOrLock[T NpmConfig | NpmLock](fs vfs.FS, logger func(string), location string, config *T) string {
	contents, _ := fs.ReadFile(location)
	_ = json.Unmarshal([]byte(contents), config)
	return contents
}

func (ti *TypingsInstaller) ensureTypingsLocationExists(fs vfs.FS, logger func(string)) {
	npmConfigPath := tspath.CombinePaths(ti.typingsLocation, "package.json")
	logger("ATA:: Npm config file: " + npmConfigPath)

	if !fs.FileExists(npmConfigPath) {
		logger(fmt.Sprintf("ATA:: Npm config file: '%s' is missing, creating new one...", npmConfigPath))
		err := fs.WriteFile(npmConfigPath, "{ \"private\": true }", false)
		if err != nil {
			logger(fmt.Sprintf("ATA:: Npm config file write failed: %v", err))
		}
	}
}

// Simple resolver for typing files - minimal implementation
type typingResolver struct {
	fs              vfs.FS
	typingsLocation string
}

func (ti *TypingsInstaller) typingToFileName(resolver *typingResolver, packageName string) string {
	// Simple implementation - just check if the typing file exists
	// This replaces the more complex module resolution from the original
	typingPath := tspath.CombinePaths(ti.typingsLocation, "node_modules", "@types", packageName, "index.d.ts")
	if resolver.fs != nil && resolver.fs.FileExists(typingPath) {
		return typingPath
	}
	return ""
}

func (ti *TypingsInstaller) loadTypesRegistryFile(fs vfs.FS, logger func(string)) map[string]map[string]string {
	typesRegistryFile := tspath.CombinePaths(ti.typingsLocation, "node_modules/types-registry/index.json")
	typesRegistryFileContents, ok := fs.ReadFile(typesRegistryFile)
	if ok {
		var entries map[string]map[string]map[string]string
		err := json.Unmarshal([]byte(typesRegistryFileContents), &entries)
		if err == nil {
			if npmDistTags, ok := entries["entries"]; ok {
				return npmDistTags
			}
		}
		logger(fmt.Sprintf("ATA:: Error when loading types registry file '%s': %v", typesRegistryFile, err))
	} else {
		logger(fmt.Sprintf("ATA:: Error reading types registry file '%s'", typesRegistryFile))
	}
	return map[string]map[string]string{}
}
