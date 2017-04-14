package resulthandler

import (
	"github.com/rssllyn/sms/card"
	"github.com/rssllyn/sms/commandsender"
	"github.com/rssllyn/sms/conf"
	"log"
	"regexp"
	"strconv"
)

var (
	RECEIVED_MEESAGE_STORAGE_INDEX_PATTERN = regexp.MustCompile(`"\w+",(\d+)`)
)

// 收到短信，先存储，后通知, +CNMI=0,1,0,0,0
// 	+CMTI: "SM",9
func cmti(data string) {
	matches := RESULT_DATA_CMTI.FindStringSubmatch(data)
	log.Println(matches)
	if matches == nil {
		log.Fatalln("CMTI data format error")
	}
	indexStr := matches[2]
	log.Printf("text message stored in %s", indexStr)
	index, _ := strconv.ParseInt(indexStr, 10, 8)
	commandsender.ReadSMSMessage(int(index))
	card.Info.CurrentlyStoredMessage++
	if card.Info.MaxMessageStorage-card.Info.CurrentlyStoredMessage < int(conf.Conf["message_storage_threshold"].(float64)) {
		commandsender.MessageDeleteAll()
	}
}
