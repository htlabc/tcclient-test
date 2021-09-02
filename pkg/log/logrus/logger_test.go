package logrus

import (
	"context"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

func TestLogrusFormatter1(t *testing.T) {
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	logrus.Info("info msg")

}

func TestLogrusFormatter2(t *testing.T) {
	//调用logrus.SetReportCaller(true)设置在输出日志中添加文件名和方法信息：

	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	})
	logrus.SetReportCaller(true)

	writers := make([]io.Writer, 0)

	file1, _ := os.Create("test1.log")
	file2, _ := os.Create("test2.log")
	file3, _ := os.Create("test3.log")
	file4 := os.Stdout

	writers = append(writers, file1, file2, file3, file4)

	logrus.SetOutput(io.MultiWriter(writers...))

	logrus.Info("test")

}

func TestLogrusHook(t *testing.T) {
	lg := GetNewFieldLoggerContext("test", "d")
	lg.Logger.WithContext(context.WithValue(context.TODO(), "test", "test")).Info("????")
}

//日志分割lfshook
func newLogrusHook(logPath, moduel string) logrus.Hook {
	logrus.SetLevel(logrus.WarnLevel)

	writer := FakeNewLoggger(logPath, moduel, time.Hour*2)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true})

	// writer 生成新的log文件类型 writer  在通过new hook函数 消费 fire 函数
	// writer 是实现了writer 接口的库，在日志调用write是做预处理
	return lfsHook
}

type LogWriter struct {
	logDir           string //日志根目录地址。
	module           string //模块 名
	curFileName      string //当前被指定的filename
	curBaseFileName  string //在使用中的file
	turnCateDuration time.Duration
	mutex            sync.RWMutex
	outFh            *os.File
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if out, err := w.FakeGetWriter(); err != nil {
		return 0, errors.New("failed to fetch targer io.Write")
	} else {
		return out.Write(p)
	}
}

func (w *LogWriter) FakeGetFileName() string {
	base := time.Now().Truncate(w.turnCateDuration)
	return fmt.Sprintf("%s/%s/%s_%s", w.logDir, "2006-01-02", w.module, base.Format("15"))
}

func (w *LogWriter) FakeGetWriter() (io.Writer, error) {
	fileName := w.curBaseFileName
	//判断是否有新的文件名
	//会有新的文件名

	baseFileName := w.FakeGetFileName()
	if baseFileName != fileName {
		fileName = baseFileName
	}

	dirname := filepath.Dir(fileName)

	if err := os.MkdirAll(dirname, 0755); err != nil {
		return nil, errors.Wrapf(err, "failed to create directory %s", dirname)
	}

	fileHandler, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Errorf("failed to open file %s", err)
	}
	w.outFh.Close()
	w.outFh = fileHandler
	w.curBaseFileName = fileName
	w.curFileName = fileName

	return fileHandler, nil

}

func FakeNewLoggger(logPath, module string, duration time.Duration) *LogWriter {
	return &LogWriter{
		logDir:           logPath,
		module:           module,
		turnCateDuration: duration,
		curFileName:      "",
		curBaseFileName:  "",
	}
}

//自实现 logrus hook
func getLogger(module string) *logrus.Logger {
	//实例化
	logger := logrus.New()
	//设置输出
	logger.Out = os.Stdout
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	//设置日志格式
	//自定writer就行， hook 交给 lfshook
	logger.AddHook(newLogrusHook(`E:\工作目录\cvte\tcclient-test\pkg\log\logrus\`, module))
	logger.SetReportCaller(true)
	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	})

	logger.WithField("app", "test")

	return logger
}

//确保每次调用使用的文件都是唯一的。
func GetNewFieldLoggerContext(module, appField string) *logrus.Entry {
	logger := getLogger(module)
	return logger.WithFields(logrus.Fields{
		"app": appField,
	})
}

func TestRotatelogs(t *testing.T) {

	path := "./go.log"
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	 `WithMaxAge` 设置文件清理前的最长保存时间
	 `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		//rotatelogs.WithLinkName(path),
		//rotatelogs.WithMaxAge(time.Duration(20)*time.Second),
		rotatelogs.WithRotationCount(2),
		rotatelogs.WithRotationTime(time.Duration(10)*time.Second),
	)

	logger := logrus.New()

	logger.SetOutput(io.MultiWriter(writer, os.Stdout))
	//log.SetFormatter(&log.JSONFormatter{})

	for {
		logger.Info("hello, world!")
		time.Sleep(time.Duration(5) * time.Second)
	}

}
