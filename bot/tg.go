package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/Chestnuts4/citrix-update-monitor/util"
	tb "gopkg.in/telebot.v3"
)

type TgBot struct {
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
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10},
		Client: client,
	})
	if err != nil {
		return nil, err
	}
	return &TgBot{
		token:   token,
		proxy:   proxy,
		bot:     bot,
		msgChan: make(chan string),
		errChan: make(chan error),
	}, nil
}

// 接收context
func (t *TgBot) Start(ctx context.Context) {

	t.bot.Start()

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
}

func (t *TgBot) sendMsg(msg string) error {

	return nil
}

func (t *TgBot) SendMsg(msg config.Msg) {
	msgStr := fmt.Sprintf("%s\n%s\n%s", msg.Title, msg.Link, msg.Date)
	t.msgChan <- msgStr
}

func (t *TgBot) GetBotName() string {
	return tgBotName
}
