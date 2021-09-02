package service

import (
	"context"
	"githup.com/htl/tcclienttest/internal/appserver/store"
	v1 "githup.com/htl/tcclienttest/internal/pkg/meta/v1"
)

type TerminalSrv interface {
	GetStoreMsg(ctx context.Context) (v1.Terminal, error)
	Execute(ctx context.Context, DeviceId int64, RunFunc func() error) error
}

type TerminalService struct {
	store store.Factory
}

func (t TerminalService) GetStoreMsg(ctx context.Context) (v1.Terminal, error) {
	panic("implement me")
}

func (t TerminalService) Execute(ctx context.Context, DeviceId int64, RunFunc func() error) error {
	panic("implement me")
}

func newTerminals(srv *service) TerminalSrv {
	return &TerminalService{srv.store}

}
