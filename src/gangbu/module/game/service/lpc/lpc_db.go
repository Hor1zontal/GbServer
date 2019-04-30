package lpc

import (
	database2 "aliens/database"
	"gangbu/constant"
	"gangbu/module/database"
)

var DBServiceProxy = &dbHandler{}

type dbHandler struct {

}

func (this *dbHandler) Insert(data interface{}, dbHandler database2.IDatabaseHandler) {
	database.ChanRPC.Go(constant.DB_COMMAND_INSERT, data, dbHandler)
}
