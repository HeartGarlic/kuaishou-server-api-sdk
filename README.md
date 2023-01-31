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

#### 2. 订单信息查询
    order, err := kuaiShou.QueryOrder("123013100433623410019")

#### 3. 支付回调验签
    // 尚未测试
    err := kuaiShou.PayCallbackCheckSignature("123", "12312321")
	if err != nil {
		t.Errorf("PayCallbackCheckSignature got a error %s", err.Error())
		return
	}