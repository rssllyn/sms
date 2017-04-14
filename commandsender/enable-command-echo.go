package commandsender

import (
	"fmt"
)

// set whether the moderm sould echo back the command it received, useful only if you are testing using tools like minicom
func SetEnableCommandEcho(enable bool) {
	value := 0
	if enable {
		value = 1
	}
	command := fmt.Sprintf("ATE%d\r", value)
	Raw(command, 1)
}
