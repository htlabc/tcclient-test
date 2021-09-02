package temple

import "golang.org/x/sync/errgroup"

func newErrgroup() errgroup.Group {
	return errgroup.Group{}
}
