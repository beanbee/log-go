/* log.go - encapsulation for log4go    */
/*
   modification history
   --------------------
   2016/02/03, by beanbee, create
   2016/02/04, by beanbee, add getter and setter
   2016/02/21, by beanbee, add string.trimspace for progName
   2016/03/16, by beanbee, set log level to 'log4go.DEBUG' in debugMode
   2016/06/02, by beanbee, add rotate size
   2016/08/01, by beanbee, add enableRotate
*/

/*
DESCRIPTION:
log: encapsulation for log4go

Usage:
    import log "github.com/beanbee/log-go"

    // Two log files will be generated in ./log
    // test.log, and test.log.wf (for log > warn)
    // The log will rotate, and with support for backup count
    logger ,err := log.NewLogger("test", "./log").Init()

    logger.Warn("warn msg")
    logger.Info("info msg")

    // it is required, to work around bug of log4go
    time.Sleep(100 * time.Millisecond)
*/

package log

import (
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/log4go"
)

const (
	FORMAT_WITHOUT_SRC  = "[%D %T] [%L] %M"
	DEFAULT_ROTATE_SIZE = 1024 * 1024 * 1024 // 1 GB
)

/*
DESCRIPTION: struct "Logger" - log4go encapsulation

PARAMS:
  - progName: program name. Name of log file will be progName.log
  - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
  - logDir: directory for log. It will be created if noexist
  - enableWf: using extra log file for 'warning, error, critical' level msg
  - enableStdout: whether to have stdout output
  - enableStdout: use log4go.FORMAT_DEFAULT_WITH_PID instead of DEFAULT_LOG_FORMAT
*/
type Logger struct {
	progName     string
	logDir       string
	rotateSize   int
	enableWf     bool
	enableDebug  bool
	enableStdout bool
	enableRotate bool

	log4go.Logger
}

// initial new logger - use default values
func NewLogger(progName string, logDir string) *Logger {
	return &Logger{
		progName: progName,
		logDir: func(dir string) string {
			if dir == "" {
				return "./" // convert empty dir to "./"
			}
			return dir
		}(logDir),
		enableStdout: false,
		enableWf:     false,
		enableDebug:  false,
		enableRotate: false,
	}
}

// func (l *Logger) SetRotateSize(size int) *Logger {
// 	l.rotateSize = size
// 	return l
// }

func (l *Logger) EnableRotate(enable bool) *Logger {
	l.enableRotate = enable
	if enable && l.rotateSize == 0 {
		l.rotateSize = DEFAULT_ROTATE_SIZE
	}
	return l
}

func (l *Logger) EnableStdout(enable bool) *Logger {
	l.enableStdout = enable
	return l
}

func (l *Logger) EnableWf(enable bool) *Logger {
	l.enableWf = enable
	return l
}

func (l *Logger) EnableDebug(enable bool) *Logger {
	l.enableDebug = enable
	return l
}

// check and create dir if nonexist
func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// filenameGen(): generate filename
func filenameGen(progName, logDir string, enableWf bool) string {
	var fileName string
	if enableWf {
		// for log file of warning, error, critical
		fileName = filepath.Join(logDir, strings.TrimSpace(progName)+".log.wf")
	} else {
		// for log file of all log
		fileName = filepath.Join(logDir, strings.TrimSpace(progName)+".log")
	}

	return fileName
}

/*
Init - initialize log lib

RETURNS:
    *Logger, nil - if succeed
    nil, error   - if fail
*/
func (l *Logger) Init() (*Logger, error) {
	// check, and create dir if nonexist
	if err := logDirCreate(l.logDir); err != nil {
		return nil, err
	}

	// default level: INFO
	level := log4go.INFO
	if l.enableDebug {
		level = log4go.DEBUG
	}

	// set logger format
	logFormat := func(enableDebug bool) string {
		if enableDebug {
			return log4go.FORMAT_DEFAULT
		}
		return FORMAT_WITHOUT_SRC
	}(l.enableDebug)

	// create logger
	logger := make(log4go.Logger)

	// create writer for stdout
	if l.enableStdout {
		logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	// create file writer for all log
	fileName := filenameGen(l.progName, l.logDir, false)
	logWriter := log4go.NewFileLogWriter(fileName, l.enableRotate)
	if l.enableRotate {
		logWriter.SetRotateSize(l.rotateSize)
		logWriter.SetRotateDaily(true)
	}
	logWriter.SetFormat(logFormat)
	logger.AddFilter("log", level, logWriter)

	if l.enableWf {
		fileNameWf := filenameGen(l.progName, l.logDir, true)
		logWriterWf := log4go.NewFileLogWriter(fileNameWf, l.enableRotate)
		if l.enableRotate {
			logWriterWf.SetRotateSize(l.rotateSize)
			logWriterWf.SetRotateDaily(true)
		}
		logWriterWf.SetFormat(logFormat)
		logger.AddFilter("log_wf", log4go.WARNING, logWriterWf)
	}

	// set Logger
	l.Logger = logger

	return l, nil
}
