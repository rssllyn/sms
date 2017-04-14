package serversync

import (
	"github.com/rssllyn/sms/conf"
	"log"
	"time"
)

// 定期将未发送成功的短信提交到服务器
func SyncMessageWithServer() {
	if conf.Conf["sync_with_server"].(bool) == false {
		return
	}

	tick := time.Tick(time.Duration(conf.Conf["sync_interval_seconds"].(float64)) * time.Second)
	for {
		log.Println("syncing")
		<-tick
		SubmitRetry()
	}
}
