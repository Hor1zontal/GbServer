package internal

import (
	"aliens/log"
	"gangbu/module/game/service/wx"
	"github.com/name5566/leaf/timer"
)

func Init() {
	cron, err := timer.NewCronExpr("0 */1 * * *")
	if err != nil {
		log.Error("init wechat accessToken timer error : %v", err)
	}
	//刷新token
	skeleton.CronFunc(cron, wx.RefreshToken)
}
