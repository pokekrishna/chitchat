package log

import (
	"log"
)

const (
	_info ="INFO:"
	_error ="ERROR:"
	_warn ="WARN:"

	MaxLogLevel = 3
)

var (
	defaultLogger *logger

	Info = makeLogger((*logger).Info)
	Error = makeLogger((*logger).Error)
	Warn = makeLogger((*logger).Warn)
)

// TODO: what are the justifications of having this interface? ...
// TODO: ...if metaLogger interface wasn't there and only  implMetaLogger...
// TODO: ... was there, would it cause any problem?
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

// TODO: why are methods exposed? starting with caps

func (l *logger) Info(v ...interface{}) {
	if l.isInitialized() && l.level >=3 {
		l.basePrinter(_info, v)
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.isInitialized() && l.level >=1 {
		l.basePrinter(_error, v)
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.isInitialized() && l.level >=2 {
		l.basePrinter(_warn, v)
	}
}
func (m *implMetaLogger) basePrinter (v ...interface{}){
	log.Println(v...)
}

// makeLogger takes in a method and returns an anonymous function called on defaultLogger
func makeLogger(fn func(l *logger, v ... interface{})) func(...interface{}) {
	return func(v ...interface{}) {
		fn(defaultLogger, v...)
	}
}

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
	} else {
		defaultLogger.Warn("Package log Initialized more than once, log level remains unchanged. Level:",
			defaultLogger.level)
	}
}

// TODO: warning code smell, get level is not dependent on isInitialised, why?
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