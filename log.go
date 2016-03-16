/* log.go - encapsulation for log4go    */
/*
   modification history
   --------------------
   2016/02/03, by Chen Jian, create
   2016/02/04, by Chen Jian, add getter and setter
   2016/02/21, by Chen Jian, add string.trimspace for progName
   2016/03/16, by Chen Jian, set log level to 'log4go.DEBUG' in debugMode
*/

/*
DESCRIPTION:
log: encapsulation for log4go

Usage:
    import log "github.com/beanbee/log-go"

    // Two log files will be generated in ./log:
    // test.log, and test.log.wf (for log > warn)
    // The log will rotate, and there is support for backup count
    logger ,err := log.NewLogger("test").SetLogDir("./log").Init()

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

// default log format
const DEFAULT_LOG_FORMAT = `[%D %T] [%L] %M`

/*
DESCRIPTION: struct "Logger" - log4go encapsulation

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
  - debugMode: use log4go.FORMAT_DEFAULT_WITH_PID instead of DEFAULT_LOG_FORMAT
*/
type Logger struct {
	log4go.Logger
	progName    string
	level       string
	logDir      string
	when        string
	backupCount int
	hasStdOut   bool
	enableWf    bool
	debugMode   bool
}

// initial new logger - use default values
func NewLogger(progName string) *Logger {
	return &Logger{
		progName:    progName,
		level:       "INFO",
		logDir:      "log",
		when:        "MIDNIGHT",
		backupCount: 7,
		hasStdOut:   false,
		enableWf:    true,
		debugMode:   false,
	}
}

// access private variable through getter/setter function
func (l *Logger) SetLevel(level string) *Logger {
	l.level = level
	return l
}

func (l *Logger) GetLevel() string {
	return l.level
}

func (l *Logger) SetLogDir(logDir string) *Logger {
	l.logDir = logDir
	return l
}

func (l *Logger) GetLogDir() string {
	return l.logDir
}

func (l *Logger) SetWhen(when string) *Logger {
	l.when = when
	return l
}

func (l *Logger) GetWhen() string {
	return l.when
}

func (l *Logger) SetBackupCount(days int) *Logger {
	l.backupCount = days
	return l
}

func (l *Logger) GetBackupCount() int {
	return l.backupCount
}

func (l *Logger) EnableWf(useWf bool) *Logger {
	l.enableWf = useWf
	return l
}

func (l *Logger) GetEnableWf() bool {
	return l.enableWf
}

func (l *Logger) SetDebugMode(debug bool) *Logger {
	l.debugMode = debug
	return l
}

func (l *Logger) GetDebugMode() bool {
	return l.debugMode
}

func (l *Logger) SetStdOut(useStd bool) *Logger {
	l.hasStdOut = useStd
	return l
}

func (l *Logger) GetStdOutMode() bool {
	return l.hasStdOut
}

// logDirCreate(): check and create dir if nonexist
func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// create directory
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
		// for log file of warning, error, critical
		fileName = filepath.Join(logDir, strings.TrimSpace(progName)+".log.wf")
	} else {
		// for log file of all log
		fileName = filepath.Join(logDir, strings.TrimSpace(progName)+".log")
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

RETURNS:
    *Logger, nil - if succeed
    nil, error   - if fail
*/
func (l *Logger) Init() (*Logger, error) {
	// check when
	if !log4go.WhenIsValid(l.when) {
		return nil, fmt.Errorf("invalid value of when: %s", l.when)
	}

	// check, and create dir if nonexist
	if err := logDirCreate(l.logDir); err != nil {
		log4go.Error("Init(), in logDirCreate(%s)", l.logDir)
		return nil, err
	}

	// convert level from string to log4go level
	level := stringToLevel(l.level)
	if l.GetDebugMode() {
		level = log4go.DEBUG
	}

	// create logger
	logger := make(log4go.Logger)

	// create writer for stdout
	if l.hasStdOut {
		logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	// set logger format
	logFormat := func(enableDebug bool) string {
		if enableDebug {
			return log4go.FORMAT_DEFAULT_WITH_PID
		}
		return DEFAULT_LOG_FORMAT
	}(l.debugMode)

	// create file writer for all log
	fileName := filenameGen(l.progName, l.logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, l.when, l.backupCount)
	if logWriter == nil {
		return nil, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	logWriter.SetFormat(logFormat)
	logger.AddFilter("log", level, logWriter)

	// create file writer for warning and fatal log
	if l.enableWf {
		fileNameWf := filenameGen(l.progName, l.logDir, true)
		logWriter = log4go.NewTimeFileLogWriter(fileNameWf, l.when, l.backupCount)
		if logWriter == nil {
			return nil, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileNameWf)
		}
		logWriter.SetFormat(logFormat)
		logger.AddFilter("log_wf", log4go.WARNING, logWriter)
	}

	// set Logger
	l.Logger = logger

	return l, nil
}
