package kuaishou_server_api_sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	accessToken "github.com/HeartGarlic/kuaishou-server-api-sdk/access-token"
	"github.com/HeartGarlic/kuaishou-server-api-sdk/cache"
	"github.com/HeartGarlic/kuaishou-server-api-sdk/util"
	"net/url"
	"sort"
	"strings"
)

// 声明常量
const (
	successCode               = 1
	code2Session              = "https://open.kuaishou.com/oauth2/mp/code2session"
	payCreateOrder            = "https://open.kuaishou.com/openapi/mp/developer/epay/create_order"              // 有收银台版本
	payCreateOrderWithChannel = "https://open.kuaishou.com/openapi/mp/developer/epay/create_order_with_channel" // 无收银台版本
	queryOrder                = "https://open.kuaishou.com/openapi/mp/developer/epay/query_order"               // 查询支付状态
)

// 快手小程序的服务端golang sdk
// 包含登陆 获取access token
// 担保支付

// KuaiShou 基础的客户端类
type KuaiShou struct {
	BaseApiHost string      // api基础地址 https://open.kuaishou.com/
	AppId       string      // 快手小程序的appid
	AppSecret   string      // 快手小程序的app secret
	Cache       cache.Cache // 基础的缓存接口
	AccessToken accessToken.AccessToken
}

// KuaiShouAppletConfig 快手小程序需要的参数
type KuaiShouAppletConfig struct {
	AppId       string      // 快手小程序的appid
	AppSecret   string      // 快手小程序的app secret
	Cache       cache.Cache // 基础的缓存接口
	AccessToken accessToken.AccessToken
}

// NewKuaiShou 实例化一个快手客户端
func NewKuaiShou(config *KuaiShouAppletConfig) *KuaiShou {
	// 如果存在cache组件就实例化一个内存缓存
	if config.Cache == nil {
		config.Cache = cache.NewMemory()
	}
	// 如果未设置token管理 就使用默认的
	if config.AccessToken == nil {
		config.AccessToken = accessToken.NewDefaultAccessToken(config.AppId, config.AppSecret, config.Cache)
	}
	return &KuaiShou{
		AppId:       config.AppId,
		AppSecret:   config.AppSecret,
		Cache:       config.Cache,
		AccessToken: config.AccessToken,
	}
}

// Code2SessionResponse ...
type Code2SessionResponse struct {
	Result     int    `json:"result,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	OpenId     string `json:"open_id,omitempty"`
	ErrorMsg   string `json:"error_msg"`
}

// Code2Session 实现具体的业务方法 登陆
func (k *KuaiShou) Code2Session(code string) (Code2SessionResponse, error) {
	post, err := util.PostForm(code2Session, url.Values{"js_code": []string{code}, "app_id": []string{k.AppId}, "app_secret": []string{k.AppSecret}})
	if err != nil {
		return Code2SessionResponse{}, err
	}
	var response Code2SessionResponse
	err = json.Unmarshal(post, &response)
	if err != nil {
		return Code2SessionResponse{}, err
	}
	if response.Result != successCode {
		return Code2SessionResponse{}, fmt.Errorf(response.ErrorMsg)
	}
	return response, nil
}

// PayCreateOrderParams 预下单所需参数
type PayCreateOrderParams struct {
	OutOrderNo           string               `json:"out_order_no,omitempty"`            // out_order_no	string[6,32]	是	是	body json	商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一 示例值：1217752501201407033233368018
	OpenId               string               `json:"open_id,omitempty"`                 // open_id	string	是	是	body json	快手用户在当前小程序的open_id，可通过login操作获取。
	TotalAmount          int64                `json:"total_amount,omitempty"`            // total_amount	number	是	是	body json	用户支付金额，单位为[分]。不允许传非整数的数值。
	Subject              string               `json:"subject,omitempty"`                 // subject	string[1,128]	是	是	body json	商品描述。注：1汉字=2字符。
	Detail               string               `json:"detail,omitempty"`                  // detail	string[1,1024]	是	是	body json	商品详情。注：1汉字=2字符。
	Type                 int64                `json:"type,omitempty"`                    // type	number	是	是	body json	商品类型，不同商品类目的编号见 担保支付商品类目编号
	ExpireTime           int64                `json:"expire_time,omitempty"`             // expire_time	number	是	是	body json	订单过期时间，单位秒，300s - 172800s
	Sign                 string               `json:"sign,omitempty"`                    // sign	string	是	否	body json	开发者对核心字段签名, 签名方式见 附录
	Attach               string               `json:"attach,omitempty"`                  // attach	string[0,128]	否	是	body json	开发者自定义字段，回调原样回传.注：1汉字=2字符；勿回传敏感信息
	NotifyUrl            string               `json:"notify_url,omitempty"`              // notify_url	string[1, 256]	是	是	body json	通知URL必须为直接可访问的URL，不允许携带查询串。
	GoodsId              string               `json:"goods_id,omitempty"`                // goods_id	string[1, 256]	否(本地生活类必填)	是	body json	下单商品id，需与商品对接 (opens new window)时的product_id一致，长度限制256个英文字符，1个汉字=2个英文字符；
	GoodsDetailUrl       string               `json:"goods_detail_url,omitempty"`        // goods_detail_url	string[1, 500]	否(本地生活类必填)	是	body json	订单详情页跳转path。长度限制500个英文字符，1个汉字=2个英文字符； 示例值：/page/index/anima
	MultiCopiesGoodsInfo MultiCopiesGoodsInfo `json:"multi_copies_goods_info,omitempty"` // multi_copies_goods_info	string[1, 500]	否(单商品多份场景必填)	是	body json	单商品购买多份场景，示例值：[{"copies":2}]， 内容见multi_copies_goods_info字段说明 multi_copies_goods_info 字段说明 字段名    类型    说明 copies    number    购买份数
	CancelOrder          int64                `json:"cancel_order,omitempty"`            // cancel_order	number	否	是	body json	该字段表示创建订单的同时是否覆盖之前已存在的订单。 取值范围: [0, 1]。 0:不覆盖 1:覆盖 使用说明：如果传值为1 重复调用接口后执行逻辑为先删除已存在的订单再创建新订单，如果传值为0 重复调用接口执行逻辑为直接返回已创建订单的订单信息。如果不传该参数则和传值为0逻辑一致
	Provider             Provider             `json:"provider,omitempty"`
}

type MultiCopiesGoodsInfo struct {
	Copies int64 `json:"copies,omitempty"`
}

type Provider struct {
	Provider            string `json:"provider,omitempty"`              // 支付方式，枚举值，目前只支持"WECHAT"、"ALIPAY"两种
	ProviderChannelType string `json:"provider_channel_type,omitempty"` // 支付方式子类型，枚举值，目前只支持"NORMAL"
}

// PayCreateOrderResponse 预下单返回结果值
type PayCreateOrderResponse struct {
	Result    int       `json:"result,omitempty"`
	ErrorMsg  string    `json:"error_msg,omitempty"`
	OrderInfo OrderInfo `json:"order_info,omitempty"`
}

// OrderInfo 预下单的订单信息
type OrderInfo struct {
	OrderNo        string `json:"order_no,omitempty"`
	OrderInfoToken string `json:"order_info_token,omitempty"`
}

// PayCreateOrder 预下单
func (k *KuaiShou) PayCreateOrder(payCreateOrderParams PayCreateOrderParams) (PayCreateOrderResponse, error) {
	token, err := k.AccessToken.GetAccessToken()
	if err != nil {
		return PayCreateOrderResponse{}, err
	}
	host := payCreateOrder
	if len(payCreateOrderParams.Provider.Provider) > 0 {
		host = payCreateOrderWithChannel
	}
	// 拼接请求地址
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", host, k.AppId, token)
	// 拼接请求参数 还需要加签
	paramsMap, err := JsonStructToMap(payCreateOrderParams)
	paramsMap["app_id"] = k.AppId
	paramsMap["multi_copies_goods_info"] = ""
	paramsMap["provider"] = ""
	// 重新修改生成map的方法 有些字段需要是 json string multi_copies_goods_info
	if payCreateOrderParams.MultiCopiesGoodsInfo != (MultiCopiesGoodsInfo{}) {
		multiCopiesGoodsInfo, _ := json.Marshal(payCreateOrderParams.MultiCopiesGoodsInfo)
		paramsMap["multi_copies_goods_info"] = multiCopiesGoodsInfo
	}
	if payCreateOrderParams.Provider != (Provider{}) {
		provider, _ := json.Marshal(payCreateOrderParams.Provider)
		paramsMap["provider"] = provider
	}
	paramsMap["sign"] = k.GenerateSign(paramsMap)
	// 开始请求接口
	postJSON, err := util.PostJSON(api, paramsMap)
	if err != nil {
		return PayCreateOrderResponse{}, err
	}
	// 解析返回值
	var payResponse PayCreateOrderResponse
	err = json.Unmarshal(postJSON, &payResponse)
	if err != nil {
		return PayCreateOrderResponse{}, err
	}
	if payResponse.Result != successCode {
		return payResponse, fmt.Errorf(payResponse.ErrorMsg)
	}
	return payResponse, nil
}

// PayCallbackCheckSignature 验证回调签名
func (k *KuaiShou) PayCallbackCheckSignature(oldSign, body string) error {
	newSignStr := fmt.Sprintf("%s%s", body, k.AppSecret)
	newSign := fmt.Sprintf("%x", md5.Sum([]byte(newSignStr)))
	if newSign != oldSign {
		return fmt.Errorf("验证签名失败 newSign: %s oldSign: %s", newSign, oldSign)
	}
	return nil
}

// QueryOrderResponse 查询订单的返回值
type QueryOrderResponse struct {
	Result      int         `json:"result,omitempty"`
	ErrorMsg    string      `json:"error_msg,omitempty"`
	PaymentInfo PaymentInfo `json:"payment_info,omitempty"`
}

type PaymentInfo struct {
	TotalAmount     int       `json:"total_amount,omitempty"`
	PayStatus       string    `json:"pay_status,omitempty"`
	PayTime         string    `json:"pay_time,omitempty"`
	PayChannel      string    `json:"pay_channel,omitempty"`
	OutOrderNo      string    `json:"out_order_no,omitempty"`
	KsOrderNo       string    `json:"ks_order_no,omitempty"`
	ExtraInfo       ExtraInfo `json:"extra_info,omitempty"`
	PromotionAmount int       `json:"promotion_amount,omitempty"`
	OpenId          string    `json:"open_id,omitempty"`
}

type ExtraInfo struct {
	Url      string `json:"url,omitempty"`
	ItemType string `json:"item_type,omitempty"`
	ItemId   string `json:"item_id,omitempty"`
	AuthorId string `json:"author_id,omitempty"`
}

// QueryOrder 查询订单状态
// outOrderNo 商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一 1217752501201407033233368018
func (k *KuaiShou) QueryOrder(outOrderNo string) (QueryOrderResponse, error) {
	token, err := k.AccessToken.GetAccessToken()
	if err != nil {
		return QueryOrderResponse{}, err
	}
	params := map[string]interface{}{
		"out_order_no": outOrderNo,
		"app_id":       k.AppId,
	}
	params["sign"] = k.GenerateSign(params)
	delete(params, "app_id")
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", queryOrder, k.AppId, token)
	postJSON, err := util.PostJSON(api, params)
	if err != nil {
		return QueryOrderResponse{}, err
	}
	var queryResponse QueryOrderResponse
	err = json.Unmarshal(postJSON, &queryResponse)
	if err != nil {
		return QueryOrderResponse{}, err
	}
	if queryResponse.Result != successCode {
		return queryResponse, fmt.Errorf(queryResponse.ErrorMsg)
	}
	return queryResponse, nil
}

// GenerateSign 生成请求签名
func (k *KuaiShou) GenerateSign(params map[string]interface{}) string {
	var paramsKey []string
	for k, v := range params {
		if k == "sign" || k == "access_token" || k == "" {
			continue
		}
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		if value == "" || value == "null" {
			continue
		}
		paramsKey = append(paramsKey, k)
	}
	sort.Strings(paramsKey)
	// 根据排序后的key拼接字符串
	var paramsVal []string
	for _, v := range paramsKey {
		paramsVal = append(paramsVal, fmt.Sprintf("%s=%+v", v, params[v]))
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(paramsVal, "&")+k.AppSecret)))
}

// JsonStructToMap ...
func JsonStructToMap(content interface{}) (map[string]interface{}, error) {
	var name map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		return name, err
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&name); err != nil {
		} else {
			for k, v := range name {
				name[k] = v
			}
		}
	}
	return name, nil
}
