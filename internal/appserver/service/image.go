package service

import (
	"context"
	"githup.com/htl/tcclienttest/internal/appserver/store"
	v1 "githup.com/htl/tcclienttest/internal/pkg/meta/v1"
)

type ImageSrv interface {
	Get(ctx context.Context, imageid int64, opts v1.GetOptions) (*v1.Image, error)
}

type ImageService struct {
	store store.Factory
}

func (i ImageService) Get(ctx context.Context, imageid int64, opts v1.GetOptions) (*v1.Image, error) {
	return i.store.Images().List(int(imageid))
}

func newImages(srv *service) ImageSrv {
	return &ImageService{srv.store}

}
