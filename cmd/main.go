package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"golang.org/x/net/proxy"
	tb "gopkg.in/telebot.v3"
)

func buildClientWithProxy(addr string) (*http.Client, error) {
	if addr != "" {
		u, err := url.Parse(addr)
		if err != nil {
			panic(err)
		}

		var auth *proxy.Auth
		if u.User != nil { // credentials are set
			pass, _ := u.User.Password()
			auth = &proxy.Auth{
				User:     u.User.Username(),
				Password: pass,
			}
		}
		dialer, err := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
		if err != nil {
			return nil, err
		}

		// Patch client transport
		httpTransport := &http.Transport{Dial: dialer.Dial}
		hc := &http.Client{Transport: httpTransport}

		return hc, nil
	}

	return nil, nil // use default
}

func main() {
	fmt.Println("bot started")
	token := config.Config.Get("bot.token").(string)
	if token == "" {
		log.Println("Set token via environment\nBOT_TOKEN=<your_token>")
		return
	}
	proxyStr := config.Config.Get("bot.proxy").(string)
	client, err := buildClientWithProxy(proxyStr)
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

	b.Handle("/start", func(c tb.Context) error {
		return c.Send("aaaa")

	})

	b.Start()
}
