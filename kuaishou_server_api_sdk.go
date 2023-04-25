package kuaishou_server_api_sdk

import (
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
	applyRefund               = "https://open.kuaishou.com/openapi/mp/developer/epay/apply_refund"              //  支付退款
	queryRefund               = "https://open.kuaishou.com/openapi/mp/developer/epay/query_refund"              // 退款查询接口
	settle                    = "https://open.kuaishou.com/openapi/mp/developer/epay/settle"                    // 结算
	querySettle               = "https://open.kuaishou.com/openapi/mp/developer/epay/query_settle"              // 结算查询
)

// KuaiShou 基础的客户端类
// 快手小程序的服务端golang sdk
// 包含登陆 获取access token
// 担保支付
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
func (k *KuaiShou) Code2Session(code string) (code2SessionResponse Code2SessionResponse, err error) {
	post, err := util.PostForm(code2Session, url.Values{"js_code": []string{code}, "app_id": []string{k.AppId}, "app_secret": []string{k.AppSecret}})
	if err != nil {
		return
	}
	err = json.Unmarshal(post, &code2SessionResponse)
	if err != nil {
		return
	}
	if code2SessionResponse.Result != successCode {
		return code2SessionResponse, fmt.Errorf(code2SessionResponse.ErrorMsg)
	}
	return
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
func (k *KuaiShou) PayCreateOrder(payCreateOrderParams PayCreateOrderParams) (payCreateOrderResponse PayCreateOrderResponse, err error) {
	token, _ := k.AccessToken.GetAccessToken()
	host := payCreateOrder
	if len(payCreateOrderParams.Provider.Provider) > 0 {
		host = payCreateOrderWithChannel
	}
	// 拼接请求地址
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", host, k.AppId, token)
	// 拼接请求参数 还需要加签
	paramsMap, err := util.JsonStructToMap(payCreateOrderParams)
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
		return
	}
	// 解析返回值
	err = json.Unmarshal(postJSON, &payCreateOrderResponse)
	if err != nil {
		return
	}
	if payCreateOrderResponse.Result != successCode {
		return payCreateOrderResponse, fmt.Errorf(payCreateOrderResponse.ErrorMsg)
	}
	return
}

// CallbackCheckSignature 验证回调签名
func (k *KuaiShou) CallbackCheckSignature(oldSign, body string) error {
	newSignStr := fmt.Sprintf("%s%s", body, k.AppSecret)
	newSign := fmt.Sprintf("%x", md5.Sum([]byte(newSignStr)))
	if newSign != oldSign {
		return fmt.Errorf("验证签名失败 newSign: %s oldSign: %s", newSign, oldSign)
	}
	return nil
}

// PayCallbackResponse 支付回调的参数解析
type PayCallbackResponse struct {
	Data      PayCallbackResponseData `json:"data,omitempty"`
	BizType   string                  `json:"biz_type,omitempty"`
	MessageId string                  `json:"message_id,omitempty"`
	AppId     string                  `json:"app_id,omitempty"`
	Timestamp int                     `json:"timestamp,omitempty"`
}

type PayCallbackResponseData struct {
	Channel         string `json:"channel,omitempty"`          // channel	string	支付渠道。取值：UNKNOWN - 未知｜WECHAT-微信｜ALIPAY-支付宝
	OutOrderNo      string `json:"out_order_no,omitempty"`     // out_order_no	string	商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一 示例值：1217752501201407033233368018
	Attach          string `json:"attach,omitempty"`           // attach	string	预下单时携带的开发者自定义信息
	Status          string `json:"status,omitempty"`           // status	string	订单支付状态。 取值： PROCESSING-处理中｜SUCCESS-成功｜FAILED-失败
	KsOrderNo       string `json:"ks_order_no,omitempty"`      // ks_order_no	string	快手小程序平台订单号
	OrderAmount     int    `json:"order_amount,omitempty"`     // order_amount	number	订单金额
	TradeNo         string `json:"trade_no,omitempty"`         // trade_no	string	用户侧支付页交易单号，具体获取方法可点击查看(opens new window)
	ExtraInfo       string `json:"extra_info,omitempty"`       // extra_info	string	订单来源信息，同支付查询接口
	EnablePromotion bool   `json:"enable_promotion,omitempty"` // enable_promotion	boolean	是否参与分销，true:分销，false:非分销
	PromotionAmount int    `json:"promotion_amount,omitempty"` // promotion_amount	number	预计分销金额，单位：分
}

// PayCallbackResponse 解析回调的参数到结构体 并返回
func (k *KuaiShou) PayCallbackResponse(oldSign, body string, checkSign bool) (payCallbackResponse PayCallbackResponse, err error) {
	// 判断是否校验签名
	if checkSign {
		err = k.CallbackCheckSignature(oldSign, body)
		if err != nil {
			return
		}
	}
	// 开始解析数据到结构体
	err = json.Unmarshal([]byte(body), &payCallbackResponse)
	if err != nil {
		return
	}
	return
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
	PayTime         int       `json:"pay_time,omitempty"`
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
func (k *KuaiShou) QueryOrder(outOrderNo string) (queryOrderResponse QueryOrderResponse, err error) {
	token, _ := k.AccessToken.GetAccessToken()
	params := map[string]interface{}{
		"out_order_no": outOrderNo,
	}
	params["sign"] = k.GenerateSign(params)
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", queryOrder, k.AppId, token)
	postJSON, err := util.PostJSON(api, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(postJSON, &queryOrderResponse)
	if err != nil {
		return
	}
	if queryOrderResponse.Result != successCode {
		return queryOrderResponse, fmt.Errorf(queryOrderResponse.ErrorMsg)
	}
	return
}

// GenerateSign 生成请求签名
func (k *KuaiShou) GenerateSign(params map[string]interface{}) string {
	params["app_id"] = k.AppId
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

// ApplyRefundParams 支付退款接口参数
type ApplyRefundParams struct {
	OutOrderNo           string               `json:"out_order_no,omitempty"`            // out_order_no	string	是	是	body json	开发者需要发起退款的支付订单号，商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一示例值：1217752501201407033233368018
	OutRefundNo          string               `json:"out_refund_no,omitempty"`           // out_refund_no	string	是	是	body json	开发者的退款单号
	Reason               string               `json:"reason,omitempty"`                  // reason	string[1,128]	是	是	body json	退款理由。1个字符=2个汉字
	Attach               string               `json:"attach,omitempty"`                  // attach	string[0,128]	否	是	body json	开发者自定义字段，回调原样回传. 注：1汉字=2字符；勿回传敏感信息
	NotifyUrl            string               `json:"notify_url,omitempty"`              // notify_url	string[1,256]	是	是	body json	通知URL必须为直接可访问的URL，不允许携带查询串。
	RefundAmount         int64                `json:"refund_amount,omitempty"`           // refund_amount	number	否	是	body json	用户退款金额，单位为分。不允许传非整数的数值
	Sign                 string               `json:"sign,omitempty"`                    // sign	string	是	否	body json	开发者对核心字段签名, 防止传输过程中出现意外，签名方式见附录
	MultiCopiesGoodsInfo MultiCopiesGoodsInfo `json:"multi_copies_goods_info,omitempty"` // multi_copies_goods_info	string[1, 500]	否(单商品多份场景必填)	是	body json	单商品购买多份场景，示例值：[{"copies":2}]， 内容见multi_copies_goods_info字段说明
}

// ApplyRefundResponse 退款接口返回值
type ApplyRefundResponse struct {
	Result   int    `json:"result,omitempty"`
	ErrorMsg string `json:"error_msg,omitempty"`
	RefundNo string `json:"refund_no,omitempty"`
}

// ApplyRefund 支付退款接口
func (k *KuaiShou) ApplyRefund(applyRefundParams ApplyRefundParams) (applyRefundResponse ApplyRefundResponse, err error) {
	token, _ := k.AccessToken.GetAccessToken()
	params, _ := util.JsonStructToMap(applyRefundParams)
	params["multi_copies_goods_info"] = ""
	if applyRefundParams.MultiCopiesGoodsInfo.Copies > 0 {
		params["multi_copies_goods_info"], _ = json.Marshal(applyRefundParams.MultiCopiesGoodsInfo)
	}
	params["sign"] = k.GenerateSign(params)
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", applyRefund, k.AppId, token)
	postJSON, err := util.PostJSON(api, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(postJSON, &applyRefundResponse)
	if err != nil {
		return
	}
	return
}

// QueryRefundResponse 退款查询接口返回值
type QueryRefundResponse struct {
	Result     int        `json:"result,omitempty"`
	ErrorMsg   string     `json:"error_msg,omitempty"`
	RefundInfo RefundInfo `json:"refund_info,omitempty"`
}

type RefundInfo struct {
	KsOrderNo    string `json:"ks_order_no,omitempty"`
	RefundStatus string `json:"refund_status,omitempty"`
	RefundNo     string `json:"refund_no,omitempty"`
	KsRefundType string `json:"ks_refund_type,omitempty"`
	RefundAmount int    `json:"refund_amount,omitempty"`
	KsRefundNo   string `json:"ks_refund_no,omitempty"`
}

// QueryRefund 退款查询接口
func (k *KuaiShou) QueryRefund(outRefundNo string) (queryRefundResponse QueryRefundResponse, err error) {
	token, _ := k.AccessToken.GetAccessToken()
	params := map[string]interface{}{
		"out_refund_no": outRefundNo,
	}
	params["sign"] = k.GenerateSign(params)
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", queryRefund, k.AppId, token)
	postJSON, err := util.PostJSON(api, params)
	if err != nil {
		return QueryRefundResponse{}, err
	}
	err = json.Unmarshal(postJSON, &queryRefundResponse)
	if err != nil {
		return
	}
	return
}

// ApplyRefundCallbackResponse 退款回调结构体
type ApplyRefundCallbackResponse struct {
	Data      ApplyRefundCallbackResponseData `json:"data,omitempty"`
	MessageId string                          `json:"message_id,omitempty"`
	BizType   string                          `json:"biz_type,omitempty"`
	AppId     string                          `json:"app_id,omitempty"`
	Timestamp int                             `json:"timestamp,omitempty"`
}

type ApplyRefundCallbackResponseData struct {
	OutRefundNo  string `json:"out_refund_no,omitempty"`
	RefundAmount int    `json:"refund_amount,omitempty"`
	Attach       string `json:"attach,omitempty"`
	Status       string `json:"status,omitempty"`
	KsOrderNo    string `json:"ks_order_no,omitempty"`
	KsRefundNo   string `json:"ks_refund_no,omitempty"`
	KsRefundType string `json:"ks_refund_type,omitempty"`
}

// ApplyRefundCallback 退款回调值解析
func (k *KuaiShou) ApplyRefundCallback(oldSign, body string, checkSign bool) (applyRefundCallbackResponse ApplyRefundCallbackResponse, err error) {
	// 判断是否需要校验签名
	if checkSign {
		err = k.CallbackCheckSignature(oldSign, body)
		if err != nil {
			return
		}
	}
	// 开始解析数据
	err = json.Unmarshal([]byte(body), &applyRefundCallbackResponse)
	if err != nil {
		return
	}
	return
}

// SettleParams 结算接口参数
type SettleParams struct {
	OutOrderNo           string               `json:"out_order_no,omitempty"`  // out_order_no	string[6,32]	是	是	body json	开发者需要发起结算的支付订单号，商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一 示例值：1217752501201407033233368018
	OutSettleNo          string               `json:"out_settle_no,omitempty"` // out_settle_no	string[6,32]	是	是	body json	开发者的结算单号，小程序唯一。
	Reason               string               `json:"reason,omitempty"`        // reason	string[1,128]	是	是	body json	结算描述，长度限制 128 个字符。1个字符=2个汉字
	Attach               string               `json:"attach,omitempty"`        // attach	string[0, 128]	否	是	body json	开发者自定义字段，回调原样回传. 注：1汉字=2字符；勿回传敏感信息
	NotifyUrl            string               `json:"notify_url,omitempty"`    // notify_url	string[1,256]	是	是	body json	通知URL必须为直接可访问的URL，不允许携带查询串。
	Sign                 string               `json:"sign,omitempty"`          // sign	string	是	否	body json	开发者对核心字段签名, 防止传输过程中出现意外，签名方式见附录
	SettleAmount         int                  `json:"settle_amount,omitempty"` // settle_amount	number	否	是	body json	当次结算金额，需传大于0的金额，单位为【分】；不传默认全额结算
	MultiCopiesGoodsInfo MultiCopiesGoodsInfo `json:"multi_copies_goods_info"` // multi_copies_goods_info	string[1, 500]	否(单商品多份场景必填)	是	body json	单商品购买多份场景，示例值：[{"copies":2}]， 内容见multi_copies_goods_info字段说明
}

// SettleResponse 结算接口返回值
type SettleResponse struct {
	Result   int    `json:"result,omitempty"`
	ErrorMsg string `json:"error_msg,omitempty"`
	SettleNo string `json:"settle_no,omitempty"`
}

// Settle 请求结算接口
func (k *KuaiShou) Settle(settleParams SettleParams) (settleResponse SettleResponse, err error) {
	token, _ := k.AccessToken.GetAccessToken()
	params, _ := util.JsonStructToMap(settleParams)
	params["multi_copies_goods_info"] = ""
	if settleParams.MultiCopiesGoodsInfo.Copies > 0 {
		params["multi_copies_goods_info"], _ = json.Marshal(settleParams.MultiCopiesGoodsInfo)
	}
	// 开始请求api
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", settle, k.AppId, token)
	postJSON, err := util.PostJSON(api, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(postJSON, &settleResponse)
	if err != nil {
		return
	}
	return
}

// QuerySettleResponse 结算查询结果
type QuerySettleResponse struct {
	Result     int        `json:"result,omitempty"`
	ErrorMsg   string     `json:"error_msg,omitempty"`
	SettleInfo SettleInfo `json:"settle_info,omitempty"`
}

type SettleInfo struct {
	SettleNo     string `json:"settle_no,omitempty"`
	TotalAmount  int    `json:"total_amount,omitempty"`
	SettleAmount int    `json:"settle_amount,omitempty"`
	SettleStatus string `json:"settle_status,omitempty"`
	KsOrderNo    string `json:"ks_order_no,omitempty"`
	KsSettleNo   string `json:"ks_settle_no,omitempty"`
}

// QuerySettle 结算结果查询
func (k *KuaiShou) QuerySettle(outSettleNo string) (querySettleResponse QuerySettleResponse, err error) {
	token, _ := k.AccessToken.GetAccessToken()
	params := map[string]interface{}{
		"out_settle_no": outSettleNo,
	}
	// 开始请求api
	api := fmt.Sprintf("%s?app_id=%s&access_token=%s", querySettle, k.AppId, token)
	postJSON, err := util.PostJSON(api, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(postJSON, &querySettleResponse)
	if err != nil {
		return
	}
	return
}

// SettleCallbackResponse 结算回调参数解析
type SettleCallbackResponse struct {
	Data      SettleCallbackResponseData `json:"data,omitempty"`
	BizType   string                     `json:"biz_type,omitempty"`
	MessageId string                     `json:"message_id,omitempty"`
	AppId     string                     `json:"app_id,omitempty"`
	Timestamp int                        `json:"timestamp,omitempty"`
}

type SettleCallbackResponseData struct {
	OutSettleNo     string `json:"out_settle_no,omitempty"`
	Attach          string `json:"attach,omitempty"`
	SettleAmount    int    `json:"settle_amount,omitempty"`
	Status          string `json:"status,omitempty"`
	KsOrderNo       string `json:"ks_order_no,omitempty"`
	KsSettleNo      string `json:"ks_settle_no,omitempty"`
	EnablePromotion bool   `json:"enable_promotion,omitempty"`
	PromotionAmount int    `json:"promotion_amount,omitempty"`
}

// SettleCallbackResponse 结算结果参数解析
func (k *KuaiShou) SettleCallbackResponse(oldSign, body string, checkSign bool) (settleCallbackResponse SettleCallbackResponse, err error) {
	if checkSign {
		err = k.CallbackCheckSignature(oldSign, body)
		if err != nil {
			return
		}
	}
	// 开始解析参数
	err = json.Unmarshal([]byte(body), &settleCallbackResponse)
	if err != nil {
		return
	}
	return
}
