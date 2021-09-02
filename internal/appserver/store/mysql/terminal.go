package mysql

import "gorm.io/gorm"

type Terminals struct {
	db *gorm.DB
}

//从存储层获取终端消息
func (t Terminals) Get() (interface{}, error) {
	panic("implement me")
}

func newTerminals(ds *datastore) *Terminals {
	return &Terminals{db: ds.db}
}
