package api

import (
	"encoding/json"
	"gangbu/constant"
	"gangbu/exception"
	"gangbu/module/game/http/helper"
	"gangbu/module/game/service"
	"github.com/gin-gonic/gin"
)

/*-------------------------------------------prop-----------------------------------------------*/

//type Prop struct {
//	ID	int32 `json:"id"`
//	Num		int32 `json:"num"`
//}
//
//type GetPropResp struct {
//	Props	[]*Prop	`json:"props"`
//}
//
//func GetProps(c *gin.Context) {
//	props := service.GetProp(helper.GetClaimUser(c).UID)
//	resp := &GetPropResp{Props:[]*Prop{}}
//	for _, prop := range props {
//		resp.Props = append(resp.Props, &Prop{Id:prop.ID.Type,Num:prop.Num})
//	}
//	helper.ResponseWithData(c, resp)
//}

//type AddPropReq struct {
//	//Uid     int32 		`json:"uid"`
//	Prop	*Prop		`form:"prop"`
//}

//type Prop struct {
//	Type	int32	`json:"type"`
//	Num 	int32 	`json:"num"`
//}

//func AddProp(c *gin.Context) {
//	req := &AddPropReq{}
//	//if err := c.ShouldBindJSON(req); err != nil {
//	//	log.Error(err.Error())
//	//	exception.GameException(exception.InvalidParam)
//	//}
//	helper.CheckReq(c, req)
//	//service.AddProp(req.Uid, req.PropType)
//	//for _, prop := range req.Prop {
//	//	service.AddProp(prop)
//	//}
//	//if req.Uid == 0 {
//	service.AddProp(helper.GetClaimUser(c).UID, req.Prop.Type, req.Prop.Num)
//	//}
//	helper.ResponseWithCode(c, exception.CodeSuccess)
//}

//type UsePropReq struct {
//	Prop	*Prop	`form:"prop"`
//}
//
//func UseProp(c *gin.Context) {
//	req := &UsePropReq{}
//	helper.CheckReq(c, req)
//	service.UseProp(helper.GetClaimUser(c).UID, req.Prop.Type, req.Prop.Num)
//	helper.ResponseWithCode(c, exception.CodeSuccess)
//}


type AddHelpPropReq struct {
	Uid		int32 		  `form:"uid"`
	//Prop	*service.Item `form:"prop"`
}

type AddHelpPropResp struct {
	Result bool `json:"result"`
}

func AddHelpProp(c *gin.Context) {
	req := &AddHelpPropReq{}
	helper.CheckReq(c, req)
	helpID := helper.GetClaimUser(c).UID
	resp := &AddHelpPropResp{}
	if helpID != req.Uid && service.CanAddHelp(helpID, req.Uid) {
		service.AddItem(req.Uid, constant.ITEM_PROP, constant.RELIVE_ITEM_ID, constant.HELP_ADD_NUM)
		resp.Result = true
	} else {
		resp.Result = false
	}
	helper.ResponseWithData(c, resp)
}

/*-------------------------------------------item-----------------------------------------------*/

type AddItemReq struct {
	Type 		int32 `form:"type"`
	ID 			int32 `form:"id"`
	Num			int32 `form:"num"`
}

func AddItem(c *gin.Context) {
	req := &AddItemReq{}
	helper.CheckReq(c, req)
	service.AddItem(helper.GetClaimUser(c).UID, req.Type, req.ID, req.Num)
	helper.ResponseWithCode(c,exception.CodeSuccess)
}

type GetItemsReq struct {
	Type		int32 `form:"type"`
}

type GetItemsResp struct {
	Items		[]*service.Item `json:"items"`
}

func GetItems(c *gin.Context) {
	req := &GetItemsReq{}
	helper.CheckReq(c, req)
	items := service.GetItems(helper.GetClaimUser(c).UID, req.Type)
	resp := &GetItemsResp{
		Items:items,
	}
	helper.ResponseWithData(c, resp)
}

type UseItemReq struct {
	Type	int32 `form:"type"`
	ID 			int32 `form:"id"`
	Num			int32 `form:"num"`
}

func UseItem(c *gin.Context) {
	req := &UseItemReq{}
	helper.CheckReq(c, req)
	service.UseItem(helper.GetClaimUser(c).UID, req.Type, req.ID, req.Num)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

type ComposeEquipReq struct {
	Relics	string `form:"relics"`
	EquipID int32 `form:"equipID"`
}

func ComposeEquip(c *gin.Context) {
	req := &ComposeEquipReq{}
	helper.CheckReq(c, req)
	relics := []int32{}
	json.Unmarshal([]byte(req.Relics), &relics)
	service.CheckComposeEquip(relics, req.EquipID)
	uid := helper.GetClaimUser(c).UID
	service.UseItems(uid, constant.ITEM_RELIC, relics)
	service.AddItem(uid, constant.ITEM_EQUIPMENT, req.EquipID, 1)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}