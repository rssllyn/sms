package models

import (
	"github.com/jinzhu/gorm"
)

func init() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tbl_" + defaultTableName
	}
}
