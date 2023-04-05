package crawltask

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/browsermobservice"
	"go-to-crawl-frame/service/crawl/vodservice"
	"go-to-crawl-frame/service/lockservice"
	"go-to-crawl-frame/task/taskdto"
	"go-to-crawl-frame/utils/browsermob"
	"go-to-crawl-frame/utils/browserutil"
)

var CrawlTask = new(CrawlUrl)

type CrawlUrl struct {
}

func (crawlUrl *CrawlUrl) CrawlUrlTask(context context.Context) {
	if !lockservice.TryLockSelenium() {
		return
	}
	defer lockservice.ReleaseLockSelenium()

	seed := getEnvPreparedSeed("", vodservice.HostTypeNormal)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}

}

func (crawlUrl *CrawlUrl) CrawlUrlType1Task(context context.Context) {
	if !lockservice.TryLockSelenium() {
		return
	}
	defer lockservice.ReleaseLockSelenium()

	seed := getEnvPreparedSeed("", vodservice.HostTypeCrawlLogin)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}
}
func (crawlUrl *CrawlUrl) CrawlUrlType2Task(context context.Context) {
	if !lockservice.TryLockSelenium() {
		return
	}
	defer lockservice.ReleaseLockSelenium()

	seed := getEnvPreparedSeed("", vodservice.HostTypeNiVod)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}
}

func (crawlUrl *CrawlUrl) CrawlUrlType3Task(context context.Context) {
	seed := getEnvPreparedSeed("", vodservice.HostTypeBananTV)
	if seed != nil {
		DoStartCrawlVodFlow(seed)
	}
}

func getEnvPreparedSeed(hostname string, hostType int) *entity.CmsCrawlQueue {
	seed := vodservice.GetSeed(vodservice.Init, hostname, hostType)
	if seed == nil {
		return nil
	}
	vodservice.UpdateStatus(seed, vodservice.Crawling)

	return seed
}

func DoStartCrawlVodFlow(seed *entity.CmsCrawlQueue) {
	ctx := taskdto.GetInitedCrawlContext()
	ctx.CrawlQueueSeed = seed
	strategy := getCrawlVodFlowStrategy(ctx.CrawlQueueSeed)
	ctx.CrawlByBrowserInterface = strategy

	DoStartCrawlVodFlowWithCtx(ctx)
}

// 复用已有Ctx. 支撑整体重试场景
func DoStartCrawlVodFlowWithCtx(ctx *taskdto.BrowserContext) {

	strategy := ctx.CrawlByBrowserInterface
	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := browserutil.GetDriverService(browserutil.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()

		if strategy.UseBrowserMobProxy() {
			xServer := browsermobservice.NewServer(ctx.ProxyPath)
			xServer.Start()
			ctx.XServer = xServer
			defer ctx.XServer.Stop()
			proxy := xServer.CreateProxy(nil)
			ctx.XClient = proxy
			defer ctx.XClient.Close()
		}

		caps := browserutil.GetAllCaps(ctx)

		webDriver, err := browserutil.NewRemote(caps, browserutil.DriverServicePort, ctx.UriPrefix)
		ctx.Wd = webDriver
		if ctx.Wd == nil {
			ctx.CrawlQueueSeed.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			vodservice.UpdateStatus(ctx.CrawlQueueSeed, vodservice.CrawlErr)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if ctx.CrawlQueueSeed.CrawlSeedParams != "" && ctx.CrawlQueueSeed.CrawlSeedParams != `{"videoitem":""}` {
			json, _ := gjson.LoadJson(ctx.CrawlQueueSeed.CrawlSeedParams)
			strategy.OpenBrowserWithParams(ctx, json)
		} else {
			if strategy.UseBrowserMobProxy() {
				browsermob.NewHarWait(ctx.Wd, ctx.XClient)
			}
			//g.Dump("打开浏览器")
			strategy.OpenBrowser(ctx)
		}
		// 业务处理-end
	}

	// 把URL,Headers信息保存起来
	strategy.FillTargetRequest(ctx)
	if ctx.CrawlQueueSeed.CrawlM3U8Url == "" && ctx.RetryFlow > 0 {
		ctx.RetryFlow -= 1
		DoStartCrawlVodFlowWithCtx(ctx)
	}
	vodservice.UpdateUrlAndStatus(ctx.CrawlQueueSeed)
}
