package pdu

import (
	"strconv"
	"strings"
)

type Address interface {
	SetAddressType(byte)
	SetAddress(string)
}

func ParseAddress(hexStartWithAddress string, address Address) (addressDataLength int) {
	read := 0

	addrLen, _ := strconv.ParseInt(hexStartWithAddress[read:read+2], 16, 16)
	read += 2

	// 当号码为奇数(比如11位手机号)是，数据中会有一个填充的F字符
	if addrLen%2 != 0 {
		addrLen += 1
	}

	addrType, _ := strconv.ParseInt(hexStartWithAddress[read:read+2], 16, 16)
	read += 2

	addr := swapSemiBytes(hexStartWithAddress[read : read+int(addrLen)])
	read += int(addrLen)

	// 移除最后的F
	if strings.ToUpper(addr[addrLen-1:]) == "F" {
		addr = addr[:addrLen-1]
	}

	if addrType == 0x91 {
		addr = "+" + addr
	}

	address.SetAddressType(byte(addrType))
	address.SetAddress(addr)

	return read
}
