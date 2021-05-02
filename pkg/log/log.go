package log

import "log"

const (
	_info ="INFO:"
	_error ="ERROR:"
	_warn ="WARN:"

	maxLogLevel = 3
)

// logLevel sets the verbosity of the logging.
// 1 is Error
// 2 is Warn (including 1)
// 3 is Info (including 2)
var logLevel int

// Initialize the package with a log `level`
func Initialize(level int) {
	if level > maxLogLevel {
		level = maxLogLevel
	}

	// set logLevel only if it is not set already
	if logLevel == 0{
		logLevel = level
	} else {
		Warn("Package log Initialized more than once, log level remains unchanged. Level:",
			logLevel)
	}
}

func isInitialized() bool {
	if logLevel == 0 {
		return false
	}
	return true
}

func Error(v ...interface{}){
	if isInitialized() && logLevel >=1 {
		log.Println(_error, v)
	}

}

func Warn(v ...interface{}) {
	if isInitialized() && logLevel >= 2 {
		log.Println(_warn, v)
	}
}

func Info(v ...interface{}){
	if isInitialized() && logLevel >=3 {
		log.Println(_info, v)
	}
}


