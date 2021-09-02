package appserver

import (
	"githup.com/htl/tcclienttest/internal/appserver/config"
	"githup.com/htl/tcclienttest/internal/appserver/options"
	"githup.com/htl/tcclienttest/pkg/app"
	"githup.com/htl/tcclienttest/pkg/log"
)

//使用多选项模式
func NewApp(baseName string) *app.App {
	opts := options.NewOptions()
	return app.NewApp("test", baseName,
		app.WithOptions(opts),
		app.WithRunFunc(run(opts)),
	)

}

func run(opts *options.Options) app.RunFunc {

	return func(basename string) error {
		log.Init(opts.Log)
		//根据option生成config
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)

	}

}
