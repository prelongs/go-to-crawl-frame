package taskdto

import "github.com/gogf/gf/v2/encoding/gjson"

// 浏览器打开后操作的接口集合
type CrawlByBrowserInterface interface {
	UseBrowser() bool         // 使用浏览器抓取URL
	UseMobileUA() bool        // 默认使用随机桌面浏览器UA，返回true就使用随机手机端UA
	UseCrawlerProxy() bool    // 数据库管理的爬虫代理池
	UseBrowserMobProxy() bool // 内嵌类似Fiddler的抓包代理, 跟 UseCrawlerProxy 互斥，优先级高于 UseCrawlerProxy
	OpenBrowser(ctx *BrowserContext)
	OpenBrowserWithParams(ctx *BrowserContext, json *gjson.Json)
	FillTargetRequest(ctx *BrowserContext)
}

type AbstractCrawlByBrowser struct {
	CrawlByBrowserInterface
}

func (r *AbstractCrawlByBrowser) UseBrowser() bool {
	// 默认抓M3U8 URL都用浏览器
	return true
}

func (r *AbstractCrawlByBrowser) UseMobileUA() bool {
	return false
}

func (r *AbstractCrawlByBrowser) UseCrawlerProxy() bool {
	return false
}

func (r *AbstractCrawlByBrowser) UseBrowserMobProxy() bool {
	// 浏览器抓M3U8 URL的时候使用BrowserMobProxy
	return true
}

func (r *AbstractCrawlByBrowser) OpenBrowser(ctx *BrowserContext) {
	// 子类根据策略选择性实现
}

func (r *AbstractCrawlByBrowser) OpenBrowserWithParams(ctx *BrowserContext, json *gjson.Json) {
	// 子类根据策略选择性实现
}

func (r *AbstractCrawlByBrowser) FillTargetRequest(ctx *BrowserContext) {
	// 子类根据策略选择性实现
}
