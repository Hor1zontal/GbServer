package app

import (
	"aliens/log"
	"flag"
	"gangbu/module/game/config"
	"gangbu/module/game/words"
	"github.com/name5566/leaf"
	"github.com/name5566/leaf/module"
)

var (
	debug = false
	//tag   = ""
	configPath = ""
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug flag")
	flag.StringVar(&configPath, "config", "conf", "configuration path")
	//flag.StringVar(&tag, "tag", "aliensboot", "log tag")
	flag.Parse()

}

func Run(mods ...module.Module) {
	//rand.Seed(time.Now().UnixNano())
	log.Info("☁☁☁☁☁☁☁☁☁☁☁☁☁☁version:1.0.0.2019032701☁☁☁☁☁☁☁☁☁☁☁☁☁☁")
	config.Init(configPath)
	words.Init()
	leaf.Run(mods...)
}