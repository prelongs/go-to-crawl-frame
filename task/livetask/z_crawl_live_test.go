package livetask

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

func TestDoStartCrawlLiveFlow(t *testing.T) {

	config := new(entity.CmsCrawlLiveConfig)
	config.ProgramName = "TNT"
	config.LiveUrl = "https://live.tv247us.com/tv247/tnt.m3u8"

	CrawlTask.CmsCrawlLiveConfig = config
	CrawlTask.DoStartCrawlLiveFlow(gctx.GetInitCtx())

}
