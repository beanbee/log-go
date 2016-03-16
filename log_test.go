/* log_test.go - test for log.go */
/*
   modification history
   --------------------
   2016/02/03, by Chen Jian, create
   2016/02/04, by Chen Jian, add setter test
   2016/03/16, by Chen Jian, add debugMode test
*/

package log

import (
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	logger, err := NewLogger("test").SetLogDir("./log").EnableWf(true).SetDebugMode(true).Init()
	if err != nil {
		t.Error("log.Init() fail")
	}

	logger.Warn("warning msg")
	logger.Info("info msg")
	logger.Debug("debug msg")
	logger.Error("error msg")
	logger.Close()

	time.Sleep(100 * time.Millisecond)

	// delete temp log directory
	os.RemoveAll("./log")
}
