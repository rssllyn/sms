package signin

import (
	"bytes"
	"encoding/json"
	"github.com/rssllyn/sms/conf"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	SIGNIN_RETRY_INTERVAL_SECONDS = 60
)

func init() {
	go signin()
}

// 程序启动后初次登录，或者由于断网很长时间导致refresh_token也失效了，也可以调用此方法重新登录
func signin() {
	if conf.Conf["sync_with_server"].(bool) == false {
		return
	}
	for {
		cred := credential{
			UserName: conf.Conf["user_name"].(string),
			Password: conf.Conf["password"].(string),
		}
		body, nil := json.Marshal(&cred)
		request, err := http.NewRequest("POST", conf.Conf["url_base"].(string)+conf.Conf["url_signin"].(string), bytes.NewReader(body))
		if err != nil {
			log.Fatal("failed created signin request")
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			// 网络未连接;
			log.Println("failed sending request to signin", err)
			time.Sleep(SIGNIN_RETRY_INTERVAL_SECONDS * time.Second)
			continue
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			// 可能是服务器维护灯原因导致错误
			log.Println("signin failed, server return status %s", response.Status)
			time.Sleep(SIGNIN_RETRY_INTERVAL_SECONDS * time.Second)
			continue
		}
		respData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("error reading signin response data", err)
			time.Sleep(SIGNIN_RETRY_INTERVAL_SECONDS * time.Second)
			continue
		}
		var t token
		err = json.Unmarshal(respData, &t)
		if err != nil {
			log.Fatal("signin response format can not be recognized")
		}
		log.Println("signin response", string(respData))
		// 保存access token
		tokenCacheMutex.Lock()
		tokenCache = &t
		tokenCache.ExpiredAt = time.Now().Add(time.Duration(tokenCache.ExpiresIn) * time.Second)
		log.Printf("signin succeeded, %#v", tokenCache)
		tokenCacheMutex.Unlock()
		break
	}
}
