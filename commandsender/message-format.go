package commandsender

import (
	"fmt"
)

const (
	MESSAGE_FORMAT_PDU  = 0
	MESSAGE_FORMAT_TEXT = 1
)

// set the message format that is returned through some unsolicited result code, such as when the model received a SMS message, or when you are reading a SMS message from memory
func SetMessageFormat(format int) {
	Raw(fmt.Sprintf("AT+CMGF=%d\r", format), 1)
}
