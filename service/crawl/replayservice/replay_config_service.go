package replayservice

import (
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
)

var (
	cc = dao.CmsCrawlReplayConfig.Columns()
)

func ListManifestConfig() []*entity.CmsCrawlReplayConfig {
	var all []*entity.CmsCrawlReplayConfig
	_ = dao.CmsCrawlReplayConfig.Ctx(gctx.GetInitCtx()).Where(cc.Status, 1).Scan(&all)
	return all
}

func GetReplayConfig(id int) *entity.CmsCrawlReplayConfig {
	var one *entity.CmsCrawlReplayConfig
	_ = dao.CmsCrawlReplayConfig.Ctx(gctx.GetInitCtx()).Where(cc.Id, id).Scan(&one)
	return one
}
