package bot

import (
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
}

func SendMsgAllBots(m *config.Msg) {
	log.Infof("Send message to all bots: %s", m.Title)
	for _, bot := range Bots {
		bot.SendMsg(m)
	}

}
