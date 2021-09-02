package mysql

import (
	v1 "githup.com/htl/tcclienttest/internal/pkg/meta/v1"
	"gorm.io/gorm"
)

type Images struct {
	db *gorm.DB
}

type ImageInfo struct {
	imageid int `json:"id" gorm:"type:INT(11);NOT NULL;"`
}

func (i Images) List(imageid int) (*v1.Image, error) {
	imageinfo := &ImageInfo{}
	err := i.db.Table("osimages").Where("id=?", imageid).First(imageinfo).Error
	return &v1.Image{}, err
}

func (i Images) Delete(imageid int) error {
	imageinfo := &ImageInfo{}
	err := i.db.Delete(imageinfo).Error
	return err
}

func newImages(ds *datastore) *Images {
	return &Images{db: ds.db}
}
