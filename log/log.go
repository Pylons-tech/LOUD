package log

import (
	"io"

	oLog "github.com/sirupsen/logrus"
)

var isSetOutput = false
var customLogOutput io.Writer = nil

func init() {
	oLog.SetFormatter(&oLog.JSONFormatter{})
	oLog.SetReportCaller(true)

	oLog.SetLevel(oLog.TraceLevel)

	// oLog.Trace("Something very low level.")
	// oLog.Debug("Useful debugging information.")
	// oLog.Info("Something noteworthy happened!")
	// oLog.Warn("You should probably take a look at this.")
	// oLog.Error("Something failed but I'm not quitting.")
	// oLog.Fatal("Bye.")
	// oLog.Panic("I'm bailing.")
}

// Trace is function to replicate logrus's Trace
func Trace(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Trace(args...)
}

// Debug is function to replicate logrus's Debug
func Debug(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Debug(args...)
}

// Print is function to replicate logrus's Print
func Print(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
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
	oLog.Warn(args...)
}

// Warning is function to replicate logrus's Warning
func Warning(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Warning(args...)
}

// Error is function to replicate logrus's Error
func Error(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Error(args...)
}

// Fatal is function to replicate logrus's Fatal
func Fatal(args ...interface{}) {
	oLog.Fatal(args...)
}

// Panic is function to replicate logrus's Panic
func Panic(args ...interface{}) {
	oLog.Panic(args...)
}

// Printf family functions

// Tracef is function to replicate logrus's Tracef
func Tracef(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Tracef(format, args...)
}

// Debugf is function to replicate logrus's Debugf
func Debugf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
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
	oLog.Printf(format, args...)
}

// Warnf is function to replicate logrus's Warnf
func Warnf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Warnf(format, args...)
}

// Warningf is function to replicate logrus's Warningf
func Warningf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Warningf(format, args...)
}

// Errorf is function to replicate logrus's Errorf
func Errorf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Errorf(format, args...)
}

// Fatalf is function to replicate logrus's Fatalf
func Fatalf(format string, args ...interface{}) {
	oLog.Fatalf(format, args...)
}

// Panicf is function to replicate logrus's Panicf
func Panicf(format string, args ...interface{}) {
	oLog.Panicf(format, args...)
}

// Println family functions

// Traceln is function to replicate logrus's Traceln
func Traceln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Traceln(args...)
}

// Debugln is function to replicate logrus's Debugln
func Debugln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
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
	oLog.Warnln(args...)
}

// Warningln is function to replicate logrus's Warningln
func Warningln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.Warningln(args...)
}

// Errorln is function to replicate logrus's Errorln
func Errorln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
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

// SetOutput is a way to setup outout to different one
func SetOutput(w io.Writer) {
	isSetOutput = true
	customLogOutput = w
	if w != nil {
		oLog.SetOutput(w)
	}
}
