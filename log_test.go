/* log_test.go - test for log.go */
/*
   modification history
   --------------------
   2016/02/03, by beanbee, create
   2016/02/04, by beanbee, add setter test
   2016/03/16, by beanbee, add debugMode test
*/

package log

import (
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	logger, err := NewLogger("test").SetLogDir("./log").EnableWf(true).EnableDebug(true).Init()
	if err != nil {
		t.Error("log.Init() fail")
	}

	logger.Warn("warning msg")
	logger.Info("info msg")
	logger.Debug("debug msg")
	logger.Error("error msg")
	for i := 0; i < 999; i++ {
		logger.Warn("%d warning msg", i)
		logger.Info("%d info msg", i)
		logger.Debug("%d debug msg", i)
		logger.Error("%d error msg", i)
	}

	time.Sleep(100 * time.Millisecond)
	logger.Close()

	// delete temp log directory
	os.RemoveAll("./log")
}
