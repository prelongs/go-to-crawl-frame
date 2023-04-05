package replay

import (
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-video/task/replaytask/programme"
	"testing"
)

func TestInitManifestTasks(t *testing.T) {
	Manifest.InitManifestTasks(gctx.GetInitCtx())
}

func TestCrawlReplayManifest(t *testing.T) {
	Manifest.CrawlReplayManifest(gctx.GetInitCtx())
}

func TestCrawlProgrammeManifest(t *testing.T) {
	doCrawlReplayManifest("https://programme.tvb.com/jade/week/")
}

func TestCrawlTbcManifest(t *testing.T) {
	doCrawlReplayManifest("https://www.tbc.net.tw/Epg/Channel?channelId=068")
}

func TestCrawlSports8Manifest(t *testing.T) {
	doCrawlReplayManifest("https://sports8.net/program/70/2.htm")
}

func TestCrawlTvkingdomManifest(t *testing.T) {
	doCrawlReplayManifest("https://www.tvkingdom.jp/chart/23.action")
}

func TestCrawlViuTvManifest(t *testing.T) {
	doCrawlReplayManifest("https://api.viu.tv/production/epgs/96")
}

func TestCrawlProgramTableManifest(t *testing.T) {
	doCrawlReplayManifest("https://節目表.tw/channel/台視/#")
}

func TestStartRecording(t *testing.T) {
	programTask := new(entity.CmsCrawlReplayProgramTask)
	programTask.Id = 227

	ctx := new(CrawlReplayProgram)
	ctx.programTask = programTask

	ctx.StartRecording(gctx.GetInitCtx())
}

func TestRequestProgramme(t *testing.T) {
	config := new(entity.CmsCrawlReplayConfig)
	config.SeedUrl = "https://programme.tvb.com/jade"
	config.Id = 2
	programmeCrawl := new(programme.ProgrammeCrawl)
	m := &entity.CmsCrawlReplayManifestTask{
		Id: 1,
	}
	programmeCrawl.CreateProgram(config, m)
}

func doCrawlReplayManifest(seedUrl string) {
	config := new(entity.CmsCrawlReplayConfig)
	config.SeedUrl = seedUrl
	manifest := new(entity.CmsCrawlReplayManifestTask)
	manifest.Id = 0

	Manifest.doCrawlReplayManifest(config, manifest)
}
