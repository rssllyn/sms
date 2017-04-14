package card

type info struct {
	Address                string // sim卡号码
	AddressType            int    // 号码类型
	MaxMessageStorage      int    // 最多可以存储多少条短信
	CurrentlyStoredMessage int    // 当前存储的短信数量
}

var Info info
