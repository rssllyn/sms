package commandsender

import ()

// read all stored SMS message from memory
func ReadSMSMessageAll() {
	Raw(`AT+CMGL="ALL"`+"\r", 1)
}
