package model

import (
	"github.com/goradd/goradd/pkg/orm/db"
)

func Database() db.DatabaseI {
	return db.GetDatabase("goradd")
}
