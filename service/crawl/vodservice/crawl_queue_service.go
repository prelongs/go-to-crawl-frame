package vodservice

import (
	"context"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/dao"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/configservice"
	"github.com/JervisPG/go-to-crawl-frame/utils/constant"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
)

var (
	c = dao.CmsCrawlQueue.Columns()
)

// hostName: 配置文件里能标识出当前节点唯一就行，eg：prod-1, dev-1, cluster-1等
func GetSeed(status int, hostname string, hostType int) *entity.CmsCrawlQueue {
	where := dao.CmsCrawlQueue.Ctx(gctx.GetInitCtx()).Where(c.CrawlStatus, status).Where(c.HostType, hostType)
	if hostname != "" {
		where = where.Where(c.HostIp, hostname)
	}

	var seed *entity.CmsCrawlQueue
	_ = where.Scan(&seed)
	return seed
}

func GetNeedNotifySeedList() []*entity.CmsCrawlQueue {
	where := dao.CmsCrawlQueue.Ctx(gctx.GetInitCtx()).
		Where(c.CrawlM3U8Notify, CrawlM3U8NotifyNo).
		Where("crawl_m3u8_cnt >= ?", constant.ServerMaxRetry)
	var all []*entity.CmsCrawlQueue
	_ = where.Scan(&all)
	return all
}

func ExistCrawling(hostType int) bool {
	//爬虫状态为爬取中的数量大于0时 不继续
	//需要新增条件限制host_type不同类型限制 避免互相冲突
	count, _ := dao.CmsCrawlQueue.Ctx(gctx.GetInitCtx()).Where(c.HostType, hostType).Count(c.CrawlStatus, Crawling)
	return count > 0
}

func UpdateStatus(seed *entity.CmsCrawlQueue, status int) {
	if configservice.GetCrawlDebugBool("disableDB") {
		return
	}
	seed.CrawlStatus = status
	seed.UpdateTime = gtime.Now()
	dao.CmsCrawlQueue.Ctx(gctx.GetInitCtx()).Data(seed).Where(c.Id, seed.Id).Update()
}

func UpdateUrlAndStatus(seed *entity.CmsCrawlQueue) {
	if configservice.GetCrawlBool("disableDB") {
		return
	}
	seed.CrawlM3U8Cnt = seed.CrawlM3U8Cnt + 1
	seed.UpdateTime = gtime.Now()
	if seed.CrawlM3U8Url == "" && seed.CrawlM3U8Text == "" {
		if seed.CrawlM3U8Cnt >= constant.ServerMaxRetry {
			// 超过允许重试的最大次数
			seed.HostIp = configservice.GetCrawlHostIp()
			if seed.ErrorMsg == "" {
				seed.ErrorMsg = "M3U8 Empty"
			}
			UpdateStatus(seed, CrawlErr)
		} else {
			UpdateStatus(seed, Init)
		}
	} else {
		UpdateStatus(seed, CrawlFinish)
	}
}

// 重置处理中状态太久的
func ResetProcessingTooLong(context context.Context) {
	ResetHangingStatus(M3U8Parsing, CrawlFinish, 6*60)
}

// 重置抓取中状态太久的
func ResetCrawlingTooLong(context context.Context) {
	ResetHangingStatus(Crawling, Init, 1)
}

// 重置挂起中的状态
func ResetHangingStatus(fromStatus, toStatus, hangingMinutes int) {
	waterMark := gtime.Now().Add(time.Duration(hangingMinutes) * -time.Minute)

	var seed *entity.CmsCrawlQueue
	_ = dao.CmsCrawlQueue.Ctx(gctx.GetInitCtx()).
		Where("update_time < ", waterMark).
		Scan(&seed, c.CrawlStatus, fromStatus)

	UpdateStatus(seed, toStatus)
}
func ResetHostType2(context context.Context) {
	//waterMark := gtime.Now().Add(time.Duration(hangingMinutes) * -time.Minute)
	dao.CmsCrawlQueue.Ctx(gctx.GetInitCtx()).Where(g.Map{
		"host_type = ?":                       2,
		"(crawl_status =2 or crawl_status=5)": "",
	}).Update(g.Map{
		"crawl_status": 0,
	})

	//UpdateStatus(seed, toStatus)
}
