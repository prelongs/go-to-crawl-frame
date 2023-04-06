package taskdto

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/configservice"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/servicedto"
	"github.com/JervisPG/go-to-crawl-frame/utils/ffmpegutil"
	"github.com/JervisPG/go-to-crawl-frame/utils/httputil"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

type BrowserContext struct {
	// 基础部分Context
	Log *glog.Logger
	servicedto.CrawlWebDriverDO

	// 业务部分Context
	CrawlQueueSeed   *entity.CmsCrawlQueue
	CrawlVodConfigDO *servicedto.CrawlVodConfigDO
	CrawlLiveDO      *servicedto.CrawlLiveDO
	VodTV            *entity.CmsCrawlVodTv
	VodTVItem        *entity.CmsCrawlVodTvItem

	// 绑定策略
	CrawlStrategyDO
}

type CrawlStrategyDO struct {
	RetryFlow int // 当前策略是否需要整体流程重试
	CrawlByBrowserInterface
}

func GetInitedCrawlContext() *BrowserContext {
	ctx := new(BrowserContext)
	ctx.Log = g.Log().Line()

	ctx.CrawlBrowserDO = configservice.GetCrawlBrowser()
	ctx.CrawlBrowserInfoDO = configservice.GetCrawlBrowserInfo()

	return ctx
}

// 抓取点播接口集合
type CrawlVodFlowInterface interface {
	CrawlByBrowserInterface

	// 下载视频接口集合
	ConvertM3U8(seed *entity.CmsCrawlQueue, filePath string) (*ffmpegutil.M3u8DO, error)
	ConvertM3U8GetBaseUrl(m3u8Url string) string
	DownLoadToMp4(m3u8DO *ffmpegutil.M3u8DO) error
}

type AbstractCrawlVodFlow struct {
	CrawlByBrowserInterface
	*AbstractCrawlByBrowser
}

func (r *AbstractCrawlVodFlow) UseBrowser() bool {
	return true
}

func (r *AbstractCrawlVodFlow) UseMobileUA() bool {
	return false
}

func (r *AbstractCrawlVodFlow) UseCrawlerProxy() bool {
	return false
}

func (r *AbstractCrawlVodFlow) UseBrowserMobProxy() bool {
	return true
}

func (r *AbstractCrawlVodFlow) OpenBrowser(ctx *BrowserContext) {
}

func (r *AbstractCrawlVodFlow) OpenBrowserWithParams(ctx *BrowserContext, json *gjson.Json) {
}

func (r *AbstractCrawlVodFlow) FillTargetRequest(ctx *BrowserContext) {
}

func (r *AbstractCrawlVodFlow) ConvertM3U8(seed *entity.CmsCrawlQueue, filePath string) (*ffmpegutil.M3u8DO, error) {
	baseUrl := r.ConvertM3U8GetBaseUrl(seed.CrawlM3U8Url)
	return ffmpegutil.ConvertM3U8(seed.CrawlM3U8Url, baseUrl, filePath)
}

func (r *AbstractCrawlVodFlow) ConvertM3U8GetBaseUrl(m3u8Url string) string {
	return httputil.GetBaseUrlBySchema(m3u8Url)
}

func (r *AbstractCrawlVodFlow) DownLoadToMp4(m3u8DO *ffmpegutil.M3u8DO) error {
	return ffmpegutil.DownloadToMp4(m3u8DO)
}
