package vodservice

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/crawl/servicedto"
	"go-to-crawl-frame/utils/timeutil"
	"time"
)

var (
	vt  = dao.CmsCrawlVodTv.Columns()
	vti = dao.CmsCrawlVodTvItem.Columns()
	vc  = dao.CmsCrawlVodConfig.Columns()
	vct = dao.CmsCrawlVodConfigTask.Columns()
)

const (
	ConfigTaskStatusInit       = 0
	ConfigTaskStatusProcessing = 1
	ConfigTaskStatusErr        = 2
	ConfigTaskStatusOk         = 3
)

func GetVodConfigById(id int) *entity.CmsCrawlVodConfig {
	var one *entity.CmsCrawlVodConfig
	_ = dao.CmsCrawlVodConfig.Ctx(gctx.GetInitCtx()).Scan(&one, vc.Id, id)
	return one
}

func GetVodConfig() *entity.CmsCrawlVodConfig {
	var one *entity.CmsCrawlVodConfig
	hourBefore := time.Now().Add(-gtime.H).Format(timeutil.YYYY_MM_DD_HH_MM_SS)
	_ = dao.CmsCrawlVodConfig.Ctx(gctx.GetInitCtx()).
		Where(fmt.Sprintf("%v < '%v' or %v is null", vc.UpdateTime, hourBefore, vc.UpdateTime)).
		Where(vc.SeedStatus, 1).
		Order(vc.UpdateTime).Scan(&one)
	return one
}

func UpdateVodConfig(vodConfig *entity.CmsCrawlVodConfig) {
	vodConfig.UpdateTime = gtime.Now()
	dao.CmsCrawlVodConfig.Ctx(gctx.GetInitCtx()).Data(vodConfig).Where(vc.Id, vodConfig.Id).Update()
}

func UpdateVodConfigTaskStatus(configTask *entity.CmsCrawlVodConfigTask, status int) {
	configTask.TaskStatus = status
	configTask.UpdateTime = gtime.Now()
	dao.CmsCrawlVodConfigTask.Ctx(gctx.GetInitCtx()).Data(configTask).Where(vct.Id, configTask.Id).Update()
}

func GetVodConfigTaskDO() *servicedto.CrawlVodConfigDO {
	var configTask *entity.CmsCrawlVodConfigTask
	_ = dao.CmsCrawlVodConfigTask.Ctx(gctx.GetInitCtx()).Scan(&configTask, vct.TaskStatus, ConfigTaskStatusInit)
	if configTask == nil {
		return nil
	}

	var config *entity.CmsCrawlVodConfig
	_ = dao.CmsCrawlVodConfig.Ctx(gctx.GetInitCtx()).Scan(&config, vc.Id, configTask.VodConfigId)

	taskDO := new(servicedto.CrawlVodConfigDO)
	taskDO.CmsCrawlVodConfigTask = configTask
	taskDO.CmsCrawlVodConfig = config

	return taskDO
}

func GetVodTvById(id int) *entity.CmsCrawlVodTv {
	var one *entity.CmsCrawlVodTv
	_ = dao.CmsCrawlVodTv.Ctx(gctx.GetInitCtx()).Where(vt.Id, id).Scan(&one)
	return one
}

func GetVodTvByStatus(crawlStatus int) *entity.CmsCrawlVodTv {
	var one *entity.CmsCrawlVodTv
	_ = dao.CmsCrawlVodTv.Ctx(gctx.GetInitCtx()).Where(vt.CrawlStatus, crawlStatus).Scan(&one)
	return one
}

func GetVodTvByMd5(vodMd5 string) *entity.CmsCrawlVodTv {
	var one *entity.CmsCrawlVodTv
	_ = dao.CmsCrawlVodTv.Ctx(gctx.GetInitCtx()).Where(vt.VodMd5, vodMd5).Scan(&one)
	return one
}

func UpdateVodTVStatus(vodTv *entity.CmsCrawlVodTv, status int) {
	vodTv.CrawlStatus = status
	vodTv.UpdateTime = gtime.Now()
	dao.CmsCrawlVodTv.Ctx(gctx.GetInitCtx()).Data(vodTv).Where(vt.Id, vodTv.Id).Update()
}

func GetPreparedVodTvItem() *entity.CmsCrawlVodTvItem {
	join := g.Model(dao.CmsCrawlVodTvItem.Table()+" vti").LeftJoin(dao.CmsCrawlVodTv.Table()+" vt", fmt.Sprintf("vti.%s = vt.%s", vti.TvId, vt.Id))
	record, _ := join.Fields("vti.*").One(fmt.Sprintf("vti.%s = %d and vt.%s = %d", vti.CrawlStatus, CrawlTVItemInit, vt.CrawlStatus, CrawlTVPadIdOk))
	if record == nil {
		return nil
	}

	tvItem := new(entity.CmsCrawlVodTvItem)
	_ = gconv.Struct(record, tvItem)
	return tvItem
}

func GetVodTvItemByMd5(vodItemMd5 string) *entity.CmsCrawlVodTvItem {
	var one *entity.CmsCrawlVodTvItem
	_ = dao.CmsCrawlVodTvItem.Ctx(gctx.GetInitCtx()).Where(vti.TvItemMd5, vodItemMd5).Scan(&one)
	return one
}

func GetVodTvItemByVideoItemId(videoItemId string) *entity.CmsCrawlVodTvItem {
	var one *entity.CmsCrawlVodTvItem
	_ = dao.CmsCrawlVodTvItem.Ctx(gctx.GetInitCtx()).Where(vti.VideoItemId, videoItemId).Scan(&one)
	return one
}

func UpdateVodTVItemStatus(vodTvItem *entity.CmsCrawlVodTvItem, status int) {
	vodTvItem.CrawlStatus = status
	vodTvItem.UpdateTime = gtime.Now()
	dao.CmsCrawlVodTvItem.Ctx(gctx.GetInitCtx()).Data(vodTvItem).Where(vti.Id, vodTvItem.Id).Update()
}
