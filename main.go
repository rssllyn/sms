package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rssllyn/sms/card"
	"github.com/rssllyn/sms/commandsender"
	"github.com/rssllyn/sms/conf"
	"github.com/rssllyn/sms/resulthandler"
	"github.com/rssllyn/sms/serversync"
	_ "github.com/rssllyn/sms/signin"
	"github.com/tarm/serial"
	"log"
)

func main() {
	config := &serial.Config{Name: "/dev/zigbeegateway", Baud: 115200}
	p, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	go commandsender.Start(p)
	go resulthandler.Start(p)
	go serversync.SyncMessageWithServer()

	initializeModem(p)

	select {}
}

// 初始化设置短信模块
func initializeModem(p *serial.Port) {
	// send this command as first command after modem is initialized, to make the moderm more stable
	commandsender.Empty()

	commandsender.SetEnableCommandEcho(false)
	commandsender.SetPreferredSMSMessageStorage("SM", "SM", "SM")
	commandsender.GetMessageStorageInfo()
	// clear SMS message if storage is almost full
	if card.Info.MaxMessageStorage-card.Info.CurrentlyStoredMessage < int(conf.Conf["message_storage_threshold"].(float64)) {
		commandsender.MessageDeleteAll()
	}

	commandsender.SetShowMessageHeader(true)
	commandsender.SetMessageFormat(commandsender.MESSAGE_FORMAT_PDU)

	// let the received SMS message be stored first, we will send another command to read it, and delete it afterwards if neccessary
	// a better way is to be nofified by an unsolicited result with all the information we need, but I cannot figure out how to prevent the modem to storing messages
	commandsender.SetNewSMSMessageIndication(0, 1, 0, 0, 0)
	commandsender.GetSubscriberNumber()
}
