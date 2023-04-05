package constant

const (
	FlowRetry      = 1 // 本地整体流程重试次数，适用于更大范围的重试
	LocalRetry     = 1 // 本地mob重试次数，适用于M3U8出现的速度较慢的场景，循环等待M3U8出现
	LocalRefresh   = 5 // 本地刷新次数，适用于浏览器遇到了概率防爬，本地多刷新几次就可以的场景
	ServerMaxRetry = 3 // 服务端重试次数，适用于本地加载不了，通过定时任务下一次再来抓取的场景
)

var (
	CrawlBrowserType        = "browser.browserType"
	CrawlBrowserHeadless    = "browser.headless"
	CrawlBrowserUserDataDir = "browser.userDataDir"
	CrawlBrowserProxyPath   = "browser.proxyPath"
)

var (
	DriverTypeChrome  = "chrome"
	DriverTypeFireFox = "firefox"
	DriverTypeEdge    = "edge"
)
