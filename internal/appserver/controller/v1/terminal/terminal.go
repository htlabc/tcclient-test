package terminal

import (
	"githup.com/htl/tcclienttest/internal/appserver/service"
	"githup.com/htl/tcclienttest/internal/appserver/store"
)

type TerminalController struct {
	srv service.Service
}

// NewUserController creates a user handler.
func NewTerminalController(store store.Factory) *TerminalController {
	return &TerminalController{
		srv: service.NewService(store),
	}
}
