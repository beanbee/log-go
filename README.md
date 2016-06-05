# log-go
Encapsulation of log4go.

Usage:
```Go
    import log "github.com/beanbee/log-go"

    // Generate two log files in "./log" directory (xx.log, xx.log.wf)
    // Log filter: -> test.log.wf (log.level > warn)
    // Log rotate: support for backup count
    logger ,err := log.NewLogger("test").SetLogDir("./log").Init()
	
    logger.Warn("warn msg")
    logger.Info("info msg")
    logger.Debug("debug msg") // only seen in debugMode which could be set by 'SetDebugMode(true)'

    // it is required, to work around the bug of log4go
    time.Sleep(100 * time.Millisecond)
```

