package resulthandler

import (
	"log"
)

// 从SIM卡中读取到的单条短信
func cmgr(data string) {
	matches := RESULT_DATA_CMGR.FindStringSubmatch(data)
	if len(matches) != 2 {
		log.Println("CMGR data format error")
		return
	}
	messagePDUHex := matches[1]
	handleMessagePDU(messagePDUHex)
}
