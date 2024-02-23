package bot

import (
	"context"
	"log"

	"github.com/Chestnuts4/gocrawler/config"
	"github.com/Chestnuts4/gocrawler/util"
	tb "gopkg.in/telebot.v3"
)

var TgGroups = make(map[string]int64)

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
				err := t.sendMsgAllGroups(msg)
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
	log.Printf("%s started", t.name)
	t.bot.Start()
	return nil
}

func registerHandle(bot *tb.Bot) {
	bot.Handle(tb.OnText, OnTextHandle)

	bot.Handle(TG_PING, func(c tb.Context) error {
		saveGroupId(c)
		return c.Send("pong")
	})
	bot.Handle(TG_START, StartHandle)
	bot.Handle(TG_HELP, HelpHandle)
	bot.Handle(tb.OnAddedToGroup, OnAddedToGroupHandle)
}

func (t *TgBot) sendMsgAllGroups(msg string) error {

	// 遍历群组，给每个群组发消息
	for _, v := range TgGroups {
		_, err := t.bot.Send(&tb.Chat{ID: v}, msg)
		if err != nil {
			return err
		}
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
