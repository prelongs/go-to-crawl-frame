package crawltask

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/tebeka/selenium"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/browsermobservice"
	"go-to-crawl-frame/service/crawl/servicedto"
	"go-to-crawl-frame/service/crawl/sysservice"
	"go-to-crawl-frame/service/crawl/vodservice"
	"go-to-crawl-frame/service/lockservice"
	"go-to-crawl-frame/task/taskdto"
	"go-to-crawl-frame/utils/browsermob"
	"go-to-crawl-frame/utils/browserutil"
)

var CrawlVodTVTask = new(crawlVodTVTask)

type crawlVodTVTask struct {
}

func (crawlUrl *crawlVodTVTask) GenVodConfigTask(context context.Context) {
	vodConfig := vodservice.GetVodConfig()
	if vodConfig == nil {
		return
	}
	vodservice.UpdateVodConfig(vodConfig)

	configTask := new(entity.CmsCrawlVodConfigTask)
	configTask.VodConfigId = vodConfig.Id
	configTask.TaskStatus = vodservice.ConfigTaskStatusInit
	configTask.CreateTime = gtime.Now()

	dao.CmsCrawlVodConfigTask.Ctx(gctx.GetInitCtx()).Insert(configTask)

}

func (crawlUrl *crawlVodTVTask) VodTVTask(context context.Context) {
	locked := lockservice.TryLockSeleniumLong()

	if !locked {
		return
	}
	defer lockservice.ReleaseLockSelenium()
	vodConfigTaskDO := vodservice.GetVodConfigTaskDO()
	if vodConfigTaskDO != nil {
		vodservice.UpdateVodConfigTaskStatus(vodConfigTaskDO.CmsCrawlVodConfigTask, vodservice.ConfigTaskStatusProcessing)
		DoStartCrawlVodTV(vodConfigTaskDO)
		vodservice.UpdateVodConfigTaskStatus(vodConfigTaskDO.CmsCrawlVodConfigTask, vodservice.ConfigTaskStatusOk)
	}
}

// 填充视频基础信息
func (crawlUrl *crawlVodTVTask) VodTVPadInfoTask(context context.Context) {
	log := g.Log().Line()
	locked := lockservice.TryLockSelenium()
	if !locked {
		return
	}
	defer lockservice.ReleaseLockSelenium()

	vodTv := vodservice.GetVodTvByStatus(vodservice.CrawlTVInit)
	if vodTv == nil {
		return
	}

	log.Infof(gctx.GetInitCtx(),
		"更新vod tv. id = %v, to status = %v", vodTv.Id, vodservice.CrawlTVPadInfo)
	vodservice.UpdateVodTVStatus(vodTv, vodservice.CrawlTVPadInfo)

	DoStartCrawlVodPadInfo(vodTv)
}

func DoStartCrawlVodTV(configTaskDO *servicedto.CrawlVodConfigDO) {
	ctx := taskdto.GetInitedCrawlContext()
	ctx.CrawlVodConfigDO = configTaskDO
	strategy := getCrawlVodTVStrategy(ctx.CrawlVodConfigDO.CmsCrawlVodConfig)
	ctx.CrawlByBrowserInterface = strategy

	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := browserutil.GetDriverService(browserutil.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()

		proxyUrl := ""
		if strategy.UseCrawlerProxy() {
			proxyUrl = sysservice.GetRandomProxyUrl()
			ctx.Log.Infof(gctx.GetInitCtx(), "visit list page via proxy. domain = %v, proxy = %v", configTaskDO.DomainKeyPart, proxyUrl)
		}

		var caps selenium.Capabilities

		if strategy.UseBrowserMobProxy() {
			xServer := browsermobservice.NewServer(ctx.ProxyPath)
			xServer.Start()
			ctx.XServer = xServer
			defer ctx.XServer.Stop()
			proxy := xServer.CreateProxy(nil)
			ctx.XClient = proxy
			defer ctx.XClient.Close()
			caps = browserutil.GetAllCaps(ctx)
		} else {
			caps = browserutil.GetAllCapsChooseProxy(ctx, proxyUrl)
		}

		webDriver, err := browserutil.NewRemote(caps, browserutil.DriverServicePort, ctx.UriPrefix)
		ctx.Wd = webDriver
		if ctx.Wd == nil {
			ctx.CrawlVodConfigDO.CmsCrawlVodConfig.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if ctx.CrawlVodConfigDO.CmsCrawlVodConfig.SeedParams != "" {
			json, _ := gjson.LoadJson(ctx.CrawlVodConfigDO.CmsCrawlVodConfig.SeedParams)
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
}

func DoStartCrawlVodPadInfo(vodTVItem *entity.CmsCrawlVodTv) {
	ctx := taskdto.GetInitedCrawlContext()
	ctx.VodTV = vodTVItem
	strategy := getCrawlVodPadInfoStrategy(ctx.VodTV)
	ctx.CrawlByBrowserInterface = strategy

	if strategy.UseBrowser() {
		//g.Dump("使用浏览器")
		service, _ := browserutil.GetDriverService(browserutil.DriverServicePort)
		ctx.Service = service
		defer ctx.Service.Stop()

		proxyUrl := ""
		if strategy.UseCrawlerProxy() {
			proxyUrl = sysservice.GetRandomProxyUrl()
			ctx.Log.Infof(gctx.GetInitCtx(),
				"visit detail page via proxy. id = %v, proxy = %v", vodTVItem.Id, proxyUrl)
		}

		var caps selenium.Capabilities

		if strategy.UseBrowserMobProxy() {
			xServer := browsermobservice.NewServer(ctx.ProxyPath)
			xServer.Start()
			ctx.XServer = xServer
			defer ctx.XServer.Stop()
			proxy := xServer.CreateProxy(nil)
			ctx.XClient = proxy
			defer ctx.XClient.Close()

			// BrowserMobProxy抓包方式
			caps = browserutil.GetAllCaps(ctx)
		} else {
			caps = browserutil.GetAllCapsChooseProxy(ctx, proxyUrl)
		}

		webDriver, err := browserutil.NewRemote(caps, browserutil.DriverServicePort, ctx.UriPrefix)
		ctx.Wd = webDriver
		if ctx.Wd == nil {
			ctx.VodTVItem.ErrorMsg = "WebDriver Init Fail"
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}
		defer ctx.Wd.Quit()

		// 业务处理-start
		if strategy.UseBrowserMobProxy() {
			browsermob.NewHarWait(ctx.Wd, ctx.XClient)
		}
		strategy.OpenBrowser(ctx)
		// 业务处理-end
	}

	// 把URL,Headers信息保存起来
	strategy.FillTargetRequest(ctx)
}
