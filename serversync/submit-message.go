package serversync

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rssllyn/sms/conf"
	"github.com/rssllyn/sms/signin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ParamMessageCreate struct {
	ReceiverAddress     string    `json:"receiver_address"`      // 接收方号码
	ReceiverAddressType int       `json:"receiver_address_type"` // 接收放号码类型
	SenderAddress       string    `json:"sender_address"`        // 发送方号码
	SenderAddressType   int       `json:"sender_address_type"`   // 发送方号码类型
	ReceivedTime        time.Time `json:"received_time"`
	Message             string    `json:"message"`
}

// 将短信内容发送到服务器
func SubmitMessage(message string, senderAddress, receiverAddress string, senderAddressType, receiverAddressType int, receivedTime time.Time) error {
	if conf.Conf["sync_with_server"].(bool) == false {
		return nil
	}

	accessToken := signin.GetAccessToken()
	if len(accessToken) == 0 {
		return errors.New("not signed in")
	}
	param := ParamMessageCreate{
		ReceiverAddress:     receiverAddress,
		ReceiverAddressType: receiverAddressType,
		SenderAddress:       senderAddress,
		SenderAddressType:   senderAddressType,
		ReceivedTime:        receivedTime,
		Message:             message,
	}
	body, err := json.Marshal(param)
	if err != nil {
		return err
	}
	log.Printf("submitting message: %v", param)
	request, err := http.NewRequest("POST", conf.Conf["url_base"].(string)+conf.Conf["url_message_create"].(string), bytes.NewReader(body))
	if err != nil {
		log.Fatal("failed creating message submit request")
	}
	request.Header.Add("Authorization", "bearer "+accessToken)
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		// 成功
		content, _ := ioutil.ReadAll(rsp.Body)
		return errors.New(fmt.Sprintf("server returned error %s: %s", rsp.Status, content))
	}
	return nil
}
