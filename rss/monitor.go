package rss

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
		Ctx:       ctx,
		Cancel:    cancel,
	}
}

func (m *Monitor) Start() {
	log.Printf("Start monitor %s", m.Url)
	go func() {
		for {
			select {
			case <-m.Ctx.Done():
				return
			default:
				updateItems, err := m.checkFeedUpdate()
				if err != nil {
					log.Printf("Failed to check feed update: %v", err)
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
				sendMsg(items)
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

	resp, err := client.Get(m.Url)
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
