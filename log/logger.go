package log

import (
	log "github.com/sirupsen/logrus"
)

//InitLogger logger
func InitLogger(serverName string) *log.Entry {
	serverLogger := log.WithFields(log.Fields{
		"serverName": serverName,
	})
	return serverLogger
}
