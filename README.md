## 快手小程序服务端 golang sdk
`没有找到现成的, 只有自己造个轮子了.`
`ps: 第一次写`

### 已实现的接口
#### 0. 初始化快手实例
    kuaiShou = NewKuaiShou(&KuaiShouAppletConfig{
		AppId:     "AppId",
		AppSecret: "AppSecret",
	})

#### 1. 小程序登录
    kuaiShou.Code2Session("0F937BAD052278250C5DAFCACE1B6FCE7C7780C27CAF4A094F972553BBCB2137")

#### 2. 担保支付
#### 1. 支付预下单
	params := PayCreateOrderParams{
		OutOrderNo:  fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), rand.Int()),
		OpenId:      "f18f5a8e7a3bb15614bf57244ac594f9",
		TotalAmount: 1,
		Subject:     "爽豆充值",
		Detail:      "爽豆充值",
		Type:        1233,
		ExpireTime:  300,
	}
	res, err := kuaiShou.PayCreateOrder(params)
#### 1.1 支付回调解析
    jsonStr := "{\"data\":{\"channel\":\"WECHAT\",\"out_order_no\":\"1627293310922demo\",\"attach\":\"小程序demo得\",\"status\":\"SUCCESS\",\"ks_order_no\":\"121112500031787702250\",\"order_amount\":1,\"trade_no\":\"4323300968202201201545417324\",\"extra_info\":\"\",\"enable_promotion\":true,\"promotion_amount\":1},\"biz_type\":\"PAYMENT\",\"message_id\":\"fa578923-347b-4158-9ae8-06c54d485da3\",\"app_id\":\"ks682576822728417112\",\"timestamp\":1627293368719}"
	response, err := kuaiShou.PayCallbackResponse("123", jsonStr, false)
	if err != nil {
		t.Errorf("PayCallbackResponse got a error %s", err.Error())
		return
	}

#### 2. 订单信息查询
    order, err := kuaiShou.QueryOrder("123013100433623410019")

#### 3. 支付回调验签
    // 尚未测试
    err := kuaiShou.PayCallbackCheckSignature("123", "12312321")
	if err != nil {
		t.Errorf("PayCallbackCheckSignature got a error %s", err.Error())
		return
	}