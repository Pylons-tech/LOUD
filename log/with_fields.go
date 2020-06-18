package log

import oLog "github.com/sirupsen/logrus"

// Fields is a type to manage json based output
type Fields oLog.Fields

// Logger is a struct to manage custom logging for loud application
type Logger struct {
	fields oLog.Fields
}

// WithFields is to manage data in json format
func WithFields(fields Fields) Logger {
	return Logger{
		fields: oLog.Fields(fields),
	}
}

// Trace is function to replicate logrus's Trace
func (logger Logger) Trace(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Trace(args...)
}

// Debug is function to replicate logrus's Debug
func (logger Logger) Debug(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Debug(args...)
}

// Print is function to replicate logrus's Print
func (logger Logger) Print(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Print(args...)
}

// Info is function to replicate logrus's Info
func (logger Logger) Info(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.WithFields(logger.fields).Info(args...)
}

// Warn is function to replicate logrus's Warn
func (logger Logger) Warn(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Warn(args...)
}

// Warning is function to replicate logrus's Warning
func (logger Logger) Warning(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Warning(args...)
}

// Error is function to replicate logrus's Error
func (logger Logger) Error(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Error(args...)
}

// Fatal is function to replicate logrus's Fatal
func (logger Logger) Fatal(args ...interface{}) {
	printCallerLine()
	oLog.WithFields(logger.fields).Fatal(args...)
}

// Panic is function to replicate logrus's Panic
func (logger Logger) Panic(args ...interface{}) {
	printCallerLine()
	oLog.WithFields(logger.fields).Panic(args...)
}

// Printf family functions

// Tracef is function to replicate logrus's Tracef
func (logger Logger) Tracef(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Tracef(format, args...)
}

// Debugf is function to replicate logrus's Debugf
func (logger Logger) Debugf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Debugf(format, args...)
}

// Infof is function to replicate logrus's Infof
func (logger Logger) Infof(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.WithFields(logger.fields).Infof(format, args...)
}

// Printf is function to replicate logrus's Printf
func (logger Logger) Printf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Printf(format, args...)
}

// Warnf is function to replicate logrus's Warnf
func (logger Logger) Warnf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Warnf(format, args...)
}

// Warningf is function to replicate logrus's Warningf
func (logger Logger) Warningf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Warningf(format, args...)
}

// Errorf is function to replicate logrus's Errorf
func (logger Logger) Errorf(format string, args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Errorf(format, args...)
}

// Fatalf is function to replicate logrus's Fatalf
func (logger Logger) Fatalf(format string, args ...interface{}) {
	printCallerLine()
	oLog.WithFields(logger.fields).Fatalf(format, args...)
}

// Panicf is function to replicate logrus's Panicf
func (logger Logger) Panicf(format string, args ...interface{}) {
	printCallerLine()
	oLog.WithFields(logger.fields).Panicf(format, args...)
}

// Println family functions

// Traceln is function to replicate logrus's Traceln
func (logger Logger) Traceln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Traceln(args...)
}

// Debugln is function to replicate logrus's Debugln
func (logger Logger) Debugln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Debugln(args...)
}

// Infoln is function to replicate logrus's Infoln
func (logger Logger) Infoln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.WithFields(logger.fields).Infoln(args...)
}

// Println is function to replicate logrus's Println
func (logger Logger) Println(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	oLog.WithFields(logger.fields).Infoln(args...)
}

// Warnln is function to replicate logrus's Warnln
func (logger Logger) Warnln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Warnln(args...)
}

// Warningln is function to replicate logrus's Warningln
func (logger Logger) Warningln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Warningln(args...)
}

// Errorln is function to replicate logrus's Errorln
func (logger Logger) Errorln(args ...interface{}) {
	if isSetOutput && customLogOutput == nil {
		return
	}
	printCallerLine()
	oLog.WithFields(logger.fields).Errorln(args...)
}

// Fatalln is function to replicate logrus's Fatalln
func (logger Logger) Fatalln(args ...interface{}) {
	oLog.WithFields(logger.fields).Fatalln(args...)
}

// Panicln is function to replicate logrus's Panicln
func (logger Logger) Panicln(args ...interface{}) {
	oLog.WithFields(logger.fields).Panicln(args...)
}
