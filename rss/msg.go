package rss

import (
	"log"
	"time"

	"github.com/Chestnuts4/citrix-update-monitor/bot"
	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/mmcdole/gofeed"
)

func buildMsg(title string, desc string, link string, guid string, date string) *config.Msg {
	m := &config.Msg{}
	m.Title = title
	m.Description = desc
	m.Guid = guid
	m.Link = link
	m.Date = date
	return m
}

func sendItemsToBots(items []*gofeed.Item) error {

	// if len(items) >= 10 {
	// 	log.Println("Items length >=10 break")

	// 	return nil
	// }
	// for _, item := range items {
	// 	<-ticker.C // 等待定时器的下一个事件
	// 	log.Printf("Send message to all bots: %s", item.Title)

	// 	bot.SendMsgAllBots(buildMsg(item.Title, item.Description, item.Link, item.GUID, item.Published))
	// }

	ticker := time.NewTicker(time.Second / 5) // 创建一个每 200 毫秒触发一次的定时器
	defer ticker.Stop()
	for _, item := range items {
		<-ticker.C // 等待定时器的下一个事件
		log.Printf("Send message to all bots: %s", item.Title)

		bot.SendMsgAllBots(buildMsg(item.Title, item.Description, item.Link, item.GUID, item.Published))
	}
	return nil
}
