package service

import (
	"gangbu/exception"
	"gangbu/module/game/config"
	"gangbu/module/game/db"
)

func CanAddHelp(helpID, beHelpId int32) bool {
	role := db.GetRoleByUid(beHelpId)
	if role == nil {
		exception.GameException(exception.RoleNotFound)
	}
	role.UpdateHelpTime()
	if !role.CheckHelp(helpID) {
		return false
	}
	db.UpdateOne(role)
	return true
}

//-----------------------------------------------------------------------------------------------------------------
type Item struct {
	ID		int32 `json:"id" form:"id"`
	Num 	int32 `json:"num" form:"num"`
}

func AddItem(uid, itemType, itemID, num int32) {
	db.AddItem(uid, itemType, itemID, num)
}

func UseItems(uid, itemType int32, itemID []int32) {
	for _, item := range itemID {
		UseItem(uid, itemType, item, 1)
	}
}

func UseItem(uid, itemType, itemId, itemNum int32) {
	db.UseItem(uid, itemType, itemId, itemNum)
}

func GetItems(uid, itemType int32) []*Item{
	items_db := db.GetItemsByType(uid, itemType)
	result := make([]*Item,len(items_db))
	for index, item_db := range items_db {
		result[index] = &Item{ID:item_db.ID.ID,Num:item_db.Num}
	}
	return result
}

func CheckComposeEquip(relics []int32, equipID int32) {
	conf_relics := config.DATA.EquipmentRelicMap[equipID]
	if len(conf_relics) == 0 || len(relics) != len(conf_relics) {
		exception.GameException(exception.ComposeEquipError)
	}
	conf_map := make(map[int32]struct{}, len(conf_relics))
	for _, conf_id := range conf_relics {
		conf_map[conf_id] = struct{}{}
	}
	for _, id := range relics {
		_, ok := conf_map[id]
		if ok {
			delete(conf_map, id)
		} else {
			exception.GameException(exception.ComposeEquipError)
		}
	}
	//if reflect.DeepEqual(relics, conf_relics) {
	//	log.Debug("compose success")
	//}
}