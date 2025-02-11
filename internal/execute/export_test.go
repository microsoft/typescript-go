package execute

import "github.com/microsoft/typescript-go/internal/tsoptions"

func CommandLineTest(sys System, cb cbType, commandLineArgs []string) (*tsoptions.ParsedCommandLine, ExitStatus) {
	parsedCommandLine := tsoptions.ParseCommandLine(commandLineArgs, sys)
	exit, _ := executeCommandLineWorker(sys, cb, parsedCommandLine)
	return parsedCommandLine, exit
}
