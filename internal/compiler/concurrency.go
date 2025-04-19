package compiler

import (
	"runtime"
	"strconv"
	"strings"
)

// Concurrency controls the number of concurrent operations used by the Program.
// If greater than 0, it specifies the number of concurrent checkers to use.
type Concurrency int

const (
	// Use 4 checkers for checking, but maximum concurrency for all other operations.
	ConcurrencyDefault Concurrency = 0
	// Use a single thread for all operations.
	ConcurrencySingleThreaded Concurrency = -1
	// Use all available threads for all operations.
	ConcurrencyMaxProcs Concurrency = -2
	// Create a checker for each file in the program.
	ConcurrencyCheckerPerFile Concurrency = -3
)

func (c Concurrency) IsSingleThreaded() bool {
	return c == ConcurrencySingleThreaded
}

func (c Concurrency) String() string {
	switch c {
	case ConcurrencyDefault:
		return "default"
	case ConcurrencySingleThreaded:
		return "single"
	case ConcurrencyMaxProcs:
		return "max"
	case ConcurrencyCheckerPerFile:
		return "checker-per-file"
	default:
		return "concurrency-" + strconv.Itoa(int(c))
	}
}

func (c Concurrency) Checkers(numFiles int) int {
	switch c {
	case ConcurrencyDefault:
		return 4
	case ConcurrencySingleThreaded:
		return 1
	case ConcurrencyMaxProcs:
		return runtime.GOMAXPROCS(0)
	case ConcurrencyCheckerPerFile:
		return numFiles
	default:
		return int(c)
	}
}

func ParseConcurrency(s string) (Concurrency, error) {
	s = strings.ToLower(s)
	switch s {
	case "default":
		return ConcurrencyDefault, nil
	case "single":
		return ConcurrencySingleThreaded, nil
	case "max":
		return ConcurrencyMaxProcs, nil
	case "checker-per-file":
		return ConcurrencyCheckerPerFile, nil
	default:
		c, err := strconv.Atoi(s)
		return Concurrency(c), err
	}
}

func MustParseConcurrency(s string) Concurrency {
	c, err := ParseConcurrency(s)
	if err != nil {
		panic(err)
	}
	return c
}
