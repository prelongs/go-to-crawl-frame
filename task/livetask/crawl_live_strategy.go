package livetask

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/task/taskdto"
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-video/task/livetask/tv247"
)

const (
	TV247 = "tv247us.com" // 站点域名为tv247.us, m3u8域名tv247us.com
)

func getCrawlLiveFlowStrategy(liveConfig *entity.CmsCrawlLiveConfig) taskdto.CrawlLiveFlowInterface {

	url := liveConfig.LiveUrl

	if gstr.Contains(url, TV247) {
		return new(tv247.TV247Crawl)
	}

	return nil
}
