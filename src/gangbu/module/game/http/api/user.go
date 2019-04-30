package api

import (
	"aliens/log"
	"fmt"
	"gangbu/constant"
	"gangbu/exception"
	"gangbu/module/game/http/helper"
	"gangbu/module/game/service"
	"gangbu/module/game/service/wx"
	"gangbu/module/game/service/wx/model"
	"github.com/gin-gonic/gin"
	"strings"
)

/**
 * @api {GET} /users/ 获取用户信息
 * @apiGroup Users
 * @apiParam {Number} code 微信的code
 * @apiSuccess {string} token 用户的token
 * @apiSuccess {string} nickname 用户昵称
 * @apiSuccess {string} avatar 用户头像地址
 * @apiSuccess {Number} errcode 错误码
 * @apiSuccess {String} errmsg 错误信息
 */
type GetUsersReq struct {
	Code 		string  `form:"code" binding:"required"`
	Channel 	int32  `form:"channel" binding:"required"`
}

type GetUsersResp struct {
	UID   	 int32  `json:"uid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Token    string `json:"token"`
}

func Login(c *gin.Context) {

	req := &GetUsersReq{}
	//if err := c.ShouldBind(req); err != nil {
	//	log.Error(err.Error())
	//	exception.GameException(exception.InvalidParam)
	//}
	helper.CheckReq(c, req)

	//log.Info("bind success")
	user := service.LoginByCode(req.Code, req.Channel)

	token := service.GenerateToken(user.ID)

	resp := &GetUsersResp{}
	resp.UID = user.ID
	resp.Token = token
	resp.Nickname = user.Nickname
	resp.Avatar = user.Avatar

	helper.ResponseWithData(c, resp)
}


/**
 * @api {PUT} /users/ 更新用户头像昵称
 * @apiGroup Users
 * @apiParam {Number} uid 用户的uid
 * @apiParam {String} token 用户的token
 * @apiParam {String} nickname 用户的昵称
 * @apiParam {String} avatar 用户头像地址
 * @apiSuccess {Number} errcode 错误码
 * @apiSuccess {String} errmsg 错误信息
 */
type PutUsersReq struct {
	Nickname string `form:"nickname" binding:"required"`
	Avatar   string `form:"avatar" binding:"required"`
}

func UpdateUser(c *gin.Context) {
	tUser := helper.GetClaimUser(c)
	req := &PutUsersReq{}
	if err := c.ShouldBind(req); err != nil {
		log.Error(err.Error())
		exception.GameException(exception.InvalidParam)
	}
	service.UpdateUser(tUser.UID, req.Nickname, req.Avatar)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

type GetGiftResp struct {
	Status	int32 `json:"status"`
	Message string `json:"message"`
	GiftCode string `json:"giftCode"`
}

func GetGift(c *gin.Context) {
	status, message, giftCode := service.GetGumBallsGift(helper.GetClaimUser(c).UID)
	resp := &GetGiftResp{
		Status:status,
		Message:message,
		GiftCode:giftCode,
	}
	helper.ResponseWithData(c, resp)
}

type DeleteUserReq struct {
	Uid		int32 `form:"uid" binding:"required"`
}

func DeleteUser(c *gin.Context) {
	req := &DeleteUserReq{}
	helper.CheckReq(c, req)
	service.DeleteUser(req.Uid)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

type GetFlagResp struct {
	Flags []*Flag `json:"flags"`
}

type Flag struct {
	Flag constant.ROLE_FLAG `json:"flag"`
	Value int32 `json:"value"`
	UpdateTime int64 `json:"updateTime"`
}

func GetFlags(c *gin.Context) {
	flags := service.GetFlags(helper.GetClaimUser(c).UID)
	flags_resp := []*Flag{}
	for _, flag := range flags {
		flags_resp = append(flags_resp, &Flag{Flag:flag.ID.Flag, Value:flag.Value, UpdateTime:flag.UpdateTime.Unix()})
	}
	resp := &GetFlagResp{
		Flags: flags_resp,
	}
	helper.ResponseWithData(c, resp)
}

type UseFlagReq struct {
	Flag constant.ROLE_FLAG `form:"flag"`
}

func UseFlag(c *gin.Context) {
	req := &UseFlagReq{}
	helper.CheckReq(c,req)
	service.UseFlag(helper.GetClaimUser(c).UID, req.Flag)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

func Wx(c *gin.Context) {
	c.Request.ParseForm()
	form := c.Request.Form
	timestamp := strings.Join(form["timestamp"], "")
	nonce := strings.Join(form["nonce"],"")
	signature := strings.Join(form["signature"], "")
	encryptType := strings.Join(form["encrypt_type"], "")
	msgSignature := strings.Join(form["msg_signature"], "")
	echoStr := strings.Join(form["echostr"], "")

	if !wx.ValidateUrl(timestamp, nonce, signature) {
		log.Info("Wechat Service: this http request is not from Wechat platform!")
		return
	}

	//第一次验证url需要直接返回echoStr
	if echoStr != "" {
		fmt.Fprintf(c.Writer, echoStr)
	}

	var requestBody *model.TextRequestBody

	isEncrypt := encryptType == "aes"

	if isEncrypt {
		encryptRequestBody := wx.ParseEncryptRequestBody(c.Request)
		// Validate mstBody signature
		if !wx.ValidateMsg(timestamp, nonce, encryptRequestBody.Encrypt, msgSignature) {
			log.Error("Wechat Service: validate msg error")
			return
		}
		plainData, err := wx.AesDecrypt(encryptRequestBody.Encrypt)
		if err != nil {
			log.Error("Wechat Service: descrypt err : %v", err)
			return
		}
		//封装struct
		requestBody = wx.ParseEncryptTextRequestBody(plainData)
	} else {
		requestBody = wx.ParseTextRequestBody(c.Request)
	}

	log.Debug("request body : %v", requestBody)

	if requestBody == nil {
		fmt.Fprintf(c.Writer, "success")
		return
	}

	if requestBody.MsgType == "text" && requestBody.Content == "【收到不支持的消息类型，暂无法显示】" {
		//安全模式下向用户回复消息也需要加密
		respBody, e := wx.MakeEncryptResponseBody(requestBody.ToUserName, requestBody.FromUserName, "一些回复给用户的消息", nonce, timestamp)
		if e != nil {
			log.Error("encrypt response error : %v", e)
			return
		}
		fmt.Fprintf(c.Writer, string(respBody))
	}

	openID := requestBody.FromUserName
	if requestBody.MsgType == constant.KEYWORD_EVENT {
		switch requestBody.Event {
		case constant.EVENT_SUBSCRIBE:
			//关注公众号
			//key := wx.HandleClickEvent(requestBody.EventKey, openID)
			wx.HandleOtherEvent(requestBody.Event, openID)
			response := wx.BuildTextResponse(requestBody.ToUserName, openID, wx.GetTextResponse(constant.EVENT_SUBSCRIBE), nonce, timestamp, isEncrypt)
			if response != nil {
				c.Writer.Write(response)
			}
		case constant.EVENT_UNSUBSCRIBE:
			//取消关注公众号
			wx.HandleOtherEvent(requestBody.Event, openID)
		case constant.EVENT_CLICK:
			key := wx.HandleClickEvent(requestBody.EventKey, openID)
			response := wx.BuildTextResponse(requestBody.ToUserName, openID, wx.GetTextResponse(key), nonce, timestamp, isEncrypt)
			if response != nil {
				c.Writer.Write(response)
			} else {
				fmt.Fprintf(c.Writer, "success")
			}
		default:
			//TODO
			fmt.Fprintf(c.Writer, "success")
		}
		//
	}

}