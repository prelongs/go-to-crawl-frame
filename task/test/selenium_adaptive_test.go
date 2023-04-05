package test

import (
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-frame/service/browsermobservice"
	"go-to-crawl-frame/task/taskdto"
	"go-to-crawl-frame/utils/browsermob"
	"go-to-crawl-frame/utils/browserutil"
	"go-to-crawl-frame/utils/processutil"
	"testing"
)

func TestStart(t *testing.T) {
	Start(false)
}

// 测试web容器环境下多次打开和关闭浏览器代理是否会暂停应用
func Start(useBrowserMob bool) {

	pid, _ := processutil.CheckRunning(browsermobservice.PORT)
	if pid != "" {
		return
	}

	ctx := taskdto.GetInitedCrawlContext()

	service, _ := browserutil.GetDriverService(browserutil.DriverServicePort)
	ctx.Service = service
	defer ctx.Service.Stop()

	if useBrowserMob {
		xServer := browsermobservice.NewServer(ctx.ProxyPath)
		xServer.Start()
		ctx.XServer = xServer
		defer ctx.XServer.Stop()
		proxy := xServer.CreateProxy(nil)
		ctx.XClient = proxy
		defer ctx.XClient.Close()

		// BrowserMobProxy抓包方式
	}

	caps := browserutil.GetAllCaps(ctx)

	webDriver, err := browserutil.NewRemote(caps, browserutil.DriverServicePort, ctx.UriPrefix)
	ctx.Wd = webDriver
	if ctx.Wd == nil {
		ctx.Log.Error(gctx.GetInitCtx(), err)
		ctx.CrawlQueueSeed.ErrorMsg = "WebDriver Init Fail"
		return
	}
	defer ctx.Wd.Quit()

	if useBrowserMob {
		browsermob.NewHarWait(ctx.Wd, ctx.XClient)
	}
	_ = ctx.Wd.Get("https://www.nivod3.tv/1AyYWd1WFag2bKjJliUuuAQFT2vgDKzB-0-0-0-0-play.html?x=1")
	if useBrowserMob {
		browsermob.GetHarRequest(ctx.XClient, ".m3u8", "", 1)
	}
}
