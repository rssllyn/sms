package commandsender

// delete all SMS messages stored in memory (whatever set via AT+CPMS)
func MessageDeleteAll() {
	Raw("AT+CMGD=1,4\r", 1)
}
