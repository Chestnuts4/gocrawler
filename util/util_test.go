package util

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Chestnuts4/gocrawler/config"
	socks5 "github.com/armon/go-socks5"
	"github.com/elazarl/goproxy"
	"github.com/stretchr/testify/assert"
)

func TestBuildClientWithProxy(t *testing.T) {

	// 创建一个模拟的HTTP服务器
	mockHTTPServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is the HTTP server"))
	}))
	defer mockHTTPServer.Close()

	// 创建一个SOCKS5服务器
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		t.Fatal(err)
	}

	// 启动SOCKS5服务器
	go func() {
		if err := server.ListenAndServe("tcp", "127.0.0.1:58080"); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Second)

	client, err := BuildClientWithProxy("socks5://127.0.0.1:58080")
	if err != nil {
		t.Errorf("Failed to build client with SOCKS5 proxy: %v", err)
	} else {
		resp, err := client.Get(mockHTTPServer.URL)
		if err != nil || resp.StatusCode != 200 {
			t.Errorf("Failed to request with SOCKS5 proxy: %v", err)
		}
		// 判断返回内容是否是this is the HTTP proxy server
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v", err)
		}
		if string(body) != "this is the HTTP server" {
			t.Errorf("Failed to read response body: %s", string(body))
		}
	}

	httpProxy := goproxy.NewProxyHttpServer()
	go func() {
		if err := http.ListenAndServe("127.0.0.1:50801", httpProxy); err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second)

	client2, err := BuildClientWithProxy("http://127.0.0.1:50801")
	if err != nil {
		t.Errorf("Failed to build client with HTTP proxy: %v", err)
	} else {
		resp, err := client2.Get(mockHTTPServer.URL)
		if err != nil || resp.StatusCode != 200 {
			t.Errorf("Failed to request with HTTP proxy: %v", err)
		}
		// 判断返回内容是否是this is the HTTP proxy server
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v", err)
		}
		if string(body) != "this is the HTTP server" {
			t.Errorf("Failed to read response body: %s", string(body))
		}
	}
}

func TestFormatMsg(t *testing.T) {
	msg := &config.Msg{
		Title:       "Test Title",
		Description: "Test Description",
		Link:        "http://test.link",
		Guid:        "Test Guid",
		Date:        "Test Date",
	}

	expected := "title: Test Title\nDesc:Test Description\nLink: http://test.link\nGuid: Test Guid\nDate: Test Date"
	result := FormatMsg(msg)

	assert.Equal(t, expected, result)
}
