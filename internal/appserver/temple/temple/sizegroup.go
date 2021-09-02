package temple

import "github.com/go-pkgz/syncs"

func newSizegroup(sizegroup int) *syncs.SizedGroup {
	return syncs.NewSizedGroup(sizegroup)
}
