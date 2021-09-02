package logrus

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"os"
	"time"

	//rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	//"githup.com/htl/tcclienttest/pkg/log"
	//"io"
	//"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Level         logrus.Level
	LogPath       string
	RotationCount uint
	RotationTime  time.Duration
	MaxAge        time.Duration
	EbableStd     bool
	EnableColor   bool
	DisableCaller bool
}

// Validate validate the options fields.

// NewLogger create a logrus logger, add hook to it and return it.
func NewLogger(config *Config) *logrus.Logger {
	logger := logrus.New()
	writer, _ := rotatelogs.New(
		config.LogPath+".%Y%m%d%H%M",
		//rotatelogs.WithMaxAge(),
		rotatelogs.WithRotationCount(config.RotationCount),
		rotatelogs.WithRotationTime(config.RotationTime),
	)

	logger.SetOutput(io.MultiWriter(writer, os.Stdout))
	if !config.EbableStd {
		logger.SetOutput(writer)
	}

	logger.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		NoColors:        config.EnableColor,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	})
	logrus.SetReportCaller(config.DisableCaller)
	logrus.SetLevel(config.Level)
	//添加一个hook点
	//可扩展。logrus 的 Hook 机制允许使用者通过 Hook 的方式，将日志分发到任意地
	//方，例如本地文件、标准输出、Elasticsearch、Logstash、Kafka 等
	return logger
}
