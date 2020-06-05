package log

import (
	"io"
	oLog "log"
)

var isSetOutput = false
var customLogOutput io.Writer = nil

// Println is doing custom output validation before Println
func Println(v ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Println(v...)
}

// Printf is doing custom output validation before Printf
func Printf(format string, v ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Printf(format, v...)
}

// Fatalln is same as log.Fatalln
func Fatalln(v ...interface{}) {
	oLog.Fatalln(v...)
}

// Fatal is same as log.Fatal
func Fatal(v ...interface{}) {
	oLog.Fatal(v...)
}

// SetOutput is a way to setup outout to different one
func SetOutput(w io.Writer) {
	isSetOutput = true
	customLogOutput = w
	if w != nil {
		oLog.SetOutput(w)
	}
}
