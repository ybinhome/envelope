package eureka

import (
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func GetLogger() *log.Logger {
	return log.StandardLogger()
}

//
//func SetLogger(loggerLog *log.Logger) {
//	logger.SetLogger(loggerLog)
//}
//
func init() {
	// Default logger uses the go default log.
	//logger = gominlog.NewClassicMinLogWithPackageName("eureka")
	//logger.SetLevel(gominlog.Linfo)
	logger = GetLogger()
}
