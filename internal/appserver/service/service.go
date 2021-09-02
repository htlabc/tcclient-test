package service

import "githup.com/htl/tcclienttest/internal/appserver/store"

type Service interface {
	Images() ImageSrv
	Terminals() TerminalSrv
}

type service struct {
	store store.Factory
}

// NewService returns Service interface.
func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) Images() ImageSrv {
	return newImages(s)
}

func (s *service) Terminals() TerminalSrv {
	return newTerminals(s)
}
