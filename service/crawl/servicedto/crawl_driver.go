package servicedto

import (
	"github.com/tebeka/selenium"
	"go-to-crawl-frame/service/browsermobservice"
)

type CrawlWebDriverDO struct {
	Service *selenium.Service
	XServer *browsermobservice.Server
	XClient *browsermobservice.Client
	Wd      selenium.WebDriver
	*CrawlBrowserDO
	*CrawlBrowserInfoDO
}
