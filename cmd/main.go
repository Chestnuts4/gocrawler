package main

import (
	"log"

	"github.com/Chestnuts4/citrix-update-monitor/bot"
	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/Chestnuts4/citrix-update-monitor/rss"
)

func main() {
	log.Println("started")
	config.LoadConf(config.ConfPath)
	go rss.StartMonitor()
	bot.StartBots()
}
