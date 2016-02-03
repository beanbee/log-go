/* log.go - encapsulation for log4go    */
/*
   modification history
   --------------------
   2016/02/03, by Chen Jian, create
*/

/*
DESCRIPTION
log: encapsulation for log4go

Usage:
    import log "github.com/beanbee/log-go"

    // Two log files will be generated in ./log:
    // test.log, and test.wf.log(for log > warn)
    // The log will rotate, and there is support for backup count
    logger ,err := log.Init("test", "INFO", "./log", true, "midnight", 5)

    logger.Warn("warn msg")
    logger.Info("info msg")

    // it is required, to work around bug of log4go
    time.Sleep(100 * time.Millisecond)
*/

package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/log4go"
)

type Logger struct{ log4go.Logger }

// private logger
var Plog Logger
var initialized bool = false

// log format - DEFAULT: log4go.FORMAT_DEFAULT
var logFormat string = "[%D %T] [%L] %M"

// logDirCreate(): check and create dir if nonexist
func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// filenameGen(): generate filename
func filenameGen(progName, logDir string, isErrLog bool) string {
	var fileName string
	if isErrLog {
		/* for log file of warning, error, critical  */
		fileName = filepath.Join(logDir, progName+".log.wf")
	} else {
		/* for log file of all log  */
		fileName = filepath.Join(logDir, progName+".log")
	}

	return fileName
}

/* convert level in string to log4go level  */
func stringToLevel(str string) log4go.LevelType {
	var level log4go.LevelType

	str = strings.ToUpper(str)
	switch str {
	case "DEBUG":
		level = log4go.DEBUG
	case "TRACE":
		level = log4go.TRACE
	case "INFO":
		level = log4go.INFO
	case "WARNING":
		level = log4go.WARNING
	case "ERROR":
		level = log4go.ERROR
	case "CRITICAL":
		level = log4go.CRITICAL
	default:
		level = log4go.INFO
	}
	return level
}

/*
Init - initialize log lib

PARAMS:
  - progName: program name. Name of log file will be progName.log
  - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
  - logDir: directory for log. It will be created if noexist
  - hasStdOut: whether to have stdout output
  - when:
      "M", minute
      "H", hour
      "D", day
      "MIDNIGHT", roll over at midnight
  - backupCount: If backupCount is > 0, when rollover is done, no more than
      backupCount files are kept - the oldest ones are deleted.
  - enableWf: using extra log file for 'warning, error, critical' level msg

RETURNS:
    *Logger, nil - if succeed
    nil, error   - if fail
*/
func Init(progName string, levelStr string, logDir string,
	hasStdOut bool, when string, backupCount int, enableWf bool) (Logger, error) {
	/* check when   */
	if !log4go.WhenIsValid(when) {
		return Logger{}, fmt.Errorf("invalid value of when: %s", when)
	}

	/* check, and create dir if nonexist    */
	if err := logDirCreate(logDir); err != nil {
		log4go.Error("Init(), in logDirCreate(%s)", logDir)
		return Logger{}, err
	}

	/* convert level from string to log4go level    */
	level := stringToLevel(levelStr)

	/* create logger    */
	logger := make(log4go.Logger)

	/* create writer for stdout */
	if hasStdOut {
		logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	/* create file writer for all log   */
	fileName := filenameGen(progName, logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return Logger{}, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	// logWriter.SetFormat(log4go.LogFormat)
	logWriter.SetFormat(logFormat)
	logger.AddFilter("log", level, logWriter)

	/* create file writer for warning and fatal log */
	if enableWf {
		fileNameWf := filenameGen(progName, logDir, true)
		logWriter = log4go.NewTimeFileLogWriter(fileNameWf, when, backupCount)
		if logWriter == nil {
			return Logger{}, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileNameWf)
		}
		logWriter.SetFormat(logFormat)
		logger.AddFilter("log_wf", log4go.WARNING, logWriter)
	}

	return Logger{logger}, nil
}

/*
InitPrivate - initialize private logger, for situation when using extra logger inside package

Usage:
    log.Init("test", "INFO", "./log", true, "midnight", 5)
    log.Plog.Warn("warn msg")
    log.Plog.Info("info msg")

PARAMS:
    same as Init()

RETURNS:
    nil, if succeed
    error, if fail
*/
func InitPrivate(progName string, levelStr string, logDir string,
	hasStdOut bool, when string, backupCount int, enableWf bool) error {
	var err error
	if initialized {
		fmt.Println("Initialized Already")
		// return errors.New("Initialized Already")
		return nil
	}

	// init private logger
	Plog, err = Init(progName, levelStr, logDir, hasStdOut, when, backupCount, enableWf)
	if err != nil {
		return err
	}

	initialized = true
	return nil
}
