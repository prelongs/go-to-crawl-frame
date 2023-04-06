package configservice

import (
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/servicedto"
	"github.com/JervisPG/go-to-crawl-frame/utils/constant"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

func GetCrawlHostIp() string {
	return GetCrawlCfg("hostIp")
}

func GetCrawlCfg(key string) string {
	return GetString(fmt.Sprintf("crawl.%s", key))
}

func GetCrawlBool(key string) bool {
	return GetBool(fmt.Sprintf("crawl.%s", key))
}

func GetCrawlDebugBool(key string) bool {
	return GetBool(fmt.Sprintf("crawl.debug.%s", key))
}

func GetCrawlBrowser() *servicedto.CrawlBrowserDO {
	browserDO := new(servicedto.CrawlBrowserDO)
	browserDO.BrowserType = GetCrawlCfg(constant.CrawlBrowserType)
	browserDO.UserDataDir = GetCrawlCfg(constant.CrawlBrowserUserDataDir)
	browserDO.ProxyPath = GetCrawlCfg(constant.CrawlBrowserProxyPath)
	browserDO.Headless = GetCrawlBool(constant.CrawlBrowserHeadless)

	browserInfoList := GetArray("crawl.browser.browserInfoList")
	for _, item := range browserInfoList {
		infoDO := new(servicedto.CrawlBrowserInfoDO)
		infoJson := gjson.New(item)

		infoDO.DriverType = infoJson.Get("driverType").String()
		infoDO.DriverPath = infoJson.Get("driverPath").String()
		infoDO.ExecutorPath = infoJson.Get("executorPath").String()
		infoDO.ProfilePath = infoJson.Get("profilePath").String()

		if constant.DriverTypeChrome == infoDO.DriverType {
			infoDO.UriPrefix = "/wd/hub"
		}

		browserDO.BrowserInfos = append(browserDO.BrowserInfos, infoDO)
	}

	return browserDO
}

func GetCrawlBrowserInfo() *servicedto.CrawlBrowserInfoDO {
	browser := GetCrawlBrowser()
	if browser.BrowserInfos == nil {
		return nil
	}

	for _, info := range browser.BrowserInfos {
		if info.DriverType == browser.BrowserType {
			return info
		}
	}
	return nil
}

func GetCrawl(key string) string {
	return GetString(fmt.Sprintf("server.%s", key))
}

func GetServerCfg(key string) string {
	return GetString(fmt.Sprintf("server.%s", key))
}

func GetString(key string) string {
	value, err := g.Cfg().Get(gctx.GetInitCtx(), key)
	if err != nil {
		return ""
	}
	return value.String()
}

func GetBool(key string) bool {
	value, err := g.Cfg().Get(gctx.GetInitCtx(), key)
	if err != nil {
		return false
	}
	return value.Bool()
}

func GetArray(key string) []interface{} {
	arr, err := g.Cfg().Get(gctx.GetInitCtx(), key)
	if err != nil {
		return nil
	}
	return arr.Array()
}

func GetStrings(key string) []string {
	var strs []string

	for _, item := range GetArray(key) {
		strs = append(strs, gconv.String(item))
	}
	return strs
}

func GetInt(key string, defValue int) int {
	value, err := g.Cfg().Get(gctx.GetInitCtx(), key)
	if err != nil {
		return defValue
	}
	return value.Int()
}
