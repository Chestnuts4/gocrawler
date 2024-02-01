package main

import (
	"log"

	"github.com/Chestnuts4/gocrawler/bot"
	"github.com/Chestnuts4/gocrawler/config"
	"github.com/Chestnuts4/gocrawler/rss"
)

func main() {
	log.Println("started")
	config.LoadConf(config.ConfPath)
	go rss.StartMonitor()
	bot.StartBots()
}
