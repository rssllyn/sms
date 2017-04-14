package signin

import (
	"errors"
	"time"
)

var (
	REFRESH_TOKEN_EXPIRED = errors.New("refresh token has expired")
)

type token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    uint   `json:"expires_in"`
	ExpiredAt    time.Time
}

// 调用token刷新接口刷新access token，如果refresh token也失效了，则重新调用登录接口获取新的access token和refresh token
func (t *token) Refresh() error {
	return REFRESH_TOKEN_EXPIRED
}
