package main


import (
	"gangbu/app"
	"gangbu/module/database"
	"gangbu/module/game"
)

func main() {
	app.Run(
		database.Module,
		game.Module,
	)
}