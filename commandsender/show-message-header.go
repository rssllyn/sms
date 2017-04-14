package commandsender

import (
	"fmt"
)

func SetShowMessageHeader(show bool) {
	value := 0
	if show {
		value = 1
	}
	command := fmt.Sprintf("AT+CSDH=%d\r", value)
	Raw(command, 1)
}
