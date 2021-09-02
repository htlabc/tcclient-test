package store

import v1 "githup.com/htl/tcclienttest/internal/pkg/meta/v1"

type ImageStore interface {
	List(imageid int) (*v1.Image, error)
	Delete(imageid int) error
}
