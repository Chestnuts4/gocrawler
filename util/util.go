package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"golang.org/x/net/proxy"
)

func BuildClientWithProxy(addr string) (*http.Client, error) {
	if addr != "" {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}

		var auth *proxy.Auth
		if u.User != nil { // credentials are set
			pass, _ := u.User.Password()
			auth = &proxy.Auth{
				User:     u.User.Username(),
				Password: pass,
			}
		}
		var dialer proxy.Dialer
		var httpTransport *http.Transport
		switch u.Scheme {
		case "http":

			httpTransport = &http.Transport{Proxy: http.ProxyURL(u)}

		case "socks5":
			dialer, err = proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
			httpTransport = &http.Transport{Dial: dialer.Dial}

		default:
			return nil, fmt.Errorf("unsupported proxy scheme: %s", u.Scheme)
		}

		if err != nil {
			return nil, err
		}

		// Patch client transport
		hc := &http.Client{Transport: httpTransport}

		return hc, nil
	}

	return &http.Client{}, nil
}

func FormatMsg(msg *config.Msg) string {
	formatStr := "title: %s\nDesc:%s\nLink: %s\nGuid: %s\nDate: %s"
	return fmt.Sprintf(formatStr, msg.Title, msg.Description, msg.Link, msg.Guid, msg.Date)
}

func LanxinSign(secret string) string {
	timestamp := time.Now().Unix()
	stringToSign := fmt.Sprintf("%v", timestamp) + "@" + secret
	h := hmac.New(sha256.New, []byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}
