package replayservice

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/dao"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/utils/timeutil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
)

var (
	mc = dao.CmsCrawlReplayManifestTask.Columns()
)

func InitManifestTasks() {
	today := time.Now().Format(timeutil.YYYY_MM_DD)
	manifestList := ListManifestConfig()
	//g.Dump(len(manifestList))
	for _, manifest := range manifestList {

		CreateTask(manifest, today)
	}

}

func CreateTask(config *entity.CmsCrawlReplayConfig, day string) {
	var manifestTask *entity.CmsCrawlReplayManifestTask
	_ = dao.CmsCrawlReplayManifestTask.Ctx(gctx.GetInitCtx()).
		Where(mc.ReplayDay, day).
		Where(mc.ReplayConfigId, config.Id).
		Where(g.Map{
			"DATE_ADD(create_time,INTERVAL ? HOUR)>now()": 6,
		}).
		Scan(&manifestTask)
	if manifestTask == nil {
		task := new(entity.CmsCrawlReplayManifestTask)
		task.CrawlStatus = ManifestTaskInit
		task.ReplayConfigId = config.Id
		task.ReplayDay = day
		task.CreateTime = gtime.Now()
		task.Type = config.Type
		dao.CmsCrawlReplayManifestTask.Ctx(gctx.GetInitCtx()).Insert(task)
	}
}

func GetManifestTask(crawlStatus int) *entity.CmsCrawlReplayManifestTask {
	var manifestTask *entity.CmsCrawlReplayManifestTask
	_ = dao.CmsCrawlReplayManifestTask.Ctx(gctx.GetInitCtx()).Where(mc.CrawlStatus, crawlStatus).Scan(&manifestTask)
	return manifestTask
}

func UpdateManifestTaskStatus(manifestTask *entity.CmsCrawlReplayManifestTask, crawlStatus int) {
	manifestTask.CrawlStatus = crawlStatus
	dao.CmsCrawlReplayManifestTask.Ctx(gctx.GetInitCtx()).Data(manifestTask).Where(mc.Id, manifestTask.Id).Update()
}
