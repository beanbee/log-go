/* log_test.go - test for log.go */
/*
   modification history
   --------------------
   2016/02/03, by Chen Jian, create
*/

package log

import (
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	var logger Logger
	var err error

	if logger, err = Init("test", "INFO", "./log/", true, "M", 2, true); err != nil {
		t.Error("log.Init() fail")
	}

	logger.Warn("warning msg")
	logger.Info("info msg")
	logger.Error("error msg")
	logger.Close()

	// test private logger
	InitPrivate("test", "INFO", "./log/", true, "M", 2, false)
	Plog.Info("log from private log info")
	Plog.Warn("log from private log warn")
	Plog.Error("log from private log error")

	time.Sleep(1000 * time.Millisecond)
	Plog.Close()

	// delete temp log directory
	os.RemoveAll("./log")

}
