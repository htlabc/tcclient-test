package temple

import (
	"context"
	"github.com/mdlayher/schedgroup"
)

func newSchedgroup() *schedgroup.Group {
	return schedgroup.New(context.Background())
}
