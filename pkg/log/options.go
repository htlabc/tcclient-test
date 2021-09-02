package log

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"time"
)

type Options struct {
	OutputPath        string        `json:"output-path"       mapstructure:"output-path"`
	EbableStd         bool          `json:"ebable-std"     mapstructure:"enable-std"`
	ErrorOutputPaths  []string      `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level             logrus.Level  `json:"level"              mapstructure:"level"`
	Format            string        `json:"format"             mapstructure:"format"`
	DisableCaller     bool          `json:"disable-caller"     mapstructure:"disable-caller"`
	RotationCount     uint          `json:"rotation-count"    mapstructure:"rotation-count"`
	RotationTime      time.Duration `json:"rotation-time"     mapstructure:"rotation-time" `
	DisableStacktrace bool          `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool          `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool          `json:"development"        mapstructure:"development"`
	Name              string        `json:"name"               mapstructure:"name"`
}

func (o *Options) Validate() []error {
	var errs []error

	return errs
}

func (o Options) AddFlags(set *pflag.FlagSet) {

}

// NewOptions creates a Options object with default parameters.
func NewOptions() *Options {
	return &Options{
		Level:             logrus.InfoLevel,
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            "",
		EnableColor:       false,
		Development:       false,
		OutputPath:        "./tcclient-test.log",
		EbableStd:         true,
		RotationCount:     2,
		RotationTime:      2 * time.Hour,
		ErrorOutputPaths:  []string{"stderr"},
	}
}
