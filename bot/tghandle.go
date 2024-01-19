package bot

import (
	"fmt"
	"log"
	"strings"

	tb "gopkg.in/telebot.v3"
)

const (
	TG_START = "/start"
	TG_HELP  = "/help"
	TG_PING  = "/ping"
	TG_TEST  = "/test"
)

func saveGroupId(c tb.Context) {
	// get group id from message
	if c.Message().Chat.Type == tb.ChatGroup {
		id := c.Message().Chat.ID
		title := c.Message().Chat.Title
		if _, ok := TgGroups[title]; !ok {
			TgGroups[title] = id
			log.Printf("group %v %v has in table", title, id)
			c.Send(fmt.Sprintf("group %s %d has in table", title, id))
		}
	}
}

func StartHandle(c tb.Context) error {
	return c.Send("Hello, I'm a bot!")
}

func HelpHandle(c tb.Context) error {
	return c.Send("Hello, I'm a bot!")
}

func OnTextHandle(c tb.Context) error {
	// 判断在艾特自己 contains
	if strings.Contains(c.Message().Text, "@"+c.Chat().Username) {
		// 回复消息
		err := c.Send("Hello, I'm a bot! can I help you?")

		if err != nil {
			log.Println(err)
		}

	}
	// get group id from message
	if c.Message().Chat.Type == tb.ChatGroup {
		saveGroupId(c)

	}
	// get group name from message
	log.Printf("message: %v", c.Message().Text)

	return nil

}

func OnAddedToGroupHandle(c tb.Context) error {
	// get group id from message
	saveGroupId(c)
	return nil
}
