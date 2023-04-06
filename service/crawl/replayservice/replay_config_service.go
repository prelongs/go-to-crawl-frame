package replayservice

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/dao"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/gogf/gf/v2/os/gctx"
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
