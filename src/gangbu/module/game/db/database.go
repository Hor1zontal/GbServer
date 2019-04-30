package db

import (
	"aliens/database"
	"aliens/database/mongo"
	"gangbu/module/game/config"
)

var Database database.IDatabase = &mongo.Database{}
var databaseHandler = Database.GetHandler()

func Init() {
	if config.Server.Database.Name == "" {
		config.Server.Database.Name = "GangBu"
	}
	err := Database.Init(config.Server.Database)
	if err != nil {
		panic(err)
	}

	databaseHandler.EnsureTable("user", &DBUser{})
	databaseHandler.EnsureTable("role", &DBRole{})
	databaseHandler.EnsureTable("loginLog", &DBLoginLog{})
	//databaseHandler.EnsureTable("mail", &DBMail{})
	//databaseHandler.EnsureTable("check", &DBCheck{})
	//databaseHandler.EnsureTable("hero", &DBHero{})
	databaseHandler.EnsureTable("item", &DBItem{})
	databaseHandler.EnsureTable("gameData", &DBGameData{})
	databaseHandler.EnsureTable("gameLog", &DBGameLog{})
	databaseHandler.EnsureTable("notice", &DBNotice{})
	databaseHandler.EnsureTable("flag",&DBFlag{})
	//databaseHandler.EnsureTable("", &DBUserToken{})
	databaseHandler.EnsureIndex("flag", []string{"_id.uid"}, false)
	databaseHandler.EnsureIndex("item", []string{"_id.uid","_id.type"}, false)
}

func Close() {

}