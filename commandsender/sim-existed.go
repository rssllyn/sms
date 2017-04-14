package commandsender

func SimExisted() {
	Raw("AT+CCID\r", 2)
}
