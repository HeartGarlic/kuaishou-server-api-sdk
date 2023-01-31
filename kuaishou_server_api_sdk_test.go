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
	res, err := kuaiShou.Code2Session("0F937BAD052278250C5DAFCACE1B6FCE7C7780C27CAF4A094F972553BBCB2137")
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
		NotifyUrl:   "https://test-api.sylangyue.xyz",
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
	order, err := kuaiShou.QueryOrder("123013100433623410019")
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
func TestKuaiShou_PayCallbackCheckSignature(t *testing.T) {
	err := kuaiShou.PayCallbackCheckSignature("123", "12312321")
	if err != nil {
		t.Errorf("PayCallbackCheckSignature got a error %s", err.Error())
		return
	}
}
