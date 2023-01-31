package access_token

import (
	"encoding/json"
	"fmt"
	"github.com/HeartGarlic/kuaishou-server-api-sdk/cache"
	"github.com/HeartGarlic/kuaishou-server-api-sdk/util"
	"net/url"
	"sync"
	"time"
)

const accessTokenURL = "https://open.kuaishou.com/oauth2/access_token"

// AccessToken 管理AccessToken 的基础接口
type AccessToken interface {
	GetCacheKey() string             // 获取缓存的key
	SetCacheKey(key string)          // 设置缓存key
	GetAccessToken() (string, error) // 获取token
}

// DefaultAccessToken 默认的token管理类
type DefaultAccessToken struct {
	AppId               string      // app_id	string	是	小程序的 app_id
	AppSecret           string      // app_secret	string	是	小程序的密钥
	GrantType           string      // grant_type	string	是	固定值“client_credentials”
	Cache               cache.Cache // 缓存组件
	accessTokenLock     *sync.Mutex // 读写锁
	accessTokenCacheKey string      // 缓存的key
}

// NewDefaultAccessToken 实例化默认的token管理类
func NewDefaultAccessToken(appId, appSecret string, cache cache.Cache) AccessToken {
	if cache == nil {
		panic(any("cache is need"))
	}
	token := &DefaultAccessToken{
		AppId:               appId,
		AppSecret:           appSecret,
		GrantType:           "client_credentials",
		Cache:               cache,
		accessTokenCacheKey: fmt.Sprintf("kuaishou_server_api_sdk_access_token_%s", appId),
		accessTokenLock:     new(sync.Mutex),
	}
	return token
}

// GetCacheKey 获取缓存key
func (dd *DefaultAccessToken) GetCacheKey() string {
	return dd.accessTokenCacheKey
}

// SetCacheKey 设置缓存key
func (dd *DefaultAccessToken) SetCacheKey(key string) {
	dd.accessTokenCacheKey = key
}

// GetAccessToken 获取token
func (dd *DefaultAccessToken) GetAccessToken() (string, error) {
	// 先尝试从缓存中获取如果不存在就调用接口获取
	if val := dd.Cache.Get(dd.GetCacheKey()); val != nil {
		return val.(string), nil
	}

	// 加锁防止并发获取接口
	dd.accessTokenLock.Lock()
	defer dd.accessTokenLock.Unlock()

	// 双捡防止重复获取
	if val := dd.Cache.Get(dd.GetCacheKey()); val != nil {
		return val.(string), nil
	}

	// 开始调用接口获取token
	reqAccessToken, err := GetTokenFromServer(accessTokenURL, dd.AppId, dd.AppSecret)
	if err != nil {
		return "", err
	}
	// 设置缓存
	expires := reqAccessToken.ExpiresIn - 1500
	err = dd.Cache.Set(dd.GetCacheKey(), reqAccessToken.AccessToken, time.Duration(expires)*time.Second)
	if err != nil {
		return "", err
	}
	return reqAccessToken.AccessToken, nil
}

// ResAccessToken 获取token的返回结构体
type ResAccessToken struct {
	Result      int    `json:"result,omitempty"`
	ErrorMsg    string `json:"error_msg,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

// GetTokenFromServer 从快手服务器获取token
func GetTokenFromServer(apiUrl string, appId, appSecret string) (resAccessToken ResAccessToken, err error) {
	params := url.Values{"app_id": []string{appId}, "app_secret": []string{appSecret}, "grant_type": []string{"client_credentials"}}
	body, err := util.PostForm(apiUrl, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.Result != 1 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.Result, resAccessToken.ErrorMsg)
		return
	}
	return
}
