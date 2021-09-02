package server

import (
	"github.com/gin-gonic/gin"
	"time"
)

type GenericAPIServer struct {
	middlewares []string
	mode        string
	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration
	Engine          *gin.Engine
	healthz         bool
	enableMetrics   bool
	enableProfiling bool
	// wrapper for gin.Engine
	//insecureServer, secureServer *http.Server
}

func initGenericAPIServer(s *GenericAPIServer) {
	// do some setup
	// s.GET(path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	//s.Setup()
	//s.InstallMiddlewares()
	//s.InstallAPIs()
}
