package pdu

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	pduHexFormat = regexp.MustCompile("^[0-9A-F]+$") // 要求是正确的16进制表示
)

func FromHexString(pduString string) (interface{}, error) {
	if !pduHexFormat.MatchString(pduString) {
		return nil, errors.New("pdu data is not a valid hex string")
	}

	read := 0

	SMSC_Length_64, _ := strconv.ParseInt(pduString[0:2], 16, 16)
	SMSC_Length := int(SMSC_Length_64)

	read += 2

	var smscNumber string
	if SMSC_Length != 0 {
		// 获取服务中心号码

		SMSC := pduString[read : read+SMSC_Length*2]
		read += SMSC_Length * 2

		SMSC_TypeOfAddress_64, _ := strconv.ParseInt(SMSC[0:2], 16, 16)
		SMSC_TypeOfAddress := int(SMSC_TypeOfAddress_64)

		smscNumber = SMSC[2:]
		smscNumber = swapSemiBytes(smscNumber)
		// 移除最后的F
		if strings.ToLower(smscNumber[len(smscNumber)-1:]) == "f" {
			smscNumber = smscNumber[:len(smscNumber)-1]
		}

		if SMSC_TypeOfAddress == 0x91 {
			smscNumber = "+" + smscNumber
		}
	}

	SMSDeliver_FirstOcet_64, _ := strconv.ParseInt(pduString[read:read+2], 16, 16)
	SMSDeliver_FirstOcet := int(SMSDeliver_FirstOcet_64)
	read += 2

	if SMSDeliver_FirstOcet&0x03 == 1 {
		// 发送短信
		pdu := SubmitMessage{SMSC_Number: smscNumber}

		// 取第4、3两位
		pdu.ValidityPeriodFormat = VPF((SMSDeliver_FirstOcet >> 3) & 0x03)

		// messageReference := pduString[read:read+2]
		read += 2

		read += ParseAddress(pduString[read:], &pdu)

		pdu.TP_PID_HEX = pduString[read : read+2]
		read += 2

		tpDCS, _ := strconv.ParseInt(pduString[read:read+2], 16, 16)
		read += 2

		pdu.TP_DCS = byte(tpDCS)
		pdu.TP_DCS_DESC = tpDCSMeaning(pdu.TP_DCS)
		log.Printf("%s %02X", pdu.TP_PID_HEX, pdu.TP_DCS)

		// 有效期
		switch pdu.ValidityPeriodFormat {
		case VPNotPresent:
			read += 0
		case VPEnhanced:
			read += 14
		case VPRelative:
			read += 2
		case VPAbsolute:
			read += 14
		}

		log.Println("validity format", pdu.ValidityPeriodFormat)
		log.Println("parsing message", pduString[read:])
		r, err := ParseMessage(pduString[read:], &pdu)
		if err != nil {
			return nil, err
		}
		read += r

		return &pdu, nil
	} else if SMSDeliver_FirstOcet&0x03 == 0 {
		// 收到短信
		pdu := DeliveryMessage{SMSC_Number: smscNumber}

		read += ParseAddress(pduString[read:], &pdu)

		pdu.TP_PID_HEX = pduString[read : read+2]
		read += 2

		tpDCS, _ := strconv.ParseInt(pduString[read:read+2], 16, 16)
		read += 2

		pdu.TP_DCS = byte(tpDCS)
		pdu.TP_DCS_DESC = tpDCSMeaning(pdu.TP_DCS)

		timestampStr := swapSemiBytes(pduString[read : read+14])
		read += 14

		// 将时间转为UTC时间，timestampBytes[6]表示时区(以一刻钟为单位，比如北京时间的时区用32表示)
		year, _ := strconv.Atoi(timestampStr[0:2]) // 年份后两位
		month, _ := strconv.Atoi(timestampStr[2:4])
		day, _ := strconv.Atoi(timestampStr[4:6])
		hour, _ := strconv.Atoi(timestampStr[6:8])
		minute, _ := strconv.Atoi(timestampStr[8:10])
		second, _ := strconv.Atoi(timestampStr[10:12])
		timeZoneInQuartersOfHour, _ := strconv.Atoi(timestampStr[12:14]) // 以1/4小时为单位表示的时区

		currentYear := time.Now().Year()
		century := currentYear - currentYear%100

		pdu.Time = time.Date(century+int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC).Add(-time.Duration(timeZoneInQuartersOfHour/4) * time.Hour)

		r, err := ParseMessage(pduString[read:], &pdu)
		if err != nil {
			return nil, err
		}
		read += r

		return &pdu, nil
	} else {
		return nil, nil
	}
}

func tpDCSMeaning(tpDCS byte) string {
	tpDCSDesc := fmt.Sprintf("%02X", tpDCS)
	switch tpDCS & 0xC0 {
	case 0x00:
	case 0x40:
	case 0x80:
	case 0xC0:
	}
	return tpDCSDesc
}

// 将16进制表示的字符串中，每个字节的高4位低4位对换
func swapSemiBytes(hexString string) string {
	out := ""
	for i := 0; i < len(hexString); i += 2 {
		out += hexString[i+1:i+2] + hexString[i:i+1]
	}
	return out
}
