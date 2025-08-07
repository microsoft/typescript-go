currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::

tsgo --build --help
ExitStatus:: Success
Output::
Version FakeTSVersion
tsc: The TypeScript Compiler - Version FakeTSVersion

BUILD OPTIONS

--help, -h
Print this message.

--help, -?


--watch, -w
Watch input files.

--preserveWatchOutput
Disable wiping the console in watch mode.
type: boolean
default: false

--listFiles
Print all of the files read during the compilation.
type: boolean
default: false

--explainFiles
Print files read during the compilation including why it was included.
type: boolean
default: false

--listEmittedFiles
Print the names of emitted files after a compilation.
type: boolean
default: false

--pretty
Enable color and formatting in TypeScript's output to make compiler errors easier to read.
type: boolean
default: true

--traceResolution
Log paths used during the 'moduleResolution' process.
type: boolean
default: false

--diagnostics
Output compiler performance information after building.
type: boolean
default: false

--extendedDiagnostics
Output more detailed compiler performance information after building.
type: boolean
default: false

--generateCpuProfile
Emit a v8 CPU profile of the compiler run for debugging.
type: string
default: profile.cpuprofile

--generateTrace
Generates an event trace and a list of types.

--incremental, -i
Save .tsbuildinfo files to allow for incremental compilation of projects.
type: boolean
default: `false`, unless `composite` is set

--declaration, -d
Generate .d.ts files from TypeScript and JavaScript files in your project.
type: boolean
default: `false`, unless `composite` is set

--declarationMap
Create sourcemaps for d.ts files.
type: boolean
default: false

--emitDeclarationOnly
Only output d.ts files and not JavaScript files.
type: boolean
default: false

--sourceMap
Create source map files for emitted JavaScript files.
type: boolean
default: false

--inlineSourceMap
Include sourcemap files inside the emitted JavaScript.
type: boolean
default: false

--noCheck
Disable full type checking (only critical parse and emit errors will be reported).
type: boolean
default: false

--noEmit
Disable emitting files from a compilation.
type: boolean
default: false

--assumeChangesOnlyAffectDirectDependencies
Have recompiles in projects that use 'incremental' and 'watch' mode assume that changes within a file will only affect files directly depending on it.
type: boolean
default: false

--locale
Set the language of the messaging from TypeScript. This does not affect emit.

--quiet, -q
Do not print diagnostics.

--singleThreaded
Run in single threaded mode.

--pprofDir
Generate pprof CPU/memory profiles to the given directory.

--verbose, -v
Enable verbose logging.

--dry, -d
Show what would be built (or deleted, if specified with '--clean')

--force, -f
Build all projects, including those that appear to be up to date.

--clean
Delete the outputs of all projects.

--stopBuildOnErrors
Skip building downstream projects on error in upstream project.


