package rss

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	monitor := &Monitor{
		Url:      "https://www.citrix.com/content/citrix/en_us/downloads/citrix-adc.rss",
		Interval: 10,
		Proxy:    "http://192.168.59.1:8083",
	}

	monitor.Start()
}

func TestCheckFeedUpdate(t *testing.T) {

	rss1, err := os.ReadFile("citrix-adc.rss")
	if err != nil {
		t.Fatalf("Failed to read RSS file: %v", err)
	}

	rss2, err := os.ReadFile("citrix-adc2.rss")
	if err != nil {
		t.Fatalf("Failed to read RSS file: %v", err)
	}
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestCount == 0 {
			w.Write(rss1)
		} else if requestCount == 1 {
			w.Write(rss2)
		} else {
			w.Write(rss2)
		}
		requestCount++
	}))
	// 关闭服务器
	defer server.Close()

	monitor := NewMonitor(server.URL, 10, "http://192.168.59.1:8083", nil, nil)
	fmt.Printf("Server URL: %s\n", server.URL)

	// 循环三次，调用CheckFeedUpdate
	for i := 0; i < 3; i++ {
		items, err := monitor.checkFeedUpdate()
		if err != nil {
			t.Fatalf("Failed to check feed update: %v", err)
		}
		if i == 1 && len(items) != 1 {
			t.Error("Failed to check feed update: items is not 1")
		}
		if i == 2 && len(items) != 0 {
			t.Error("Failed to check feed update: items is not empty")
		}
		if len(items) != 0 && i != 0 {
			for _, item := range items {
				fmt.Printf("Title: %s\n", item.Title)
				fmt.Printf("Link: %s\n", item.Link)
				fmt.Printf("Date: %s\n", item.Published)
			}
			time.Sleep(time.Duration(monitor.Interval) * time.Second)
		}
	}
}

func TestDownloadRSS(t *testing.T) {
	monitor := NewMonitor("https://www.citrix.com/content/citrix/en_us/downloads/citrix-adc.rss", 10, "192.168.59.1:8083", nil, nil)

	monitor.downloadRSS()
}
