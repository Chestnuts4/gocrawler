package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/Chestnuts4/citrix-update-monitor/util"
	tb "gopkg.in/telebot.v3"
)

func main() {
	fmt.Println("bot started")
	config.LoadConf()
	token := config.Config.Get("tgbot.token").(string)
	if token == "" {
		log.Println("Set token via environment\nBOT_TOKEN=<your_token>")
		return
	}
	proxyStr := config.Config.Get("tgbot.proxy").(string)
	client, err := util.BuildClientWithProxy(proxyStr)
	if err != nil {
		log.Fatal(err)
		return
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: client,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Bot started[%s]", b.Me.Username)

	// get all groups
	// b.Send(&tb.Chat{ID: }, "Hello, I'm a bot!aaa")

	b.Handle(tb.OnText, func(c tb.Context) error {

		// 判断在艾特自己 contains
		if strings.Contains(c.Message().Text, "@"+b.Me.Username) {
			// 回复消息
			err := c.Send("Hello, I'm a bot!")
			if err != nil {
				log.Println(err)
			}

		}
		// get group id from message
		if c.Message().Chat.Type == tb.ChatGroup {
			log.Printf("group id: %v\nname: %v", c.Message().Chat.ID, c.Message().Chat.Title)

		}
		// get group name from message
		log.Printf("message: %v", c.Message().Text)

		return nil

	})

	b.Start()
}
