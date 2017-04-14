package pdu

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Message interface {
	SetMessage(msg string)
	GetTPDCS() byte
}

func ParseMessage(hexStartWithMessage string, msg Message) (int, error) {
	read := 0

	// 以字节为单位的消息长度
	messageLength, _ := strconv.ParseInt(hexStartWithMessage[read:read+2], 16, 16)
	read += 2

	if read+2*int(messageLength) > len(hexStartWithMessage) {
		return 0, errors.New(fmt.Sprintf("the message denotes %d bytes of its length, but is contains only %d hex characters", messageLength, len(hexStartWithMessage)-2))
	}

	messageBytes, _ := hex.DecodeString(hexStartWithMessage[read : read+int(2*messageLength)])
	bitSize := bitSizeFromTPDCS(msg.GetTPDCS())
	switch bitSize {
	case 7:
		msg.SetMessage(getUserMessage(messageBytes))
	case 8:
		msg.SetMessage(getUserMessage8(messageBytes))
	case 16:
		msg.SetMessage(getUserMessage16(messageBytes))
	}
	return read, nil
}

func bitSizeFromTPDCS(tpDCS byte) (bitSize int) {
	bitSize = 7
	switch tpDCS & 0xC0 {
	case 0x00:
		switch tpDCS & 0x0C {
		case 0x04:
			bitSize = 8
		case 0x08:
			bitSize = 16
		}
	case 0xC0:
		switch tpDCS & 0x30 {
		case 0x20:
			bitSize = 16
		case 0x30:
			if tpDCS&0x04 != 0 {
				bitSize = 8
			}
		}
	}
	return
}

func getUserMessage(messageBytes []byte) string {
	log.Println("7-bit message")
	var messageRunes []rune

	var restBinary []string
	var septetsBinary []string

	for i, b := range messageBytes {
		// 将每个字节转换成2进制表示，并拆分成两部分
		byteBinary := fmt.Sprintf("%08b", b)
		restBinary = append(restBinary, byteBinary[0:i%7+1])      // 包含byteBinary里面的前1到7位，循环
		septetsBinary = append(septetsBinary, byteBinary[i%7+1:]) // 包含byteBinary里面的后7到1位，循环
	}
	log.Println(restBinary)
	log.Println(septetsBinary)
	for i, b := range septetsBinary {
		if i%7 == 0 {
			if i != 0 {
				messageRunes = append(messageRunes, binaryToRune(restBinary[i-1]))
			}
			messageRunes = append(messageRunes, binaryToRune(b))
		} else {
			messageRunes = append(messageRunes, binaryToRune(b+restBinary[i-1]))
		}
	}
	if len(restBinary[len(restBinary)-1]) == 7 {
		messageRunes = append(messageRunes, binaryToRune(restBinary[len(restBinary)-1]))
	}
	return string(messageRunes)
}

func getUserMessage8(messageBytes []byte) string {
	log.Println("8-bit message")
	var runes []rune
	for _, b := range messageBytes {
		runes = append(runes, rune(b))
	}
	return string(runes)
}

// messageHex中每2个子额接表示一个unicode字符的unicode code point
func getUserMessage16(messageBytes []byte) string {
	var runes []rune
	for i := 0; i < len(messageBytes); i += 2 {
		var unicodeCodePoint uint16
		binary.Read(bytes.NewReader(messageBytes[i:i+2]), binary.BigEndian, &unicodeCodePoint)
		runes = append(runes, rune(unicodeCodePoint))
	}
	return string(runes)
}

func binaryToRune(bin string) rune {
	value, _ := strconv.ParseInt(bin, 2, 8)
	return rune(value)
}
