package bot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/prometheus/common/log"
)

// bot 数组
var Bots []Bot

type Bot interface {
	// 发送消息
	SendMsg(msg *config.Msg) error
	// bot名字
	GetBotName() string
	Start(ctx context.Context) error
}

func SendMsgAllBots(m *config.Msg) {
	for _, bot := range Bots {
		bot.SendMsg(m)
	}
}

func StartBots() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bots := config.GlobalConfig.Get("bot.list").([]interface{})
	for _, bot := range bots {
		botMap := bot.(map[string]interface{})
		if tgbot, ok := botMap["tgbot"]; ok {
			tgbotMap := tgbot.(map[string]interface{})
			tgToken := tgbotMap["token"].(string)
			tgProxy := tgbotMap["proxy"].(string)
			tgbot, err := NewTgbot(tgToken, tgProxy)
			if err != nil {
				log.Fatalf("NewTgbot error: %v", err)
			}
			go tgbot.Start(ctx)
			Bots = append(Bots, tgbot)
		}
		if lanxinBot, ok := botMap["lanxin"]; ok {
			lanxinBots := lanxinBot.(interface{})
			for _, lanxinSetting := range lanxinBots.([]interface{}) {
				lanxinBotMap := lanxinSetting.(map[string]interface{})
				lanxinSecret := lanxinBotMap["secret"].(string)
				lanxinWebHook := lanxinBotMap["webhook"].(string)
				lanxinProxy := lanxinBotMap["proxy"].(string)
				lanxinBot, err := NewLangxinBot(lanxinSecret, lanxinWebHook, lanxinProxy)
				if err != nil {
					log.Fatalf("NewTgbot error: %v", err)
				}
				go lanxinBot.Start(ctx)
				Bots = append(Bots, lanxinBot)
			}

		}
	}

	// lanxinList := config.GlobalConfig.Get("bot.list.lanxin").([]interface{})
	// for _, v := range lanxinList {
	// 	lanxin := v.(map[interface{}]interface{})
	// 	lanxinSecret := lanxin["secret"].(string)
	// 	lanxinWebHook := lanxin["webhook"].(string)
	// 	lanxinProxy := lanxin["proxy"].(string)
	// 	lanxinBot, err := NewLangxinBot(lanxinSecret, lanxinWebHook, lanxinProxy)
	// 	if err != nil {
	// 		log.Fatalf("NewTgbot error: %v", err)
	// 	}
	// 	go lanxinBot.Start(ctx)
	// 	Bots = append(Bots, lanxinBot)
	// }
	//
	// lanxinSecret := config.GlobalConfig.Get("lanxin.secret").(string)
	// lanxinWebHook := config.GlobalConfig.Get("lanxin.webhook").(string)
	// lanxinProxy := config.GlobalConfig.Get("lanxin.proxy").(string)
	// lanxinBot, err := NewLangxinBot(lanxinSecret, lanxinWebHook, lanxinProxy)
	// if err != nil {
	// 	log.Fatalf("NewTgbot error: %v", err)
	// }
	// go lanxinBot.Start(ctx)
	// Bots = append(Bots, lanxinBot)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

}
