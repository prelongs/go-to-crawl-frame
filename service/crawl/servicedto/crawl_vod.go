package servicedto

import "go-to-crawl-frame/db/mysql/model/entity"

type CrawlVodConfigDO struct {
	*entity.CmsCrawlVodConfig
	*entity.CmsCrawlVodConfigTask
}
