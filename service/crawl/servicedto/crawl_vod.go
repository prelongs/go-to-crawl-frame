package servicedto

import "github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"

type CrawlVodConfigDO struct {
	*entity.CmsCrawlVodConfig
	*entity.CmsCrawlVodConfigTask
}
