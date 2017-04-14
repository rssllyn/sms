package commandsender

import (
	"github.com/rssllyn/sms/card"
	"regexp"
	"strconv"
)

var (
	CARD_NUMBER_PATTERN = regexp.MustCompile(`((\r\n)|^)\+CNUM: ".*","(.*)",(\d+)`)
)

// get subscribe number, that is, the number of SIM card that is pluged into the modem
func GetSubscriberNumber() {
	result := Raw("AT+CNUM\r", 1)
	if result.Success {
		matches := CARD_NUMBER_PATTERN.FindStringSubmatch(result.Frames[0])
		if len(matches) != 5 {
			return
		}
		numberType, _ := strconv.Atoi(matches[4])

		card.Info.AddressType = int(numberType)
		card.Info.Address = matches[3]
	}
}
