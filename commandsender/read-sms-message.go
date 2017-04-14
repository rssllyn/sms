package commandsender

import (
	"fmt"
)

// read a stored SMS message from specified storage index
func ReadSMSMessage(index int) {
	command := fmt.Sprintf("AT+CMGR=%d\r", index)
	Raw(command, 2)
}
