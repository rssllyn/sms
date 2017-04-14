package commandsender

// an empty command, send this as the first command after modem initizlied would make the moderm more stable somehow
func Empty() {
	Raw("\r", 0)
}
