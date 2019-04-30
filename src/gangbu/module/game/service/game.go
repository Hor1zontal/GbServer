package service

import (
	"gangbu/constant"
	"gangbu/exception"
	"gangbu/module/game/db"
	"gangbu/module/game/service/lpc"
	"time"
)

func EnterGameData(uid int32/*, props []*db.Prop*/) {
	//初始化，数据清空
	gameData := db.GetGameDataByUid(uid)
	if gameData == nil {
		exception.GameException(exception.GameDataNotFound)
	}
	//初始化带入游戏的道具
	if !gameData.InGame {
		////扣掉道具
		//for _, v := range props {
		//	prop := db.GetPropByUidAndType(uid, v.Type)
		//	if prop == nil {
		//		exception.GameException(exception.PropNotFound)
		//	}
		//	if !prop.TakeOut(prop.Num) {
		//		exception.GameException(exception.PropNotEnough)
		//	}
		//	prop.Update()
		//}
		//初始化带入游戏的道具
		//gameData.Init(props)
		gameData.Clean()
		gameData.InGame = true
		db.UpdateOne(gameData)
	}
}

func UpdateGameData(uid int32, score int32, floor int32, /*props []*db.Prop, boxIDs []int32,*/ end bool, monsterNum, boxNum int32) {
	gameData := db.GetGameDataByUid(uid)
	if gameData == nil {
		exception.GameException(exception.GameDataNotFound)
	}
	if score - gameData.Score < 0 || score - gameData.Score > constant.MAX_FLOOR_SCORE {
		exception.GameException(exception.UpdateScoreError)
	}
	if floor != gameData.Floor + 1 {
		//不是下一层
		exception.GameException(exception.UpdateScoreError)
	}
	lpc.DBServiceProxy.Insert(&db.DBGameLog{
		UID:gameData.UID,
		End:end,
		Floor:gameData.Floor,
		Boxes:gameData.Boxes,
		Score:gameData.Score,
		MonsterNum: gameData.MonsterNum,
		Time:time.Now(),
	}, db.Database.GetHandler())

	if end {
		//游戏结束 游戏数据清零
		gameData.Clean()
	} else {
		//更新分数、道具等信息
		gameData.Update(floor, score/*, props, boxIDs*/, monsterNum, boxNum)
	}
	db.UpdateOne(gameData)
}

func GetGameData(uid int32) *db.DBGameData {
	gameData := db.GetGameDataByUid(uid)
	if gameData == nil {
		gameData = db.NewGameData(uid)
	}
	return gameData
}

func UpdateBoxes(uid int32, boxes []*db.DBBox) {
	gameData := db.GetGameDataByUid(uid)
	if gameData == nil {
		exception.GameException(exception.GameDataNotFound)
	}
	if boxes != nil {
		gameData.Boxes = boxes
	}
	db.UpdateOne(gameData)
}