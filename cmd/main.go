package main

import (
	"fmt"
	"log"
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

	b.Handle(tb.OnText, func(c tb.Context) error {

		// 判断在艾特自己
		if c.Message().Text == "@"+b.Me.Username {
			// 回复消息
			err := c.Send("Hello, I'm a bot!")
			if err != nil {
				log.Println(err)
			}
		}
		return nil

	})

	b.Start()
}
