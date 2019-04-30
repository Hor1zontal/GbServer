package model

import (
	"encoding/xml"
	"time"
)

type LoginResponse struct {
	OpenID 		string `json:"openid"`
	SessionKey 	string `json:"session_key"`
	UnionID 	string `json:"unionid"`
}

type ErrorResponse struct {
	ErrCode 	int `json:"errcode"`
	ErrMsg 		string `json:"errmsg"`
}

type Config struct {
	Token string			`yaml:"Token"`
	AppID string			`yaml:"AppID"`
	AppSecret string		`yaml:"AppSecret"`
	EncodingAESKey string	`yaml:"EncodingAESKey"`
	AesKey []byte  `json:"id"`
}

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Event 	     string
	EventKey     string
	Url          string
	PicUrl       string
	MediaId      string
	ThumbMediaId string
	Content      string
	MsgId        int
	Location_X   string
	Location_Y   string
	Label        string
}

type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   string
	MsgType      CDATAText
	Content      CDATAText
}

type EncryptRequestBody struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    string
}


type EncryptResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATAText
	MsgSignature CDATAText
	TimeStamp    string
	Nonce        CDATAText
}


type EncryptResponseBody1 struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string
	MsgSignature string
	TimeStamp    string
	Nonce        string
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}