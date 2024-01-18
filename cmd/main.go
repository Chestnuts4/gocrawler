package main

import (
	"fmt"

	"github.com/Chestnuts4/citrix-update-monitor/bot"
	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/Chestnuts4/citrix-update-monitor/rss"
)

func main() {
	fmt.Println("bot started")
	config.LoadConf()
	go rss.Start()
	bot.Start()
}
