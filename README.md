# log-go
Encapsulation of log4go.

Usage:
    import log "github.com/beanbee/log-go"

    // Two log files will be generated in ./log
    // Log file: test.log, and test.log.wf (log.level > warn)
    // The log will rotate, and with support for backup count
    logger ,err := log.NewLogger("test").SetLogDir("./log").Init()
	
    logger.Warn("warn msg")
    logger.Info("info msg")
    logger.Debug("debug msg") // only seen in debugMode which could be set by 'SetDebugMode(true)'

    // it is required, to work around bug of log4go
    time.Sleep(100 * time.Millisecond)

