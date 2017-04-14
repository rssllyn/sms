package commandsender

import (
	"github.com/rssllyn/sms/card"
	"log"
	"regexp"
	"strconv"
)

var (
	STORAGE_INFO_PATTERN = regexp.MustCompile(`\+CPMS: "\w+",\d+,\d+,"\w+",\d+,\d+,"\w+",(\d+),(\d+)`)
)

// Set the location where various SMS messages operations would be read or deleted from, or be saved to
func SetPreferredSMSMessageStorage(readAndDeleteFrom, writeAndSendFrom, receivedSMSTo string) {
	Raw(`AT+CPMS="SM","SM","SM"`+"\r", 2)
}

// Get home many SMS messages the set memory location can store at maximum, and how many has been stored
func GetMessageStorageInfo() {
	result := Raw(`AT+CPMS?`+"\r", 2)
	if result.Success && len(result.Frames) == 2 {
		storageInfo := result.Frames[0]
		matches := STORAGE_INFO_PATTERN.FindStringSubmatch(storageInfo)
		if len(matches) != 3 {
			log.Println("invalid cpms result", storageInfo)
			return
		}
		card.Info.CurrentlyStoredMessage, _ = strconv.Atoi(matches[1])
		card.Info.MaxMessageStorage, _ = strconv.Atoi(matches[2])
		log.Printf("we could store at most %d messages, and %d stored now, after that we will deleting all stored messages", card.Info.MaxMessageStorage, card.Info.CurrentlyStoredMessage)
	}

}
