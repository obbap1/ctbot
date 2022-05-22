package main

import (
	"github.com/obbap1/ctbot/bot"
	"github.com/obbap1/ctbot/db"
)

func main() {
	// initialize database
	d := db.Init()
	// initialize bot
	bot.Init(d)
}
