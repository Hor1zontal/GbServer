package api

import (
	"encoding/json"
	"gangbu/exception"
	"gangbu/module/game/db"
	"gangbu/module/game/http/helper"
	"gangbu/module/game/service"
	"github.com/gin-gonic/gin"
)

//type EnterGameReq struct {
//	Props []*db.Prop	`form:"props"`
//}

type EnterGameResp struct {
	Energy	int32	`json:"energy"`
	EnergyTime int64 `json:"energyTime"`
}

func EnterGame(c *gin.Context) {
	//req := &EnterGameReq{}
	//helper.CheckReq(c, req)
	service.EnterGameData(helper.GetClaimUser(c).UID/*, req.Props*/)
	role := service.UseEnergy(helper.GetClaimUser(c).UID)
	resp := &EnterGameResp{
		Energy:role.Energy,
		EnergyTime:role.EnergyTime.Unix(),
	}
	helper.ResponseWithData(c, resp)
}

type UpdateGameReq struct {
	Floor	int32 `form:"floor"`
	Score	int32 `form:"score"`
	//Props   []*db.Prop `form:"props"`
	//BoxIds 	[]int32 `form:"boxIds"`
	MonsterNum	int32 `form:"monsterNum"`
	BoxNum	int32 `form:"boxNum"`
	End		bool  `form:"end"` //游戏是否结束
}

func UpdateGame(c *gin.Context) {
	req := &UpdateGameReq{}
	helper.CheckReq(c, req)
	service.UpdateGameData(helper.GetClaimUser(c).UID, req.Score, req.Floor, /*req.Props, req.BoxIds,*/ req.End, req.MonsterNum, req.BoxNum)
	service.UpdateScore(helper.GetClaimUser(c).UID, req.Floor, req.Score, req.End)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

type UpdateBoxReq struct {
	Boxes string `form:"boxes"`
}

//type Box struct {
//	PropID	int32 `bson:"propID" json:"propID"`
//	RelicID int32 `bson:"relicID" json:"relicID"`
//	GetTime int32 `bson:"getTime" json:"getTime"`
//}

func UpdateBoxData(c *gin.Context) {
	req := &UpdateBoxReq{}
	helper.CheckReq(c, req)
	boxes := []*db.DBBox{}
	json.Unmarshal([]byte(req.Boxes), &boxes)
	service.UpdateBoxes(helper.GetClaimUser(c).UID, boxes)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

type GetGameResp struct {
	InGame	bool       `json:"inGame"`
	Score   int32        `json:"score"`
	Floor	int32       `json:"floor"`
	Boxes	[]*db.DBBox `json:"boxes"`
	MonsterNum int32 	`json:"monsterNum"`
	BoxNum	int32		`json:"boxNum"`
}

func GetGame(c *gin.Context) {
	gameData := service.GetGameData(helper.GetClaimUser(c).UID)
	resp := &GetGameResp{
		InGame:gameData.InGame,
		Floor:gameData.Floor,
		Score:gameData.Score,
		Boxes:gameData.Boxes,
		MonsterNum:gameData.MonsterNum,
		//Props:gameData.Props,
		BoxNum:gameData.BoxNum,
	}
	helper.ResponseWithData(c, resp)
}