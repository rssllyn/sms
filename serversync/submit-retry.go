package serversync

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rssllyn/sms/conf"
	"github.com/rssllyn/sms/models"
	"log"
)

func SubmitRetry() {
	db, err := gorm.Open("mysql", conf.Conf["db_connection_string"].(string))
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	var messages []*models.Message
	ret := db.Find(&messages, "submited=?", false)
	if ret.Error != nil {
		log.Println(fmt.Sprintf("failed to load messages from database", err))
	}
	for _, msg := range messages {
		err := SubmitMessage(msg.Message, msg.SenderAddress, msg.ReceiverAddress, msg.SenderAddressType, msg.ReceiverAddressType, msg.ReceivedTime)
		if err == nil {
			msg.Submited = true
			db.Save(msg)
		} else {
			log.Printf("failed to submit message: %#v", err)
		}

	}
}
