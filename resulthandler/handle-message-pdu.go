package resulthandler

import (
	"github.com/rssllyn/sms/card"
	"github.com/rssllyn/sms/localstore"
	"github.com/rssllyn/sms/pdu"
	"github.com/rssllyn/sms/serversync"
	"log"
)

func handleMessagePDU(messagePDUHex string) {
	message, err := pdu.FromHexString(messagePDUHex)
	if err != nil {
		log.Println("invalid CMT data format", err)
	}
	log.Printf("got message %#v", message)
	messageDelivery, ok := message.(*pdu.DeliveryMessage)
	if !ok {
		log.Println("not a *pdu.DeliveryMessage")
		return
	}
	log.Println("about to submit")
	submited := false
	err = serversync.SubmitMessage(messageDelivery.Message, messageDelivery.SenderAddress, card.Info.Address, messageDelivery.SenderAddressType, card.Info.AddressType, messageDelivery.Time)
	if err != nil {
		log.Println("failed to submit message to server", err)
	} else {
		submited = true
	}
	if err := localstore.StoreMessage(messageDelivery.Message, messageDelivery.SenderAddress, card.Info.Address, messageDelivery.SenderAddressType, card.Info.AddressType, messageDelivery.Time, submited); err != nil {
		log.Println("failed to store message locally", err)
	}
}
