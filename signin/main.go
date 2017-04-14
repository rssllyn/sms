package signin

import (
	"sync"
	"time"
)

var (
	tokenCache      *token
	tokenCacheMutex sync.Mutex
)

// 判断当前是否已经登录，并且access token还有效，立即返回结果，返回的字符串非空表示当前处于登录状态。如果access token已经失效，通过重新登录或者刷新token的方式保证token有效
func GetAccessToken() string {
	tokenCacheMutex.Lock()
	defer tokenCacheMutex.Unlock()

	if tokenCache == nil {
		// 未登录，有其他逻辑正在尝试登录，这里不做任何处理直接返回
		return ""
	}
	if time.Now().Before(tokenCache.ExpiredAt) {
		// access token还在有效期
		return tokenCache.AccessToken
	}
	err := tokenCache.Refresh()
	if err == nil {
		// 刷新token成功
		return tokenCache.AccessToken
	}
	if err == REFRESH_TOKEN_EXPIRED {
		// refresh token失效，重新登录
		tokenCache = nil
		go signin()
		return ""
	} else {
		// 其他错误，可能是网路问题或者服务器问题，直接返回，之后重新尝试
		return ""
	}
}
