package resulthandler

import (
	"log"
)

// 收到短信通知, +CNMI=0,2,0,0,0
func cmt(data string) {
	matches := RESULT_DATA_CMT.FindStringSubmatch(data)
	if len(matches) != 2 {
		log.Println("invalid CMT data format", data)
	}

	messagePDUHex := matches[1]
	handleMessagePDU(messagePDUHex)
}
