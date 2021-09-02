package appserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type httpAppServer struct {
	address string
	*gin.Engine
}

func (h *httpAppServer) Run() error {
	//listen, err := net.Listen("tcp", s.address)
	//if err != nil {
	//	log.Fatalf("failed to listen: %s", err.Error())
	//}

	//go func() {
	//	if err := s.Serve(listen); err != nil {
	//		log.Fatalf("failed to start grpc server: %s", err.Error())
	//	}
	//}()
	//
	//log.Infof("start grpc server at %s", s.address)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     h, //imageServerRouter.InitRouter(),
		ReadTimeout: 50 * time.Second,
		//WriteTimeout: 6 * time.Second,
		//MaxHeaderBytes: maxHeaderBytes,
	}
	err := server.ListenAndServe()
	return err
}

func (h *httpAppServer) Close() {
	//s.GracefulStop()
	//log.Infof("GRPC server on %s stopped", s.address)
}
