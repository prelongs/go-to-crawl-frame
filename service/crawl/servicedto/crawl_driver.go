package servicedto

import (
	"github.com/JervisPG/go-to-crawl-frame/service/browsermobservice"
	"github.com/tebeka/selenium"
)

type CrawlWebDriverDO struct {
	Service *selenium.Service
	XServer *browsermobservice.Server
	XClient *browsermobservice.Client
	Wd      selenium.WebDriver
	*CrawlBrowserDO
	*CrawlBrowserInfoDO
}
