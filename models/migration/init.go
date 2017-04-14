package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rssllyn/sms/conf"
	"github.com/rssllyn/sms/models"
	"log"
)

func main() {
	db, err := gorm.Open("mysql", conf.Conf["db_connection_string"].(string))
	if err != nil {
		log.Fatalln("unable to connect to database", err)
	}

	db.DropTableIfExists(
		&models.Message{},
	)

	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8").CreateTable(
		&models.Message{},
	)
}
