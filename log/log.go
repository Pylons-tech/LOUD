package log

import (
	"fmt"
	"io"
	"runtime"

	oLog "github.com/sirupsen/logrus"
)

var isSetOutput = false
var customLogOutput io.Writer = nil

func init() {
	oLog.SetLevel(oLog.TraceLevel)
}

// SetOutput is a way to setup outout to different one
func SetOutput(w io.Writer) {
	isSetOutput = true
	customLogOutput = w
	if w != nil {
		oLog.SetOutput(w)
		// oLog.SetFormatter(&oLog.JSONFormatter{})
		oLog.SetFormatter(&oLog.TextFormatter{
			DisableTimestamp: true,
			DisableColors:    true,
		})
		// oLog.SetReportCaller(true)
	}
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

func printCallerLine() {
	frame := getFrame(2)
	oLog.WithFields(oLog.Fields{
		"file_line": fmt.Sprintf("%s:%d", frame.File, frame.Line),
		"func":      frame.Function,
	}).Trace("debug caller line")
}

// Trace is function to replicate logrus's Trace
func Trace(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Trace(args...)
}

// Debug is function to replicate logrus's Debug
func Debug(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Debug(args...)
}

// Print is function to replicate logrus's Print
func Print(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Print(args...)
}

// Info is function to replicate logrus's Info
func Info(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Info(args...)
}

// Warn is function to replicate logrus's Warn
func Warn(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Warn(args...)
}

// Warning is function to replicate logrus's Warning
func Warning(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Warning(args...)
}

// Error is function to replicate logrus's Error
func Error(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Error(args...)
}

// Fatal is function to replicate logrus's Fatal
func Fatal(args ...interface{}) {
	printCallerLine()
	oLog.Fatal(args...)
}

// Panic is function to replicate logrus's Panic
func Panic(args ...interface{}) {
	printCallerLine()
	oLog.Panic(args...)
}

// Printf family functions

// Tracef is function to replicate logrus's Tracef
func Tracef(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Tracef(format, args...)
}

// Debugf is function to replicate logrus's Debugf
func Debugf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Debugf(format, args...)
}

// Infof is function to replicate logrus's Infof
func Infof(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Infof(format, args...)
}

// Printf is function to replicate logrus's Printf
func Printf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Printf(format, args...)
}

// Warnf is function to replicate logrus's Warnf
func Warnf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Warnf(format, args...)
}

// Warningf is function to replicate logrus's Warningf
func Warningf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Warningf(format, args...)
}

// Errorf is function to replicate logrus's Errorf
func Errorf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Errorf(format, args...)
}

// Fatalf is function to replicate logrus's Fatalf
func Fatalf(format string, args ...interface{}) {
	printCallerLine()
	oLog.Fatalf(format, args...)
}

// Panicf is function to replicate logrus's Panicf
func Panicf(format string, args ...interface{}) {
	printCallerLine()
	oLog.Panicf(format, args...)
}

// Println family functions

// Traceln is function to replicate logrus's Traceln
func Traceln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Traceln(args...)
}

// Debugln is function to replicate logrus's Debugln
func Debugln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Debugln(args...)
}

// Infoln is function to replicate logrus's Infoln
func Infoln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Infoln(args...)
}

// Println is function to replicate logrus's Println
func Println(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Infoln(args...)
}

// Warnln is function to replicate logrus's Warnln
func Warnln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Warnln(args...)
}

// Warningln is function to replicate logrus's Warningln
func Warningln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Warningln(args...)
}

// Errorln is function to replicate logrus's Errorln
func Errorln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.Errorln(args...)
}

// Fatalln is function to replicate logrus's Fatalln
func Fatalln(args ...interface{}) {
	oLog.Fatalln(args...)
}

// Panicln is function to replicate logrus's Panicln
func Panicln(args ...interface{}) {
	oLog.Panicln(args...)
}
