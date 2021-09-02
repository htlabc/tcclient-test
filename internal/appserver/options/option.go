package options

import (
	"githup.com/htl/tcclienttest/internal/pkg/options"
	"githup.com/htl/tcclienttest/pkg/app"
	"githup.com/htl/tcclienttest/pkg/log"
)

//appserver option

type Options struct {
	GenericServerRunOptions *options.ServerRunOptions `json:"server"   mapstructure:"server"`
	HttpOptions             *options.HttpOptions      `json:"http"    mapstructure:"http"`
	MySQLOptions            *options.MySQLOptions     `json:"mysql"    mapstructure:"mysql"`
	Log                     *log.Options              `json:"log"      mapstructure:"log"`
}

func (o Options) Flags() (fss app.NamedFlagSets) {
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.HttpOptions.AddFlags(fss.FlagSet("http"))
	o.Log.AddFlags(fss.FlagSet("logs"))
	return fss
}

func (o Options) Validate() []error {
	var errs []error
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.HttpOptions.Validate()...)
	return errs
}

func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: options.NewServerRunOptions(),
		HttpOptions:             options.NewHttpOptions(),
		MySQLOptions:            options.NewMySQLOptions(),
		Log:                     log.NewOptions(),
	}

	return &o
}
