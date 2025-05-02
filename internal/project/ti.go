package project

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/semver"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type TypesMapFile struct {
	TypesMap  map[string]SafeListEntry `json:"typesMap"`
	SimpleMap map[string]string        `json:"simpleMap,omitzero"`
}

type SafeListEntry struct {
	Match   string   `json:"match"`
	Exclude [][]any  `json:"exclude"`
	Types   []string `json:"types"`
}

type PendingRequest struct {
	requestId              int32
	packageNames           []string
	filteredTypings        []string
	currentlyCachedTypings []string
	p                      *Project
}

type NpmInstallOperation func(string, []string) ([]byte, error)

type TypingsInstaller struct {
	TypingsLocation string
	NpmInstall      NpmInstallOperation

	// 	const typingSafeListLocation = ts.server.findArgument(ts.server.Arguments.TypingSafeListLocation);
	// const typesMapLocation = ts.server.findArgument(ts.server.Arguments.TypesMapLocation);
	// const npmLocation = ts.server.findArgument(ts.server.Arguments.NpmLocation);
	// const validateDefaultNpmLocation = ts.server.hasArgument(ts.server.Arguments.ValidateDefaultNpmLocation);
	ThrottleLimit int

	initialized   bool
	initializedMu sync.Mutex

	packageNameToTypingLocation collections.SyncMap[string, *CachedTyping]
	missingTypingsSet           collections.SyncMap[string, bool]

	typesRegistry map[string]map[string]string
	typesMap      *TypesMapFile
	safeList      map[string]string

	installRunCount      atomic.Int32
	inFlightRequestCount int
	pendingRunRequests   []*PendingRequest
	pendingRunRequestsMu sync.Mutex
}

func (ti *TypingsInstaller) IsKnownTypesPackageName(p *Project, name string) bool {
	// We want to avoid looking this up in the registry as that is expensive. So first check that it's actually an NPM package.
	validationResult, _, _ := ValidatePackageName(name)
	if validationResult != NameOk {
		return false
	}
	// Strada did this lazily - is that needed here to not waiting on and returning false on first request
	ti.init(p)
	_, ok := ti.typesRegistry[name]
	return ok
}

func (ti *TypingsInstaller) InstallPackage(p *Project, fileName string, packageName string) {
	cwd, ok := tspath.ForEachAncestorDirectory(tspath.GetDirectoryPath(fileName), func(directory string) (string, bool) {
		if p.FS().FileExists(tspath.CombinePaths(directory, "package.json")) {
			return directory, true
		}
		return "", false
	})
	if !ok {
		cwd = p.GetCurrentDirectory()
	}
	if cwd != "" {
		go ti.installWorker(p, -1, []string{packageName}, cwd, func(
			p *Project,
			requestId int32,
			packageNames []string,
			success bool,
		) {
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
		})
	} else {
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

func (ti *TypingsInstaller) EnqueueInstallTypingsRequest(p *Project, typingInfo *TypingsCacheInfo) {
	// because we arent using buffers, no need to throttle for requests here
	if p.HasLogLevel(LogLevelVerbose) {
		p.Log("TIAdapter:: Got install request for: " + p.Name())
	}
	go ti.discoverAndInstallTypings(
		p,
		typingInfo,
		p.GetFileNames( /*excludeFilesFromExternalLibraries*/ true /*excludeConfigFiles*/, true),
		p.GetCurrentDirectory(),
	) //.concat(project.getExcludedFiles())
}

func (ti *TypingsInstaller) discoverAndInstallTypings(p *Project, typingsInfo *TypingsCacheInfo, fileNames []string, projectRootPath string) {
	ti.init((p))

	ti.initializeSafeList(p)

	cachedTypingPaths, newTypingNames, filesToWatch := DiscoverTypings(
		p,
		typingsInfo,
		fileNames,
		projectRootPath,
		ti.safeList,
		&ti.packageNameToTypingLocation,
		ti.typesRegistry,
	)

	// start watching files
	p.WatchTypingLocations(filesToWatch)

	// install typings
	if len(newTypingNames) > 0 {
		ti.installTypings(p, cachedTypingPaths, newTypingNames)
	} else {
		p.Logf("No new typings were requested as a result of typings discovery")
		p.UpdateTypingFiles(cachedTypingPaths)
		// DO we really need these events
		// this.event(response, "setTypings");
	}
}

func (ti *TypingsInstaller) installTypings(
	p *Project,
	currentlyCachedTypings []string,
	typingsToInstall []string,
) {
	p.Logf("Installing typings %v", typingsToInstall)
	filteredTypings := ti.filterTypings(p, typingsToInstall)
	if len(filteredTypings) == 0 {
		p.Logf("All typings are known to be missing or invalid - no need to install more typings")
		p.UpdateTypingFiles(currentlyCachedTypings)
		// DO we really need these events
		// this.event(response, "setTypings");
		return
	}

	// ti.ensureTypingsLocationExists(p)

	requestId := ti.installRunCount.Add(1)

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
		scopedTypings[i] = fmt.Sprintf("@types/%s@ts%s", packageName, core.VersionMajorMinor)
	}

	ti.pendingRunRequestsMu.Lock()
	request := &PendingRequest{
		requestId:              requestId,
		packageNames:           scopedTypings,
		filteredTypings:        filteredTypings,
		currentlyCachedTypings: currentlyCachedTypings,
		p:                      p,
	}
	if (ti.inFlightRequestCount + 1) < ti.ThrottleLimit {
		ti.inFlightRequestCount++
		ti.pendingRunRequestsMu.Unlock()
		ti.invokeRoutineToInstallTypings(request)
	} else {
		ti.pendingRunRequests = append(ti.pendingRunRequests, request)
		ti.pendingRunRequestsMu.Unlock()
	}
}

func (ti *TypingsInstaller) invokeRoutineToInstallTypings(
	request *PendingRequest,
) {
	go ti.installWorker(
		request.p,
		request.requestId,
		request.packageNames,
		ti.TypingsLocation,
		func(
			p *Project,
			requestId int32,
			packageNames []string,
			success bool,
		) {
			ti.pendingRunRequestsMu.Lock()
			pendingRequestsCount := len(ti.pendingRunRequests)
			var nextRequest *PendingRequest
			if pendingRequestsCount == 0 {
				ti.inFlightRequestCount--
			} else {
				nextRequest = ti.pendingRunRequests[0]
				if pendingRequestsCount == 1 {
					ti.pendingRunRequests = ti.pendingRunRequests[0:0]
				} else {
					ti.pendingRunRequests = ti.pendingRunRequests[1:pendingRequestsCount]
				}
			}
			ti.pendingRunRequestsMu.Unlock()
			if nextRequest != nil {
				ti.invokeRoutineToInstallTypings(nextRequest)
			}

			if success {
				p.Logf("Installed typings %v", packageNames)
				var installedTypingFiles []string
				resolver := module.NewResolver(p, &core.CompilerOptions{ModuleResolution: core.ModuleResolutionKindNodeNext})
				for _, packageName := range request.filteredTypings {
					typingFile := ti.typingToFileName(resolver, packageName)
					if typingFile == "" {
						ti.missingTypingsSet.Store(packageName, true)
						continue
					}

					// packageName is guaranteed to exist in typesRegistry by filterTypings
					distTags := ti.typesRegistry[packageName]
					useVersion, ok := distTags["ts"+core.VersionMajorMinor]
					if !ok {
						useVersion = distTags["latest"]
					}
					newVersion := semver.MustParse(useVersion)
					newTyping := &CachedTyping{Location: typingFile, Version: newVersion}
					ti.packageNameToTypingLocation.Store(packageName, newTyping)
					installedTypingFiles = append(installedTypingFiles, typingFile)
				}
				p.Logf("Installed typing files %v", installedTypingFiles)
				p.UpdateTypingFiles(append(request.currentlyCachedTypings, installedTypingFiles...))
				// DO we really need these events
				// this.event(response, "setTypings");

			} else {
				p.Logf("install request failed, marking packages as missing to prevent repeated requests: %v", request.filteredTypings)
				for _, typing := range request.filteredTypings {
					ti.missingTypingsSet.Store(typing, true)
				}
			}

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
		},
	)
}

func (ti *TypingsInstaller) installWorker(
	p *Project,
	requestId int32,
	packageNames []string,
	cwd string,
	onRequestComplete func(
		p *Project,
		requestId int32,
		packageNames []string,
		success bool,
	),
) {
	p.Logf("#%d with cwd: %s arguments: %v", requestId, cwd, packageNames)
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
			ti.queueTypingsInstall(wg, p, cwd, packageNames[currentCommandStart:currentCommandEnd], &hasError)
			currentCommandStart = currentCommandEnd
			currentCommandSize = 100 + len(packageName) + 1
			currentCommandEnd++
		}
	}
	ti.queueTypingsInstall(wg, p, cwd, packageNames[currentCommandStart:currentCommandEnd], &hasError)
	wg.RunAndWait()

	p.Logf("npm install #%d completed", requestId)
	onRequestComplete(p, requestId, packageNames, !hasError.Load())
}

func (ti *TypingsInstaller) queueTypingsInstall(
	wg core.WorkGroup,
	p *Project,
	cwd string,
	packages []string,
	hasError *atomic.Bool,
) {
	ti.ensureNpmInstall()
	wg.Queue(func() {
		var npmArgs []string
		npmArgs = append(npmArgs, "install", "--ignore-scripts")
		npmArgs = append(npmArgs, packages...)
		npmArgs = append(npmArgs, "--save-dev", "--user-agent=\"typesInstaller/"+core.Version+"\"")
		output, err := ti.NpmInstall(cwd, npmArgs)
		if err != nil {
			p.Logf("Output is: %s", output)
			hasError.Store(true)
		}
	})
}

func (ti *TypingsInstaller) filterTypings(
	p *Project,
	typingsToInstall []string,
) []string {
	var result []string
	for _, typing := range typingsToInstall {
		typingKey := module.MangleScopedPackageName(typing)
		if _, ok := ti.missingTypingsSet.Load(typingKey); ok {
			p.Logf("'%s':: '%s' is in missingTypingsSet - skipping...", typing, typingKey)
			continue
		}
		validationResult, name, isScopeName := ValidatePackageName(typing)
		if validationResult != NameOk {
			// add typing name to missing set so we won't process it again
			ti.missingTypingsSet.Store(typingKey, true)
			p.Log(RenderPackageNameValidationFailure(typing, validationResult, name, isScopeName))
			continue
		}
		typesRegistryEntry, ok := ti.typesRegistry[typingKey]
		if !ok {
			p.Logf(`'%s':: Entry for package '%s' does not exist in local types registry - skipping...`, typing, typingKey)
			continue
		}
		if typingLocation, ok := ti.packageNameToTypingLocation.Load(typingKey); ok && IsTypingUpToDate(typingLocation, typesRegistryEntry) {
			p.Logf("'%s':: '%s' already has an up-to-date typing - skipping...", typing, typingKey)
			continue
		}
		result = append(result, typingKey)
	}
	return result
}

func (ti *TypingsInstaller) init(p *Project) {
	ti.initializedMu.Lock()
	if ti.initialized {
		ti.initializedMu.Unlock()
		return
	}
	p.Log("Global cache location '" + ti.TypingsLocation + "'") //, safe file path '" + safeListPath + "', types map path '" + typesMapLocation + "`")
	ti.processCacheLocation(p)

	//     // If the NPM path contains spaces and isn't wrapped in quotes, do so.
	//     if (this.npmPath.includes(" ") && this.npmPath[0] !== `"`) {
	//         this.npmPath = `"${this.npmPath}"`;
	//     }
	//     if (this.log.isEnabled()) {
	//         this.log.writeLine(`Process id: ${process.pid}`);
	//         this.log.writeLine(`NPM location: ${this.npmPath} (explicit '${ts.server.Arguments.NpmLocation}' ${npmLocation === undefined ? "not " : ""} provided)`);
	//         this.log.writeLine(`validateDefaultNpmLocation: ${validateDefaultNpmLocation}`);
	//     }

	ti.ensureTypingsLocationExists(p)
	p.Log("Updating types-registry@latest npm package...")
	ti.ensureNpmInstall()
	if _, err := ti.NpmInstall(ti.TypingsLocation, []string{"install", "--ignore-scripts", "types-registry@latest"}); err == nil {
		p.Log("Updated types-registry npm package")
	} else {
		p.Logf("Error updating types-registry package: %v", err)
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

	ti.typesRegistry = ti.loadTypesRegistryFile(p)
	ti.initialized = true
	ti.initializedMu.Unlock()
}

func (ti *TypingsInstaller) ensureNpmInstall() {
	if ti.NpmInstall == nil {
		ti.NpmInstall = npmInstall
	}
}

type NpmConfig struct {
	DevDependencies map[string]any `json:"devDependencies"`
}

type NpmDependecyEntry struct {
	Version string `json:"version"`
}
type NpmLock struct {
	Dependencies map[string]NpmDependecyEntry `json:"dependencies"`
}

func (ti *TypingsInstaller) processCacheLocation(p *Project) {
	p.Log("Processing cache location " + ti.TypingsLocation)
	packageJson := tspath.CombinePaths(ti.TypingsLocation, "package.json")
	packageLockJson := tspath.CombinePaths(ti.TypingsLocation, "package-lock.json")
	p.Log("Trying to find '" + packageJson + "'...")
	if p.FS().FileExists(packageJson) && p.FS().FileExists((packageLockJson)) {
		var npmConfig NpmConfig
		npmConfigContents := parseNpmConfigOrLock(p, packageJson, &npmConfig)
		var npmLock NpmLock
		npmLockContents := parseNpmConfigOrLock(p, packageLockJson, &npmLock)

		p.Log("Loaded content of " + packageJson + ": " + npmConfigContents)
		p.Log("Loaded content of " + packageLockJson + ": " + npmLockContents)

		// TODO:: Not next but Node10 in strada
		resolver := module.NewResolver(p, &core.CompilerOptions{ModuleResolution: core.ModuleResolutionKindNodeNext})
		if npmConfig.DevDependencies != nil && npmLock.Dependencies != nil {
			for key := range npmConfig.DevDependencies {
				npmLockValue, npmLockValueExists := npmLock.Dependencies[key]
				if !npmLockValueExists {
					// if package in package.json but not package-lock.json, skip adding to cache so it is reinstalled on next use
					continue
				}
				// key is @types/<package name>
				packageName := tspath.GetBaseFileName(key)
				if packageName == "" {
					continue
				}
				typingFile := ti.typingToFileName(resolver, packageName)
				if typingFile == "" {
					ti.missingTypingsSet.Store(packageName, true)
					continue
				}
				if existingTypingFile, existingTypingsFilePresent := ti.packageNameToTypingLocation.Load(packageName); existingTypingsFilePresent {
					if existingTypingFile.Location == typingFile {
						continue
					}
					p.Log("New typing for package " + packageName + " from " + typingFile + " conflicts with existing typing file " + existingTypingFile.Location)
				}
				p.Log("Adding entry into typings cache: " + packageName + " => " + typingFile)
				version := npmLockValue.Version
				if version == "" {
					continue
				}

				newTyping := &CachedTyping{
					Location: typingFile,
					Version:  semver.MustParse(version),
				}
				ti.packageNameToTypingLocation.Store(packageName, newTyping)
			}
		}
	}
	p.Log("Finished processing cache location " + ti.TypingsLocation)
}

func parseNpmConfigOrLock[T NpmConfig | NpmLock](p *Project, location string, config *T) string {
	contents, _ := p.FS().ReadFile(location)
	_ = json.Unmarshal([]byte(contents), config)
	return contents
}

func (ti *TypingsInstaller) ensureTypingsLocationExists(p *Project) {
	npmConfigPath := tspath.CombinePaths(ti.TypingsLocation, "package.json")
	p.Log("Npm config file: " + npmConfigPath)

	if !p.FS().FileExists(npmConfigPath) {
		p.Logf("Npm config file: '%s' is missing, creating new one...", npmConfigPath)
		err := p.FS().WriteFile(npmConfigPath, "{ \"private\": true }", false)
		if err != nil {
			p.Logf("Npm config file write failed: %v", err)
		}
	}
}

func (ti *TypingsInstaller) typingToFileName(resolver *module.Resolver, packageName string) string {
	result := resolver.ResolveModuleName(packageName, tspath.CombinePaths(ti.TypingsLocation, "index.d.ts"), core.ModuleKindNone, nil)
	return result.ResolvedFileName
}

func (ti *TypingsInstaller) loadTypesRegistryFile(p *Project) map[string]map[string]string {
	typesRegistryFile := tspath.CombinePaths(ti.TypingsLocation, "node_modules/types-registry/index.json")
	typesRegistryFileContents, ok := p.FS().ReadFile(typesRegistryFile)
	if ok {
		var entries map[string]map[string]map[string]string
		err := json.Unmarshal([]byte(typesRegistryFileContents), &entries)
		if err == nil {
			if typesRegistry, ok := entries["entries"]; ok {
				return typesRegistry
			}
		}
		p.Logf("Error when loading types registry file '%s': %v", typesRegistryFile, err)
	} else {
		p.Logf("Error reading types registry file '%s'", typesRegistryFile)
	}
	return map[string]map[string]string{}
}

func (ti *TypingsInstaller) initializeSafeList(p *Project) {
	if ti.safeList != nil {
		return
	}
	ti.loadTypesMap(p)
	if ti.typesMap.SimpleMap != nil {
		p.Logf("Loaded safelist from types map file '%s'", tspath.CombinePaths(p.DefaultLibraryPath(), "typesMap.json"))
		ti.safeList = ti.typesMap.SimpleMap
		return
	}

	p.Logf("Failed to load safelist from types map file '$%s'", tspath.CombinePaths(p.DefaultLibraryPath(), "typesMap.json"))
	ti.loadSafeList(p)
}

func (ti *TypingsInstaller) loadTypesMap(p *Project) {
	if ti.typesMap != nil {
		return
	}
	typesMapLocation := tspath.CombinePaths(p.DefaultLibraryPath(), "typesMap.json")
	typesMapContents, ok := p.FS().ReadFile(typesMapLocation)
	if ok {
		err := json.Unmarshal([]byte(typesMapContents), &ti.typesMap)
		if err != nil {
			return
		}
		p.Logf("Error when parsing typesMapLocation '%s': %v", typesMapLocation, err)
	} else {
		p.Logf("Error reading typesMapLocation '%s'", typesMapLocation)
	}
	ti.typesMap = &TypesMapFile{}
}

func (ti *TypingsInstaller) loadSafeList(p *Project) {
	safeListLocation := tspath.CombinePaths(p.DefaultLibraryPath(), "typingSafeList.json")
	safeListContents, ok := p.FS().ReadFile(safeListLocation)
	if ok {
		err := json.Unmarshal([]byte(safeListContents), &ti.safeList)
		if err != nil {
			return
		}
		p.Logf("Error when parsing safeListLocation '%s': %v", safeListLocation, err)
	} else {
		p.Logf("Error reading safeListLocation '%s'", safeListLocation)
	}
	ti.safeList = map[string]string{}
}

func npmInstall(cwd string, npmInstallArgs []string) ([]byte, error) {
	cmd := exec.Command("npm", npmInstallArgs...)
	cmd.Dir = cwd
	return cmd.Output()
}

func GetGlobalTypingsCacheLocation() string {
	switch runtime.GOOS {
	case "windows":
		{
			basePath, err := os.UserCacheDir()
			if err != nil {
				if basePath, err = os.UserConfigDir(); err != nil {
					if basePath, err = os.UserHomeDir(); err != nil {
						if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
							basePath = userProfile
						} else if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); homeDrive != "" && homePath != "" {
							basePath = homeDrive + homePath
						} else {
							basePath = os.TempDir()
						}
					}
				}
			}
			return tspath.CombinePaths(tspath.CombinePaths(basePath, "Microsoft/TypeScript"), core.VersionMajorMinor)
		}
	case "openbsd", "freebsd", "netbsd", "darwin", "linux", "android":
		{
			cacheLocation := getNonWindowsCacheLocation()
			return tspath.CombinePaths(tspath.CombinePaths(cacheLocation, "typescript"), core.VersionMajorMinor)
		}
	default:
		panic("unsupported platform: " + runtime.GOOS)
	}
}

func getNonWindowsCacheLocation() string {
	if xdgCacheHome := os.Getenv("XDG_CACHE_HOME"); xdgCacheHome != "" {
		return xdgCacheHome
	}
	const platformIsDarwin = runtime.GOOS == "darwin"
	var usersDir string
	if platformIsDarwin {
		usersDir = "Users"
	} else {
		usersDir = "home"
	}
	homePath, err := os.UserHomeDir()
	if err != nil {
		if home := os.Getenv("HOME"); home != "" {
			homePath = home
		} else {
			var userName string
			if logName := os.Getenv("LOGNAME"); logName != "" {
				userName = logName
			} else if user := os.Getenv("USER"); user != "" {
				userName = user
			}
			if userName != "" {
				homePath = "/" + usersDir + "/" + userName
			} else {
				homePath = os.TempDir()
			}
		}
	}
	var cacheFolder string
	if platformIsDarwin {
		cacheFolder = "Library/Caches"
	} else {
		cacheFolder = ".cache"
	}
	return tspath.CombinePaths(homePath, cacheFolder)
}
