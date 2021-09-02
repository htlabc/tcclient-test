package image

import (
	"githup.com/htl/tcclienttest/internal/appserver/service"
	"githup.com/htl/tcclienttest/internal/appserver/store"
)

type ImageController struct {
	srv service.Service
}

// NewUserController creates a user handler.
func NewImageController(store store.Factory) *ImageController {
	return &ImageController{
		srv: service.NewService(store),
	}
}
