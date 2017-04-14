package localstore

import (
	"github.com/jinzhu/gorm"
	"github.com/rssllyn/sms/conf"
	"github.com/rssllyn/sms/models"
	"time"
)

// store received message locally in a mysql database, the database can be initilized by models/migration package
func StoreMessage(message string, senderAddress, receiverAddress string, senderAddressType, receiverAddressType int, receivedTime time.Time, submited bool) error {
	if conf.Conf["saved_to_db"].(bool) == false {
		return nil
	}
	msg := models.Message{
		ReceiverAddress:     receiverAddress,
		ReceiverAddressType: receiverAddressType,
		SenderAddress:       senderAddress,
		SenderAddressType:   senderAddressType,
		ReceivedTime:        receivedTime,
		Message:             message,
		Submited:            submited,
	}
	db, err := gorm.Open("mysql", conf.Conf["db_connection_string"].(string))
	if err != nil {
		return err
	}
	ret := db.Create(&msg)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}
