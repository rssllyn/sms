package commandsender

import (
	"fmt"
)

// mt: if set to 1, received SMS message would be indicated by a message in the format: +CMTI: "SM",4
func SetNewSMSMessageIndication(mode, mt, bm, ds, brf int) {
	Raw(fmt.Sprintf("AT+CNMI=%d,%d,%d,%d,%d\r", mode, mt, bm, ds, brf), 1)
}
