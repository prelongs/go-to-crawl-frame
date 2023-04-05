package crawltask

import (
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/crawl/servicedto"
	"testing"
)

func TestBananTV(t *testing.T) {
	doStartVodTV(Banan, "https://banan.tv/vodplay/78923-1-1.html")
}

func TestBananTVPadInfo(t *testing.T) {
	doStartVodPadInfo("https://banan.tv/vodplay/78923-1-1.html")
}

// 偏单元测试
func TestNiVodTV(t *testing.T) {
	doStartVodTV(NiVod, "https://www.mudvod.tv/filter.html?x=1&channelId=7&showTypeId=145")
}

func TestNiVodTVPadInfo(t *testing.T) {
	doStartVodPadInfo("https://www.mudvod.tv/l01tKF0bT0wpcHiJVV2CsJpLH3gDDP8T-0-0-0-0-detail.html?x=1")
}

func doStartVodTV(domain, url string) {
	seed := new(servicedto.CrawlVodConfigDO)
	seed.CmsCrawlVodConfig = new(entity.CmsCrawlVodConfig)
	seed.SeedUrl = url
	seed.PageSize = 6
	seed.DomainKeyPart = domain

	seed.CmsCrawlVodConfigTask = new(entity.CmsCrawlVodConfigTask)
	seed.VodConfigId = -1
	DoStartCrawlVodTV(seed)
}

func doStartVodPadInfo(seedUrl string) {
	vodTvItem := new(entity.CmsCrawlVodTv)
	vodTvItem.SeedUrl = seedUrl
	DoStartCrawlVodPadInfo(vodTvItem)
}

func TestGenVodConfigTask(t *testing.T) {
	CrawlVodTVTask.GenVodConfigTask(gctx.GetInitCtx())
}

func TestVodTVTask(t *testing.T) {
	CrawlVodTVTask.VodTVTask(gctx.GetInitCtx())
}

func TestVodTVPadInfoTask(t *testing.T) {
	CrawlVodTVTask.VodTVPadInfoTask(gctx.GetInitCtx())
}
