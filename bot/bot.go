package bot

import "github.com/Chestnuts4/citrix-update-monitor/config"

// bot 数组
var Bots []Bot

type Bot interface {
	// 发送消息
	SendMsg(msg *config.Msg) error
	// bot名字
	GetBotName() string
}

func SendMsgAllBots(m *config.Msg) {
	for _, bot := range Bots {
		bot.SendMsg(m)
	}

}
