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

func sendItemsToBots(items []*gofeed.Item) error {
	if len(items) == 0 {
		return nil
	}
	if len(items) < 10 {
		for _, item := range items {
			bot.SendMsgAllBots(buildMsg(item.Title, item.Link, item.Published))
		}
	}

	return nil
}
