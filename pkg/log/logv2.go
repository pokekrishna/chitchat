package log

import "log"

const (
	_info ="INFO:"
	_error ="ERROR:"
	_warn ="WARN:"

	MaxLogLevel = 3
)

var defaultLogger *logger

type metaLogger interface{
	basePrinter(v ...interface{})
}

// A logger is the actual defaultLogger. level sets the verbosity of the logging.
// logger is meant to be private (not exposed) because the recommended usage of the
// logging functions are <pkg_name>.<method_name>(v ...).
// For instance , log.Info("message")
//
// level 1 is Error
// level 2 is Warn (including 1)
// level 3 is Info (including 2)
type logger struct {
	level int
	metaLogger
}

// implMetaLogger is a simple implementation of the metaLogger interface
type implMetaLogger struct {}

func (l *logger) isInitialized() bool {
	if l.level == 0 {
		return false
	}
	return true
}

func (l *logger) Info(v ...interface{}) {
	if l.isInitialized() && l.level >=3 {
		l.basePrinter(_info, v)
	}
}

func (m *implMetaLogger) basePrinter (v ...interface{}){
	log.Println(v...)
}

// TODO : expose the methods as all logging functions to be used as <pkg_name>.<function_name>()
	// TODO: would using an interface for exposed functions be useful, as a contract?

//Initialize the package with a log `level`
func Initialize(level int) {
	if level < 0{
		level = 0
	} else if level > MaxLogLevel {
		level = MaxLogLevel
	}

	// set logLevel only if it is not set already
	if defaultLogger == nil || defaultLogger.level == 0{
		defaultLogger = &logger{
			level:      level,
			metaLogger: &implMetaLogger{},
		}
	} //else {
//		Warn("Package log Initialized more than once, log level remains unchanged. Level:",
//			logLevel)
//	}
	// TODO: Complete the else branch ^^ once Warn() is onboarded, and function vs method is decided
}

func GetLevel() int{
	if defaultLogger != nil {
		return defaultLogger.level
	} else{
		return 0
	}
}

// ResetForTests resets the package as if Initialize() was never called.
// Convenience method for testing. This should only be called from tests.
func ResetForTests(){
	defaultLogger = nil
}