package pdu

type VPF uint

const (
	VPNotPresent = iota
	VPEnhanced
	VPRelative
	VPAbsolute
)

// 发送的短信
type SubmitMessage struct {
	SMSC_Number            string // 服务中心号码
	ValidityPeriodFormat   VPF
	ValidityPeriod         string
	MessageReference       byte
	DestinationAddress     string
	DestinationAddressType int
	TP_PID_HEX             string
	TP_DCS                 byte
	TP_DCS_DESC            string
	Message                string // 短信内容
}

func (m *SubmitMessage) SetAddressType(addrType byte) {
	m.DestinationAddressType = int(addrType)
}

func (m *SubmitMessage) SetAddress(addr string) {
	m.DestinationAddress = addr
}

func (m *SubmitMessage) GetTPDCS() byte {
	return m.TP_DCS
}

func (m *SubmitMessage) SetMessage(msg string) {
	m.Message = msg
}
