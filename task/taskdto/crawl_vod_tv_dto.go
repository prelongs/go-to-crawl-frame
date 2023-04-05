package taskdto

import (
	"github.com/gogf/gf/v2/encoding/gjson"
)

// 抓取点播清单接口集合
type CrawlVodTVInterface interface {
	CrawlByBrowserInterface
}

type AbstractCrawlVodTV struct {
	CrawlByBrowserInterface
	*AbstractCrawlByBrowser
}

func (r *AbstractCrawlVodTV) UseBrowser() bool {
	return true
}

func (r *AbstractCrawlVodTV) UseCrawlerProxy() bool {
	return false
}

func (r *AbstractCrawlVodTV) UseBrowserMobProxy() bool {
	return true
}

func (r *AbstractCrawlVodTV) OpenBrowser(ctx *BrowserContext) {
}

func (r *AbstractCrawlVodTV) OpenBrowserWithParams(ctx *BrowserContext, json *gjson.Json) {
}

func (r *AbstractCrawlVodTV) FillTargetRequest(ctx *BrowserContext) {
}
