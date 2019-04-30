package api

import (
	"gangbu/exception"
	"gangbu/module/game/http/helper"
	"gangbu/module/game/service"
	"github.com/gin-gonic/gin"
)

type GetNoticeResp struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetNotice(c *gin.Context) {
	resp := &GetNoticeResp{}
	notice := service.GetNotice()
	if notice != nil {
		resp.ID = notice.ID
		resp.Title = notice.Title
		resp.Content = notice.Content
	}
	helper.ResponseWithData(c, resp)
}

type PubNoticeReq struct {
	Title string `form:"title"`
	Context string `form:"context"`
	PubTime string `form:"pubTime"`
}

func PubNotice(c *gin.Context) {
	req := &PubNoticeReq{}
	helper.CheckReq(c, req)
	service.PubicNotice(req.Title, req.Context, req.PubTime)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}