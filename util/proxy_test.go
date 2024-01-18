package util

import (
	"testing"
)

func TestBuildClientWithProxy(t *testing.T) {
	httpProxy := "http://192.168.59.1:10809"
	socks5Proxy := "socks5://192.168.59.1:10808"
	testURL := "https://www.citrix.com/content/citrix/en_us/downloads/citrix-adc.rss"

	// Test HTTP proxy
	client, err := BuildClientWithProxy(httpProxy)
	if err != nil {
		t.Errorf("Failed to build client with HTTP proxy: %v", err)
	} else {
		resp, err := client.Get(testURL)
		if err != nil || resp.StatusCode != 200 {
			t.Errorf("Failed to request with HTTP proxy: %v", err)
		}
	}

	// Test SOCKS5 proxy
	client, err = BuildClientWithProxy(socks5Proxy)
	if err != nil {
		t.Errorf("Failed to build client with SOCKS5 proxy: %v", err)
	} else {
		resp, err := client.Get(testURL)
		if err != nil || resp.StatusCode != 200 {
			t.Errorf("Failed to request with SOCKS5 proxy: %v", err)
		}
	}

}
