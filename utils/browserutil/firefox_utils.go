package browserutil

import (
	"github.com/JervisPG/go-to-crawl-frame/service/configservice"
	"github.com/JervisPG/go-to-crawl-frame/task/taskdto"
	"github.com/tebeka/selenium/firefox"
)

func GetFirefoxCaps(browserCtx *taskdto.BrowserContext, crawlerProxy string) firefox.Capabilities {
	args := []string{
		"--no-sandbox",
		//"--disable-blink-features=AutomationControlled", // 隐藏自己是selenium. window.navigator.webdrive=true
		"--user-agent=" + getRandomUA(browserCtx),
		//"--acceptSslCerts=true",
	}

	args = appendProxyArgs(args, browserCtx, crawlerProxy)
	args = appendConfigArgs(args)

	browserInfoConfig := configservice.GetCrawlBrowserInfo()
	specCaps := firefox.Capabilities{
		Binary: browserInfoConfig.ExecutorPath,
		Args:   args,
		Prefs: map[string]interface{}{
			//"profile.managed_default_content_settings.images": 2,
			//"permissions.default.stylesheet": 2,
			"dom.webdriver.enabled": false, // 去除window.navigator.webdriver属性的核心语句
		},
	}

	if browserInfoConfig.ProfilePath != "" {
		specCaps.Args = append(specCaps.Args, "--profile", browserInfoConfig.ProfilePath)
	}

	return specCaps
}
