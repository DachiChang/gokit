package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(logPath string, logLevel logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logLevel)
	logRotate := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    30, // megabytes
		MaxBackups: 10,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	}
	logger.SetOutput(io.MultiWriter(logRotate, os.Stdout))
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote: true,
		ForceColors:  true,
	})
	logger.Info("Logger is ready.")

	return logger
}
