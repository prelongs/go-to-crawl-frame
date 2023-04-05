package servicedto

type CrawlBrowserDO struct {
	BrowserType string
	UserDataDir string
	ProxyPath   string
	Headless    bool

	BrowserInfos []*CrawlBrowserInfoDO
}

type CrawlBrowserInfoDO struct {
	DriverType   string
	DriverPath   string
	ExecutorPath string
	ProfilePath  string
	UriPrefix    string
}
