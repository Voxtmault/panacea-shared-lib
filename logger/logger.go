package logger

import (
	"panacea-shared-lib/config"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	serverLogger *lumberjack.Logger
	errorLogger  *lumberjack.Logger
)

func InitLogger(conf *config.LoggingConfig) error {
	serverLogger = &lumberjack.Logger{
		// Log path
		Filename: conf.ServerLogPath,
		// Log size MB
		MaxSize: conf.LogMaxSize,
		// Backup count
		MaxBackups: conf.LogMaxBackup,
		// expire days
		MaxAge: conf.LogMaxAge,
		// gzip compress
		Compress: conf.LogCompress,
	}
	errorLogger = &lumberjack.Logger{
		// Log path
		Filename: conf.ErrLogPath,
		// Log size MB
		MaxSize: conf.LogMaxSize,
		// Backup count
		MaxBackups: conf.LogMaxBackup,
		// expire days
		MaxAge: conf.LogMaxAge,
		// gzip compress
		Compress: conf.LogCompress,
	}

	return nil
}

func GetServerLogger() *lumberjack.Logger {
	return serverLogger
}

func GetErrorLogger() *lumberjack.Logger {
	return errorLogger
}
