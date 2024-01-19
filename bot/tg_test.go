package bot

import (
	"context"
	"testing"

	"github.com/Chestnuts4/citrix-update-monitor/config"
)

func TestTgBotStart(t *testing.T) {
	config.LoadConf(config.ConfPath)
	tgToken := config.GlobalConfig.Get("tgbot.token").(string)
	tgProxyStr := config.GlobalConfig.Get("tgbot.proxy").(string)
	tgbot, err := NewTgbot(tgToken, tgProxyStr)
	if err != nil {
		t.Fatalf("NewTgbot error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tgbot.Start(ctx)

}
