package browserutil

import (
	"github.com/JervisPG/go-to-crawl-frame/service/configservice"
	"github.com/JervisPG/go-to-crawl-frame/task/taskdto"
	"github.com/tebeka/selenium/chrome"
)

const (
	DriverServicePort = 8088
	ScrollBottomJs    = "window.scrollTo(0,document.body.scrollHeight);"
)

func GetChromeCaps(browserCtx *taskdto.BrowserContext, crawlerProxy string) chrome.Capabilities {
	args := []string{
		"--no-sandbox",
		"--ignore-certificate-errors",
		"--disable-blink-features=AutomationControlled", // 隐藏自己是selenium. window.navigator.webdrive=true
		"--user-agent=" + getRandomUA(browserCtx),
		"--acceptSslCerts=true",
	}

	args = appendProxyArgs(args, browserCtx, crawlerProxy)
	args = appendConfigArgs(args)

	specCaps := chrome.Capabilities{
		Path:  configservice.GetCrawlBrowserInfo().ExecutorPath,
		Args:  args,
		Prefs: map[string]interface{}{
			//"profile.managed_default_content_settings.images": 2,
			//"permissions.default.stylesheet": 2,
		},
		ExcludeSwitches: []string{"enable-automation"},
	}
	return specCaps
}

func getRandomUA(browserCtx *taskdto.BrowserContext) string {
	if browserCtx == nil || !browserCtx.UseMobileUA() {
		return GetRandomUA(false)
	}

	return GetRandomUAByPlatform(Android)
}
