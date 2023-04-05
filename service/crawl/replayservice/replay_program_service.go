package replayservice

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/configservice"
	"go-to-crawl-frame/service/crawl/servicedto"
)

var (
	pc = dao.CmsCrawlReplayProgramTask.Columns()
)

func GetProgramTask(programTask *entity.CmsCrawlReplayProgramTask) *entity.CmsCrawlReplayProgramTask {
	var one *entity.CmsCrawlReplayProgramTask
	_ = dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).
		Where(pc.ConfigId, programTask.ConfigId).
		Where(pc.ProgramStartTime, programTask.ProgramStartTime).Scan(&one)
	return one
}

func GetById(id int) *entity.CmsCrawlReplayProgramTask {
	var one *entity.CmsCrawlReplayProgramTask
	_ = dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Where(pc.Id, id).Scan(&one)
	return one
}

func GetDetailById(id int) *servicedto.CmsCrawlReplayProgramTaskDO {
	join := g.DB().Ctx(gctx.GetInitCtx()).Model(dao.CmsCrawlReplayProgramTask.Table()+" rpt").LeftJoin(dao.CmsCrawlReplayManifestTask.Table()+" rmt", "rmt.id = rpt.manifest_id")
	join = join.LeftJoin(dao.CmsCrawlReplayConfig.Table()+" rc", " rc.id = rmt.replay_config_id")
	join = join.Where("rpt.id = ", id)
	taskDO := new(servicedto.CmsCrawlReplayProgramTaskDO)
	_ = join.Fields("rpt.*,rmt.replay_day,rc.channel_name,rc.channel_no,rc.play_url").Scan(&taskDO)
	return taskDO
}

func GetAllByStatus(status int) []*entity.CmsCrawlReplayProgramTask {
	var list []*entity.CmsCrawlReplayProgramTask
	_ = dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).
		Where(pc.CrawlStatus, status).
		// 只把当前时间还没到的记录加入到录制任务
		Where(pc.HostIp, configservice.GetCrawlHostIp()).Where("program_start_time > ", gtime.Now()).
		Scan(list)
	return list
}

func GetOneTest(status int) *entity.CmsCrawlReplayProgramTask {
	var one *entity.CmsCrawlReplayProgramTask
	// 只把当前时间还没到的记录加入到录制任务
	_ = dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Scan(&one)
	return one
}

// eg: 11:15重启，但有一个节目是11:00-11:30的，需要把这个节目的尾段录下来
func GetAllCurrentByStatus(status int) []*entity.CmsCrawlReplayProgramTask {
	var list []*entity.CmsCrawlReplayProgramTask
	now := gtime.Now()
	_ = dao.CmsCrawlReplayProgramTask.
		Ctx(gctx.GetInitCtx()).
		Where(pc.CrawlStatus, status).
		Where("program_start_time <= ", now).
		Where("program_end_time > ", now).Scan(&list)
	return list
}

func GetOneByIpStatus(hostIp string, status int) *entity.CmsCrawlReplayProgramTask {
	var one *entity.CmsCrawlReplayProgramTask
	_ = dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).
		Where(pc.CrawlStatus, status).
		Where(pc.HostIp, hostIp).Scan(&one)
	return one
}

func GetAllByIpStatus(hostIp string, status int) []*entity.CmsCrawlReplayProgramTask {
	var list []*entity.CmsCrawlReplayProgramTask
	where := dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Where(pc.CrawlStatus, status)
	where = where.Where(pc.HostIp, hostIp)
	_ = where.Scan(&list)
	return list
}

func GetOneByStatus(status int) (one *entity.CmsCrawlReplayProgramTask, err error) {
	var task *entity.CmsCrawlReplayProgramTask
	e := dao.CmsCrawlReplayProgramTask.
		Ctx(gctx.GetInitCtx()).
		Where(pc.CrawlStatus, status).
		Where(pc.HostIp, configservice.GetCrawlHostIp()).
		Order("id desc").Scan(&task)
	return task, e
}

//是否需要转码 根据配置文件最大转码数量判断
func NeedToTran(status int) bool {
	where := dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Where(pc.CrawlStatus, status)
	where = where.Where(pc.HostIp, configservice.GetCrawlHostIp()).Order("id")
	count, err := where.Count()
	if err != nil {
		return false
	}
	conf_count := configservice.GetInt("crawl.maxTrans", 5)
	//g.Dump(count,conf_count)
	if count >= conf_count {
		return false
	}
	return true
}

func UpdateProgramTaskStatus(programTask *entity.CmsCrawlReplayProgramTask, crawlStatus int) {
	programTask.CrawlStatus = crawlStatus
	dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Data(programTask).Where(pc.Id, programTask.Id).Update()
}

func CheckAndSave(programTask *entity.CmsCrawlReplayProgramTask) {
	dbProgramTask := GetProgramTask(programTask)
	if dbProgramTask == nil {
		//不存在则加入数据库
		dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Insert(programTask)
	}
}
