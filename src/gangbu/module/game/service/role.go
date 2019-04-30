package service

import (
	"aliens/common/util"
	"gangbu/constant"
	"gangbu/exception"
	"gangbu/module/game/db"
	"time"
)

func GetRoleInfo(uid int32) *db.DBRole {
	role := db.GetRoleByUid(uid)
	if role == nil {
		newRole := db.NewRoleByUid(uid)
		role = newRole
	}
	return role
}

func UpdateEnergy(uid int32, isAd bool) *db.DBRole {
	role := db.GetRoleByUid(uid)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}
	timestamp := time.Now()
	interval := timestamp.Sub(role.EnergyTime).Seconds()
	restoreTime := constant.ENERGY_RESTORE_TIME // 60秒恢复x点能量
	var ratio int32
	if !isAd {
		ratio = int32(interval / float64( restoreTime ))
		if ratio <= 0  {
			// 不加能量
			return role
		}
	} else {
		//恢复看广告次数
		refreshTime := util.GetTodayHourTime(0)
		if role.LastWatchAd.Before(refreshTime) && time.Now().After(refreshTime) {
			role.AdTimes = constant.MAX_DAY_AD_RESTORE
			db.UpdateOne(role)
		}
		if !CheckAdEnergy(role) {
			return role
		}
		ratio = 1
	}
	if role.TakeInEnergy(ratio*constant.ENERGY_RESTORE_NUM, true) {
		if isAd {
			role.AdTimes --
			role.LastWatchAd = time.Now()
		} else {
			duration := time.Duration(int64(restoreTime) * int64(time.Second))
			if role.Energy == role.EnergyLimit {
				role.EnergyTime = timestamp
			} else {
				role.EnergyTime = role.EnergyTime.Add(duration)
			}
		}
		db.UpdateOne(role)
	}
	return role
}

//查看是否还能通过看广告更新体力
func CheckAdEnergy(role *db.DBRole) bool {
	//小于看广告的间隔时间
	if time.Now().Sub(role.LastWatchAd).Seconds() < constant.WATCH_AD_INTERVAL {
		return false
	}
	//可用看广告恢复体力次数 <= 0
	if role.AdTimes <= 0 {
		return false
	}
	return true
}

func UseEnergy(uid int32) *db.DBRole {
	role := db.GetRoleByUid(uid)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}
	if !role.TakeOutEnergy(constant.ENERGY_COST) {
		exception.GameException(exception.EnergyNotEnough)
	} else {
		//err := role.Update()
		//if err != nil {
		//	exception.GameExceptionCustom("UseEnergy-Update", exception.DatabaseError, err)
		//}
		db.UpdateOne(role)
	}
	return role
}

func PassGuide(uid int32, isEnd bool) {
	role := db.GetRoleByUid(uid)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}
	if isEnd { //新手引导结束
		if !role.Guide {
			role.Guide = true
			db.UpdateOne(role)
		}
	} else {
		flag := db.GetFlag(uid, constant.FLAG_GUIDE, 0)
		flag.SetValue(flag.Value + 1)
		db.UpdateOne(flag)
	}
}

func ReadLetter(uid int32/*, share bool*/) {
	role := db.GetRoleByUid(uid)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}
	flag := db.GetFlag(uid, constant.FLAG_DRAW_LETTER, constant.FLAG_VALUE_FALSE)
	if flag.Value != constant.FLAG_VALUE_FALSE {
		exception.GameException(exception.AlreadyDrawLetter)
	}
	for relicId, relicNum := range constant.CONFIG_LETTER_GIFT {
		AddItem(uid, constant.ITEM_RELIC, relicId, relicNum /** mul*/)
	}
	flag.SetValue(constant.FLAG_VALUE_TRUE)
	db.UpdateOne(flag)
}

func UpdateScore(uid int32, currFloor int32, currScore int32, end bool) {
	role := db.GetRoleByUid(uid)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}
	role.UpdateRoleScore(currScore, currFloor, end)
	db.UpdateOne(role)
	//return role
}

type RankRole struct {
	Nickname	string		`json:"nickname"`
	Avatar		string		`json:"avatar"`
	Score 		int32		`json:"score"`
	Rank		int32 		`json:"rank"` //排行数据
}

func GetRankings(skip int32, limit int32) []*RankRole {
	if limit > constant.RANK_MAX_LIMIT {
		exception.GameException(exception.InvalidParam)
	}
	roles, err := db.GetRolesRank(skip, limit)
	if err != nil {
		exception.GameExceptionCustom("GetRankings-GetRolesRank", exception.DatabaseError, err)
	}
	var result []*RankRole
	for index, role := range roles {
		user := db.GetUserByUid(role.UID)
		if user == nil {
			continue
		}
		result = append(result, &RankRole{Nickname:user.Nickname, Avatar:user.Avatar, Score:role.Score, Rank:int32(index+1) + skip*limit})
	}
	return result
}

func GetRankCount() int32 {
	return db.GetRankCount()
}

func GetRoleRank(uid int32, limit int32) (int32, int32, []*RankRole) {
	role := db.GetRoleByUid(uid)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}

	rankList, rank := GetLimitRankList(limit, role.Score)
	return int32(rank), role.Score, rankList
}

//获取排行榜某个分数的上下limit个的人的具体信息
func GetLimitRankList(limit int32, score int32) ([]*RankRole, int) {
	rank := db.GetScoreRank(score)
	rank += 1
	var rankList []*RankRole
	if limit != 0 {
		gt_roles, lt_roles := db.GetScoreLimitRank(score, limit)
		for index, role := range gt_roles {
			user := db.GetUserByUid(role.UID)
			if user == nil {
				continue
			}
			rankList = append(rankList, &RankRole{Nickname:user.Nickname, Avatar:user.Avatar, Score:role.Score, Rank:int32(rank -(index + 1))})
		}
		for index, role := range lt_roles {
			user := db.GetUserByUid(role.UID)
			if user == nil {
				continue
			}
			rankList = append(rankList, &RankRole{Nickname:user.Nickname, Avatar:user.Avatar, Score:role.Score, Rank:int32(rank + (index + 1))})
		}
	}
	return rankList, rank
}

func GetHelper(uid int32) []int32{
	role := db.GetRoleByUid(uid)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}
	role.UpdateHelpTime()
	db.UpdateOne(role)
	return role.TodayHelper
}