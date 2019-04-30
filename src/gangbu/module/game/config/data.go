package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var DATA struct{
	EquipmentRelicMap	map[int32][]int32		// equipmentID -- [relicID]
}

func initData() {
	updateFragmentBaseData(ReadJsonData("/data/fragment.json"))
}


func ReadJsonData(path string) []byte {
	rfile, _ := os.Open(GetConfigPath(path))
	byte, _ := ioutil.ReadAll(rfile)
	return byte
}
type FragmentBase struct {
	ID int32	`json:"id"`
	EquipID int32 `json:"equipID"`
}

func updateFragmentBaseData(data []byte) {
	var datas []*FragmentBase
	json.Unmarshal(data, &datas)
	result := make(map[int32][]int32)
	for _, data := range datas {
		result[data.EquipID] = append(result[data.EquipID], data.ID)
	}
	DATA.EquipmentRelicMap = result
}