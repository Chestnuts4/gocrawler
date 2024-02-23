package rss

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Chestnuts4/gocrawler/bot"
	"github.com/Chestnuts4/gocrawler/config"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

type Bot struct {
	name     string
	Messages []*config.Msg
}

func (b *Bot) SendMsg(msg *config.Msg) error {
	b.Messages = append(b.Messages, msg)
	return nil
}

func (b *Bot) GetBotName() string {
	return b.name
}

func (b *Bot) Start(ctx context.Context) error {
	// implementation for starting the bot
	return nil
}

func TestSendItemsToBots(t *testing.T) {
	// 重定向 log 输出到标准输出，以便在测试中查看
	log.SetOutput(os.Stdout)

	// 创建一个模拟的 BotSender
	mockBot := &Bot{name: "test"}
	bot.Bots = append(bot.Bots, mockBot)
	// 创建一些模拟的 gofeed.Item
	items := []*gofeed.Item{
		{Title: "Item 1", Description: "Description 1", Link: "http://example.com/1", GUID: "1", Published: "Date 1"},
		{Title: "Item 2", Description: "Description 2", Link: "http://example.com/2", GUID: "2", Published: "Date 2"},
	}

	// 调用 sendItemsToBots 函数
	err := sendItemsToBots(items)

	// 检查是否有错误
	assert.NoError(t, err)

	// 检查是否发送了正确数量的消息
	assert.Equal(t, len(items), len(mockBot.Messages))

	// 检查消息内容
	for i, item := range items {
		expectedMsg := buildMsg(item.Title, item.Description, item.Link, item.GUID, item.Published)
		assert.Equal(t, expectedMsg, mockBot.Messages[i])
	}
}
