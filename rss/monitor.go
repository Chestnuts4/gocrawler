package rss

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chestnuts4/citrix-update-monitor/config"
	"github.com/Chestnuts4/citrix-update-monitor/util"
	"github.com/mmcdole/gofeed"
)

type Monitor struct {
	Url       string
	Interval  int
	Proxy     string
	LastItems map[string]*gofeed.Item
	// channel
	Updates chan []*gofeed.Item // 添加一个channel字段
	Errors  chan error
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func NewMonitor(url string, interval int, proxy string, ctx context.Context, cancel context.CancelFunc) *Monitor {
	return &Monitor{
		Url:       url,
		Interval:  interval,
		Proxy:     proxy,
		LastItems: make(map[string]*gofeed.Item),
		Updates:   make(chan []*gofeed.Item),
		Errors:    make(chan error),
		Ctx:       ctx,
		Cancel:    cancel,
	}
}

func StartMonitor() {
	url := config.GlobalConfig.Get("monitor.url").(string)
	interval := config.GlobalConfig.Get("monitor.interval").(int)
	proxy := config.GlobalConfig.Get("monitor.proxy").(string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monitor := NewMonitor(url, interval, proxy, ctx, cancel)
	monitor.Start()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}

func (m *Monitor) Start() {
	log.Printf("StartMonitor monitor %s", m.Url)
	go func() {
		for {
			select {
			case <-m.Ctx.Done():
				return
			default:
				updateItems, err := m.checkFeedUpdate()
				if err != nil {
					m.Errors <- err
					continue
				}
				m.Updates <- updateItems

				time.Sleep(time.Duration(m.Interval) * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-m.Ctx.Done():
				return
			case items := <-m.Updates:
				err := sendItemsToBots(items)
				if err != nil {
					m.Errors <- err
				}
			}
		}
	}()
}

func (m *Monitor) checkFeedUpdate() ([]*gofeed.Item, error) {
	rssContent, err := m.downloadRSS()
	fp := gofeed.NewParser()

	if err != nil {
		return nil, fmt.Errorf("failed to download RSS: %v", err)
	}

	feed, err := fp.ParseString(string(rssContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS: %v", err)
	}
	var updateItems []*gofeed.Item

	// 如果LastItems为空，说明是第一次运行，直接将所有item加入LastItems
	if len(m.LastItems) == 0 {
		for _, item := range feed.Items {
			m.LastItems[item.GUID] = item
			updateItems = append(updateItems, item)
		}
		return updateItems, nil
	}

	// 比较LastItems和feed.Items，如果有新的item，就发送消息
	for _, item := range feed.Items {
		// 如果LastItems中没有这个item，说明是新的item
		if _, exists := m.LastItems[item.GUID]; !exists {
			// log.Printf("New item: %s - %s\n", item.Title, item.Link)
			m.LastItems[item.GUID] = item
			updateItems = append(updateItems, item)
		}
	}
	return updateItems, nil
}

func (m *Monitor) downloadRSS() ([]byte, error) {
	client, err := util.BuildClientWithProxy(m.Proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to build client with proxy: %v", err)
	}

	req, err := http.NewRequest("GET", m.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download RSS: %v", err)
	}
	defer resp.Body.Close()

	rssData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSS data: %v", err)
	}

	return rssData, nil
}
