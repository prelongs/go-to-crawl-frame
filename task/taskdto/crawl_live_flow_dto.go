package taskdto

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/utils/ffmpegutil"
	"github.com/JervisPG/go-to-crawl-frame/utils/httputil"
	"github.com/gogf/gf/v2/encoding/gjson"
)

// 抓取直播接口集合
type CrawlLiveFlowInterface interface {
	CrawlByBrowserInterface

	LoadLiveStream(ctx *BrowserContext)
	// 下载视频接口集合
	ConvertM3U8(seed *entity.CmsCrawlLiveConfig, filePath string) (*ffmpegutil.M3u8DO, error)
	ConvertM3U8GetBaseUrl(m3u8Url string) string
}

type AbstractCrawlLiveFlow struct {
	CrawlByBrowserInterface
	*AbstractCrawlByBrowser
}

func (r *AbstractCrawlLiveFlow) UseBrowser() bool {
	return false
}

func (r *AbstractCrawlLiveFlow) UseMobileUA() bool {
	return false
}

func (r *AbstractCrawlLiveFlow) UseCrawlerProxy() bool {
	return false
}

func (r *AbstractCrawlLiveFlow) UseBrowserMobProxy() bool {
	return false
}

func (r *AbstractCrawlLiveFlow) OpenBrowser(ctx *BrowserContext) {
}

func (r *AbstractCrawlLiveFlow) OpenBrowserWithParams(ctx *BrowserContext, json *gjson.Json) {
}

func (r *AbstractCrawlLiveFlow) FillTargetRequest(ctx *BrowserContext) {
}

func (r *AbstractCrawlLiveFlow) LoadLiveStream(ctx *BrowserContext) {
}

func (r *AbstractCrawlLiveFlow) ConvertM3U8(liveConfig *entity.CmsCrawlLiveConfig, filePath string) (*ffmpegutil.M3u8DO, error) {
	baseUrl := r.ConvertM3U8GetBaseUrl(liveConfig.LiveUrl)
	return ffmpegutil.ConvertM3U8WithOriginTsName(liveConfig.LiveUrl, baseUrl, filePath)
}

func (r *AbstractCrawlLiveFlow) ConvertM3U8GetBaseUrl(m3u8Url string) string {
	return httputil.GetBaseUrlBySchema(m3u8Url)
}
