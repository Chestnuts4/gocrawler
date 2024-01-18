package bot

import (
	"context"
	"log"
	"strings"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/Chestnuts4/citrix-update-monitor/util"
	tb "gopkg.in/telebot.v3"
)

type TgBot struct {
	name    string
	token   string
	proxy   string
	bot     *tb.Bot
	msgChan chan string
	errChan chan error
}

const tgBotName = "telegram bot"

func NewTgbot(token string, proxy string) (*TgBot, error) {
	client, err := util.BuildClientWithProxy(proxy)
	if err != nil {
		return nil, err
	}
	bot, err := tb.NewBot(tb.Settings{
		Token: token,
		Poller: &tb.LongPoller{Timeout: 10, AllowedUpdates: []string{"message",
			"chat_member",
			"inline_query",
			"callback_query"}},
		Client: client,
	})
	if err != nil {
		return nil, err
	}

	registerHandle(bot)
	return &TgBot{
		name:    tgBotName,
		token:   token,
		proxy:   proxy,
		bot:     bot,
		msgChan: make(chan string),
		errChan: make(chan error),
	}, nil
}

// 接收context
func (t *TgBot) Start(ctx context.Context) error {

	go func() {
		for {
			select {
			case msg := <-t.msgChan:
				err := t.sendMsg(msg)
				if err != nil {
					t.errChan <- err
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-t.errChan:
				log.Println(err)
			case <-ctx.Done():
				return
			}
		}
	}()
	t.bot.Start()
	return nil
}

func registerHandle(bot *tb.Bot) {
	bot.Handle(tb.OnText, func(c tb.Context) error {
		// 判断在艾特自己 contains
		if strings.Contains(c.Message().Text, "@"+bot.Me.Username) {
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
}

func (t *TgBot) sendMsg(msg string) error {
	_, err := t.bot.Send(&tb.Chat{ID: -4193925869}, msg)
	if err != nil {
		return err
	}
	//log.Printf("send msg: %v", retMsg)
	return nil
}

func (t *TgBot) SendMsg(msg *config.Msg) error {

	msgStr := util.FormatMsg(msg)
	t.msgChan <- msgStr
	return nil
}

func (t *TgBot) GetBotName() string {
	return t.name
}
