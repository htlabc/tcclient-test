package appserver

import "githup.com/htl/tcclienttest/internal/appserver/config"

//run.go->server.go

func Run(cfg *config.Config) error {
	server, err := createAppServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}
