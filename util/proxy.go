package util

import (
	"fmt"
	"net/http"
	"net/url"

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
