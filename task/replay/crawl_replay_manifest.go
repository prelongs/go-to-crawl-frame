package replay

import (
	"context"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/replayservice"
)

var Manifest = new(CrawlReplayManifest)

type CrawlReplayManifest struct {
	programTask *entity.CmsCrawlReplayProgramTask
}

func (receiver *CrawlReplayManifest) InitManifestTasks(context context.Context) {
	replayservice.InitManifestTasks()
}

func (receiver *CrawlReplayManifest) CrawlReplayManifest(context context.Context) {

	task := replayservice.GetManifestTask(replayservice.ManifestTaskInit)
	if task == nil {
		return
	}

	replayservice.UpdateManifestTaskStatus(task, replayservice.ManifestTaskCrawling)
	replayConfig := replayservice.GetReplayConfig(task.ReplayConfigId)
	receiver.doCrawlReplayManifest(replayConfig, task)
	replayservice.UpdateManifestTaskStatus(task, replayservice.ManifestTaskCrawlFinish)

}

func (receiver *CrawlReplayManifest) doCrawlReplayManifest(replayConfig *entity.CmsCrawlReplayConfig, task *entity.CmsCrawlReplayManifestTask) {
	strategy := getCrawlReplayStrategy(replayConfig)
	strategy.CreateProgram(replayConfig, task)
}
