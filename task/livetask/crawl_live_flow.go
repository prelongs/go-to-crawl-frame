package livetask

import (
	"context"
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/liveservice"
	"github.com/JervisPG/go-to-crawl-frame/task/taskdto"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

var CrawlTask = new(CrawlLive)

type CrawlLive struct {
	*entity.CmsCrawlLiveConfig
}

func (r *CrawlLive) Execute(context context.Context) {

	for _, liveConfig := range liveservice.GetLiveConfigList() {
		r.CmsCrawlLiveConfig = liveConfig
		r.DoStartCrawlLiveFlow(context)
	}

}

func (r *CrawlLive) DoStartCrawlLiveFlow(context context.Context) {
	liveConfig := r.CmsCrawlLiveConfig
	ctx := taskdto.GetInitedCrawlContext()
	ctx.CrawlLiveDO.CmsCrawlLiveConfig = liveConfig
	strategy := getCrawlLiveFlowStrategy(liveConfig)
	ctx.CrawlByBrowserInterface = strategy

	strategy.LoadLiveStream(ctx)

	nextLoadTime := ctx.CrawlLiveDO.NextLoadTime
	if nextLoadTime != nil {

		if gtime.Now().After(nextLoadTime) {
			// 兼容下载太久导致定时任务直接关闭的场景
			nextLoadTime = gtime.Now()
		}

		pattern := fmt.Sprintf("%d %d %d * * *", nextLoadTime.Second(), nextLoadTime.Minute(), nextLoadTime.Hour())
		taskName := fmt.Sprintf("CrawlLive-%s-%v", liveConfig.ProgramName, nextLoadTime.TimestampMilliStr())
		_, _ = gcron.AddOnce(gctx.GetInitCtx(), pattern, r.DoStartCrawlLiveFlow, taskName)
		gcron.Start(taskName)
	}
}
