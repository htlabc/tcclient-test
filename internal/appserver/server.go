package appserver

import (
	"githup.com/htl/tcclienttest/internal/appserver/config"
	"githup.com/htl/tcclienttest/internal/appserver/store"
	"githup.com/htl/tcclienttest/internal/appserver/store/mysql"
	"githup.com/htl/tcclienttest/internal/pkg/options"
	"githup.com/htl/tcclienttest/internal/pkg/server"
)

type appServer struct {
	httpAppServer    *httpAppServer
	genericAPIServer *server.GenericAPIServer
}

type preparedAPIServer struct {
	*appServer
}

//构建 extra config
func buildExtraConfig(cfg *config.Config) (*ExtraConfig, error) {
	return &ExtraConfig{
		Addr:         "",
		MysqlOptions: *cfg.MySQLOptions,
		// etcdOptions:      cfg.EtcdOptions,
	}, nil
}

func createAppServer(cfg *config.Config) (*appServer, error) {
	//gs := shutdown.New()
	//gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := buildGenericConfig(cfg)
	//if err != nil {
	//	return nil, err
	//}

	extraConfig, err := buildExtraConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	//if err != nil {
	//	return nil, err
	//}
	extraServer, err := extraConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &appServer{
		httpAppServer:    extraServer,
		genericAPIServer: genericServer,
	}

	return server, nil
}

func (c *ExtraConfig) complete() *completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedExtraConfig{c}
}

type completedExtraConfig struct {
	*ExtraConfig
}

type ExtraConfig struct {
	Addr         string
	MaxMsgSize   int
	MysqlOptions options.MySQLOptions
}

func (c *completedExtraConfig) New() (*httpAppServer, error) {
	storeIns, _ := mysql.GetMySQLFactoryOr(&c.MysqlOptions)
	// storeIns, _ := etcd.GetEtcdFactoryOr(c.etcdOptions, nil)
	store.SetClient(storeIns)
	return &httpAppServer{c.Addr, nil}, nil
}

func buildGenericConfig(cfg *config.Config) (genericConfig *server.Config, lastErr error) {
	genericConfig = server.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func (s *appServer) PrepareRun() preparedAPIServer {
	//初始化httpserver gin engine

	if s.genericAPIServer.Engine != nil {

		initRouter(s.genericAPIServer.Engine)

		s.httpAppServer.Engine = s.genericAPIServer.Engine
	}

	//s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
	//	mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
	//	if mysqlStore != nil {
	//		return mysqlStore.Close()
	//	}
	//
	//	s.gRPCAPIServer.Close()
	//	s.genericAPIServer.Close()
	//
	//	return nil
	//}))

	return preparedAPIServer{s}
}

func (s preparedAPIServer) Run() error {

	// start shutdown managers
	//if err := s.gs.Start(); err != nil {
	//	log.Fatalf("start shutdown manager failed: %s", err.Error())
	//}

	return s.httpAppServer.Run()
}
