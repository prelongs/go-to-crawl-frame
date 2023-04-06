package taskdto

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
)

type CrawlReplayInterface interface {
	CreateProgram(replayConfig *entity.CmsCrawlReplayConfig, manifestTask *entity.CmsCrawlReplayManifestTask)
}

type AbstractCrawlReplayUrl struct {
	CrawlReplayInterface
}

func (receiver *AbstractCrawlReplayUrl) CreateProgram(replayConfig *entity.CmsCrawlReplayConfig, manifestTask *entity.CmsCrawlReplayManifestTask) {

}
