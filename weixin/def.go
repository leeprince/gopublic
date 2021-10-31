package weixin

import wxpay "gopkg.in/go-with/wxpay.v1"

// UserInfo 用户信息
type UserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

// APITicket ...
type APITicket struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// WxInfo 微信配置信息
type WxInfo struct {
	AppID     string // 微信公众平台应用ID
	AppSecret string // 微信支付商户平台商户号
	APIKey    string // 微信支付商户平台API密钥
	MchID     string // 商户号
	NotifyURL string // 通知地址
	ShearURL  string // 微信分享url
}

// TempMsg 订阅消息头
type TempMsg struct {
	Touser     string                       `json:"touser"`      //	是	接收者（用户）的 openid
	TemplateID string                       `json:"template_id"` //	是	所需下发的模板消息的id
	Page       string                       `json:"page"`        //	否	点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	Data       map[string]map[string]string `json:"data"`        //是	模板内容，不填则下发空模板
}

// ResTempMsg 模版消息返回值
type ResTempMsg struct {
	Errcode int    `json:"errcode"` //
	Errmsg  string `json:"errmsg"`
}

type wxTools struct {
	client     *wxpay.Client
	wxInfo     WxInfo
	certFile   string // 微信支付商户平台证书路径
	keyFile    string
	rootcaFile string
}

// QrcodeRet ...
type QrcodeRet struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//
type wxPostdata struct {
	Scene string `json:"scene"`
	Page  string `json:"page"`
}

//
type wxQrcodedata struct {
	Path  string `json:"path"`  //路径
	Width int    `json:"width"` //二维码宽度
}
