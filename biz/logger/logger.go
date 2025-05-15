package logger

import (
	"core/conf"
	"io"
	"os"
	"path"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var c = conf.GetConf()

func ImplyZapLogger() {
	// 可定制的输出目录。
	logFilePath := c.LOG.Path
	if err := os.MkdirAll(logFilePath, 0o777); err != nil {
		panic(err)
	}
	// 将文件名设置为日期
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)

	logger := hertzzap.NewLogger(hertzzap.WithZapOptions(zap.AddCaller(), zap.AddCallerSkip(3)))
	// Provides compression and deletion
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    c.LOG.MaxSize,    // A file can be up to 20M.
		MaxBackups: c.LOG.MaxBackups, // Save up to 5 files at the same time.
		MaxAge:     c.LOG.MaxAge,     // A file can exist for a maximum of 10 days.
		Compress:   c.LOG.Compress,   // Compress with gzip.
	}

	// logger.SetOutput(lumberjackLogger)
	logger.SetLevel(conf.LogLevel())
	// if you want to output the log to the file and the stdout at the same time, you can use the following codes

	fileWriter := io.MultiWriter(lumberjackLogger, os.Stdout)
	logger.SetOutput(fileWriter)
	hlog.SetLogger(logger)
}
