package util

import (
	"fmt"
	"os"
)

var (
	logFile *os.File
)

func openLog() {
	var err error
	logFile, err = os.Create("log.txt")
	must(err)
}

func closeLog() {
	if logFile == nil {
		return
	}
	logFile.Close()
	logFile = nil
}

/*
// TODO: should take additional format and args for optional message
func logError(err error) {
	if err != nil {
		return
	}
	logf("%s", err.Error())
}
*/

func Logf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if logFile != nil {
		fmt.Fprint(logFile, s)
	}
	fmt.Print(s)
}

// TODO: have just one
func LogVerbose(format string, args ...interface{}) {
	Verbose(format, args...)
}

func Verbose(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	if logFile != nil {
		fmt.Fprint(logFile, s)
	}
}

var (
	doTempLog = false
)

func LogTemp(format string, args ...interface{}) {
	if !doTempLog {
		return
	}
	Logf(format, args...)
}
