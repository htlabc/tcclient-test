package log

import (
	"github.com/sirupsen/logrus"
	logrus2 "githup.com/htl/tcclienttest/pkg/log/logrus"
	"sync"
)

var (
	once sync.Once
	std  *logrusLogger = New(NewOptions())
)

func Init(opts *Options) {
	once.Do(func() {
		std = New(opts)
	})
}

type Logger interface {
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
}

type logrusLogger struct {
	// NB: this looks very similar to zap.SugaredLogger, but
	// deals with our desire to have multiple verbosity levels.
	Logrus *logrus.Logger
}

func (l logrusLogger) Infof(format string, v ...interface{}) {
	l.Logrus.Infof(format, v)
}

func (l logrusLogger) Info(v ...interface{}) {
	l.Logrus.Info(v)
}

func (l logrusLogger) Warnf(format string, v ...interface{}) {
	l.Logrus.Warnf(format, v)
}

func (l logrusLogger) Errorf(format string, v ...interface{}) {
	l.Logrus.Errorf(format, v)
}

func (l logrusLogger) Debugf(format string, v ...interface{}) {
	l.Logrus.Debugf(format, v)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

func Info(v ...interface{}) {
	std.Info(v...)
}

func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	std.Warnf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func New(opts *Options) *logrusLogger {

	c := &logrus2.Config{Level: opts.Level,
		LogPath:       opts.OutputPath,
		RotationTime:  opts.RotationTime,
		RotationCount: opts.RotationCount,
		EbableStd:     opts.EbableStd,
		EnableColor:   opts.EnableColor,
		DisableCaller: opts.DisableCaller,
	}
	return &logrusLogger{logrus2.NewLogger(c)}
	//logrus2.NewLogger(opts)
}

func StdInfoLogger() *logrusLogger {
	if std == nil {
		return nil
	}
	return std
}
