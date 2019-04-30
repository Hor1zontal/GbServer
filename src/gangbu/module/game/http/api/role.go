package api

import (
	"aliens/log"
	"gangbu/constant"
	"gangbu/exception"
	"gangbu/module/game/http/helper"
	"gangbu/module/game/service"
	"github.com/gin-gonic/gin"
)

type GetRoleInfoResp struct {
	Score	int32 `json:"score"` //历史最高分
	Floor	int32 `json:"floor"` //历史到达最高关卡
	InGame   bool 	`json:"inGame"` //游戏是否在进行
	Energy	int32	`json:"energy"`
	EnergyLimit int32 `json:"energyLimit"`
	EnergyTime int64 `json:"energyTime"`
	Guide    bool	`json:"guide"`
	LastWatchAd int64 `json:"lastWatchAd"`
	AdTimes 	int32 `json:"adTimes"`
	GiftCode	string `json:"giftCode"`
}

func GetRoleInfo(c *gin.Context) {
	role := service.GetRoleInfo(helper.GetClaimUser(c).UID)
	role = service.UpdateEnergy(helper.GetClaimUser(c).UID, false)
	resp := &GetRoleInfoResp{
		Score:	role.Score,
		Floor: 	role.Floor,
		Energy: role.Energy,
		EnergyLimit: role.EnergyLimit,
		EnergyTime:role.EnergyTime.Unix(),
		Guide:role.Guide,
		LastWatchAd:role.LastWatchAd.Unix(),
		AdTimes:role.AdTimes,
		GiftCode:role.GiftCode,
		//TodayHelper: role.TodayHelper,
	}
	helper.ResponseWithData(c, resp)
}

type GetRankingsReq struct {
	Skip	int32 `form:"skip"`
	Limit	int32 `form:"limit"`

}

type GetRankingsResp struct {
	RankList	[]*service.RankRole		`json:"rankList"`
	Count		int32 `json:"count"`
}

func GetRankings(c *gin.Context) {
	req := &GetRankingsReq{}
	helper.CheckReq(c,req)
	rankRoles := []*service.RankRole{}
	if req.Skip != 0 || req.Limit != 0 {
		rankRoles = service.GetRankings(req.Skip, req.Limit)
	}
	count := service.GetRankCount()
	resp := &GetRankingsResp{
		RankList:rankRoles,
		Count:count,
	}
	helper.ResponseWithData(c, resp)
}

type GetRoleRankReq struct {
	Limit	int32	`form:"limit"`
}

type GetRoleRankResp struct {
	Rank 	int32	`json:"rank"`
	Score 	int32	`json:"score"`
	RankList	[]*service.RankRole		`json:"rankList"` //排名在附近的用户
}

func GetRoleRank(c *gin.Context) {
	req := &GetRoleRankReq{}
	if err := c.ShouldBind(req); err != nil {
		log.Error(err.Error())
		exception.GameException(exception.InvalidParam)
	}
	rank, score, rankList := service.GetRoleRank(helper.GetClaimUser(c).UID, req.Limit)
	resp := &GetRoleRankResp{
		Rank:rank,
		Score:score,
		RankList:rankList,
	}
	helper.ResponseWithData(c, resp)
}

type ReqGetScoreRankReq struct {
	Score	int32 `form:"score"`
}

type ReqGetScoreRankResp struct {
	RankList	[]*service.RankRole		`json:"rankList"` //排名在附近的用户
}

func GetScoreRank(c *gin.Context) {
	req := &ReqGetScoreRankReq{}
	helper.CheckReq(c, req)
	rankList, _ := service.GetLimitRankList(1, req.Score)
	resp := &ReqGetScoreRankResp{
		//Rank:int32(rank),
		RankList:rankList,
	}
	helper.ResponseWithData(c, resp)
}

type UpdateEnergyReq struct {
	IsAd	bool `form:"isAd"`
}

type UpdateEnergyResp struct {
	Energy	int32	`json:"energy"`
	AdTimes int32 `json:"adTimes"`
	EnergyTime int64 `json:"energyTime"`
	LastWatchAd int64 `json:"lastWatchAd"`
}

func UpdateEnergy(c *gin.Context) {
	req := &UpdateEnergyReq{}
	helper.CheckReq(c, req)
	role := service.UpdateEnergy(helper.GetClaimUser(c).UID, req.IsAd)
	resp := &UpdateEnergyResp{
		Energy:	role.Energy,
		EnergyTime: role.EnergyTime.Unix(),
		AdTimes:role.AdTimes,
		LastWatchAd:role.LastWatchAd.Unix(),
	}
	helper.ResponseWithData(c, resp)
}

type PassGuideReq struct {
	End		bool `form:"end"`
}

func PassGuide(c *gin.Context) {
	req := &PassGuideReq{}
	helper.CheckReq(c, req)
	service.PassGuide(helper.GetClaimUser(c).UID, req.End)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

func DrawLetter(c *gin.Context) {
	uid := helper.GetClaimUser(c).UID
	service.ReadLetter(uid/*, req.Share*/)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

func BeatPeter(c *gin.Context) {
	uid := helper.GetClaimUser(c).UID
	service.SetFlagTrue(uid, constant.FLAG_BEAT_PETER)
	helper.ResponseWithCode(c, exception.CodeSuccess)
}

type GetHelperResp struct {
	TodayHelper		[]int32 `json:"todayHelper"`
}

func GetHelper(c *gin.Context) {
	uid := helper.GetClaimUser(c).UID
	todayHelper := service.GetHelper(uid)
	resp := &GetHelperResp{
		TodayHelper:todayHelper,
	}
	helper.ResponseWithData(c, resp)
}