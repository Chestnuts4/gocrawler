package bot

import (
	"context"
	"testing"

	"github.com/Chestnuts4/gocrawler/config"
)

func TestLanxinBotStart(t *testing.T) {
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()
	config.LoadConf(config.ConfPath)
	lanxinSecret := config.GlobalConfig.Get("lanxin.secret").(string)
	lanxinWebHook := config.GlobalConfig.Get("lanxin.webhook").(string)
	lanxinProxy := config.GlobalConfig.Get("lanxin.proxy").(string)
	lanxinbot, err := NewLangxinBot(lanxinSecret, lanxinWebHook, lanxinProxy)
	if err != nil {
		t.Fatalf("NewTgbot error: %v", err)
	}
	lanxinbot.Start(ctx)
}
