package models

import (
	"time"
)

// the SMS message DB model that is to be stored locally in mysql
type Message struct {
	BaseModel
	ReceiverAddress     string    // receiver number, that is, the subscriber number
	ReceiverAddressType int       //
	SenderAddress       string    // sender number
	SenderAddressType   int       // type of sender number
	ReceivedTime        time.Time // the time when this message is received
	Message             string    // message contnet
	Submited            bool      // whether the message has bee uploaded to server via api
}
