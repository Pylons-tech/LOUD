package log

import (
	"io"
	oLog "log"
)

var isSetOutput = false
var customLogOutput io.Writer = nil

func Println(v ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Println(v...)
}

func Printf(format string, v ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Printf(format, v...)
}

func Fatalln(v ...interface{}) {
	oLog.Fatalln(v...)
}

func Fatal(v ...interface{}) {
	oLog.Fatal(v...)
}

func SetOutput(w io.Writer) {
	isSetOutput = true
	customLogOutput = w
	if w != nil {
		oLog.SetOutput(w)
	}
}
