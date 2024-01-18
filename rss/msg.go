package rss

import (
	"github.com/Chestnuts4/citrix-update-monitor/bot"
	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/mmcdole/gofeed"
)

func buildMsg(title string, link string, date string) *config.Msg {
	m := &config.Msg{}
	m.Title = title
	m.Link = link
	m.Date = date
	return m
}

func sendMsg(items []*gofeed.Item) error {
	for _, item := range items {
		go bot.SendMsgAllBots(buildMsg(item.Title, item.Link, item.Published))
	}
	return nil
}
