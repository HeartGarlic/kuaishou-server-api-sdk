package kuaishou_server_api_sdk

import (
	"fmt"
	"testing"
	"time"
)

// 声明测试所用的小程序的AppId 与 秘钥
const (
	AppId     = ""
	AppSecret = ""
)

// kuaiShou 快手实例
var (
	kuaiShou *KuaiShou
)

// init 初始化一个快手实例
func init() {
	kuaiShou = NewKuaiShou(&KuaiShouAppletConfig{
		AppId:     AppId,
		AppSecret: AppSecret,
	})
}

// TestKuaiShou_Code2Session 测试Code2Session登陆
func TestKuaiShou_Code2Session(t *testing.T) {
	res, err := kuaiShou.Code2Session("0F937BAD052278250C5DAFCACE1B6FCEE81780C27CAF4A094CAD357C7F67787C")
	if err != nil {
		t.Errorf("code2Sessing got a error %s", err.Error())
		return
	}
	if res.Result != 1 {
		t.Errorf("code2Sessing got a error %s", res.ErrorMsg)
		return
	}
	t.Logf("code2Sessing got OpenId: %s", res.OpenId)
}

// TestKuaiShou_PayCreateOrder 测试支付预下单 目前只测试了有收银台版本
func TestKuaiShou_PayCreateOrder(t *testing.T) {
	// 测试用户的openId
	params := PayCreateOrderParams{
		OutOrderNo:  fmt.Sprintf("%s", time.Now().Format("20060102150405")),
		OpenId:      "f18f5a8e7a3bb15614bf57244ac594f9",
		TotalAmount: 1,
		Subject:     "爽豆充值",
		Detail:      "爽豆充值",
		Type:        1233,
		ExpireTime:  300,
		NotifyUrl:   "",
	}
	res, err := kuaiShou.PayCreateOrder(params)
	if err != nil {
		t.Errorf("PayCreateOrder got a error %s", err.Error())
		return
	}
	if res.Result != 1 {
		t.Errorf("PayCreateOrder got a error %s", res.ErrorMsg)
		return
	}
	t.Logf("PayCreateOrder got value %+v", res)
}

// TestKuaiShou_QueryOrder 测试订单查询接口
func TestKuaiShou_QueryOrder(t *testing.T) {
	order, err := kuaiShou.QueryOrder("123013110250639679019")
	if err != nil {
		t.Errorf("QueryOrder got a error %s code=%d", err.Error(), order.Result)
		return
	}
	if order.Result != 1 {
		t.Errorf("QueryOrder got a error %s", order.ErrorMsg)
		return
	}
	t.Logf("QueryOrder got a value %+v", order)
}

// PayCallbackCheckSignature 回调验签
func TestKuaiShou_CallbackCheckSignature(t *testing.T) {
	err := kuaiShou.CallbackCheckSignature("49f3189f85f7019d33b40b546c87d16a", "123")
	if err != nil {
		t.Errorf("CallbackCheckSignature got a error %s", err.Error())
		return
	}
}

// 增加回调参数解析
func TestKuaiShou_PayCallbackResponse(t *testing.T) {
	jsonStr := "{\"data\":{\"channel\":\"WECHAT\",\"out_order_no\":\"1627293310922demo\",\"attach\":\"小程序demo得\",\"status\":\"SUCCESS\",\"ks_order_no\":\"121112500031787702250\",\"order_amount\":1,\"trade_no\":\"4323300968202201201545417324\",\"extra_info\":\"\",\"enable_promotion\":true,\"promotion_amount\":1},\"biz_type\":\"PAYMENT\",\"message_id\":\"fa578923-347b-4158-9ae8-06c54d485da3\",\"app_id\":\"ks682576822728417112\",\"timestamp\":1627293368719}"
	response, err := kuaiShou.PayCallbackResponse("123", jsonStr, false)
	if err != nil {
		t.Errorf("PayCallbackResponse got a error %s", err.Error())
		return
	}
	if response.Data.Status != "SUCCESS" {
		t.Errorf("PayCallbackResponse got a error %s", response.Data.Status)
		return
	}
	t.Logf("PayCallbackResponse got a value %+v", response)
}

// 支付退款接口
func TestKuaiShou_ApplyRefund(t *testing.T) {
	params := ApplyRefundParams{
		OutOrderNo:           "123456",
		OutRefundNo:          "123456",
		Reason:               "申请退款",
		NotifyUrl:            "",
		RefundAmount:         0,
		Sign:                 "",
		MultiCopiesGoodsInfo: MultiCopiesGoodsInfo{},
	}
	refund, err := kuaiShou.ApplyRefund(params)
	if err != nil {
		t.Errorf("ApplyRefund got a error %s", err.Error())
		return
	}
	if refund.Result != 1 {
		t.Errorf("ApplyRefund got a error %s", refund.ErrorMsg)
		return
	}
	t.Logf("ApplyRefund got a value %+v", refund)
}
