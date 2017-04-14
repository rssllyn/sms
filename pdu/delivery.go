package pdu

import (
	"time"
)

// 接收的短信
type DeliveryMessage struct {
	SMSC_Number        string // 服务中心号码
	SMSC_TypeOfAddress byte
	SenderAddressType  int // 0x91, 0xa1
	SenderAddress      string
	TP_PID_HEX         string
	TP_DCS             byte
	TP_DCS_DESC        string
	Time               time.Time
	Message            string // 短信内容
}

func (m *DeliveryMessage) SetAddressType(addrType byte) {
	m.SenderAddressType = int(addrType)
}

func (m *DeliveryMessage) SetAddress(addr string) {
	m.SenderAddress = addr
}

func (m *DeliveryMessage) GetTPDCS() byte {
	return m.TP_DCS
}

func (m *DeliveryMessage) SetMessage(msg string) {
	m.Message = msg
}
