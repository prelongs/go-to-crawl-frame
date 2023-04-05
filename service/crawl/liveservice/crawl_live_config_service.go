package liveservice

import (
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/configservice"
)

const (
	StatusEnable  = 1
	StatusDisable = 2
)

var (
	lc = dao.CmsCrawlLiveConfig.Columns()
)

// 根据yaml配置的分组, 获取直播配置列表
func GetLiveConfigList() []*entity.CmsCrawlLiveConfig {
	grp := configservice.GetServerCfg("openLiveTaskGroup")
	if grp == "" {
		return nil
	}

	var all []*entity.CmsCrawlLiveConfig
	_ = dao.CmsCrawlLiveConfig.Ctx(gctx.GetInitCtx()).
		Where(lc.GroupName, grp).
		Where(lc.Status, StatusEnable).Scan(&all)
	return all
}
