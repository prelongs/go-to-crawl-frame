package crawltask

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/vodservice"
	"github.com/JervisPG/go-to-crawl-frame/task/taskdto"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-video/task/vodtask/banan"
	"go-to-crawl-video/task/vodtask/bilibili"
	"go-to-crawl-video/task/vodtask/iqiyi"
	"go-to-crawl-video/task/vodtask/nivod"
	"go-to-crawl-video/task/vodtask/nunuyy"
	"go-to-crawl-video/task/vodtask/olevod"
	"go-to-crawl-video/task/vodtask/qq"
	"go-to-crawl-video/task/vodtask/tangrenjie"
)

const (
	Nunuyy     = "nunuyy"
	Bilibbili  = "bilibili"
	Ole        = "ole"
	TangRenJie = "tangrenjie.tv"
	QQ         = "v.qq.com"
	NiVod      = "nivod\\d*\\.tv"
	MudVod     = "mudvod\\d*\\.tv"
	Banan      = "banan.tv"
	Iqiyi      = "iqiyi.com"
)

func getCrawlVodFlowStrategy(seed *entity.CmsCrawlQueue) taskdto.CrawlVodFlowInterface {
	url := seed.CrawlSeedUrl
	return doGetCrawlVodFlowStrategy(url)
}

func doGetCrawlVodFlowStrategy(url string) taskdto.CrawlVodFlowInterface {
	if gstr.Contains(url, Nunuyy) {
		return new(nunuyy.NunuyyCrawl)
	} else if gstr.Contains(url, Bilibbili) {
		return new(bilibili.BilibiliCrawl)
	} else if gstr.Contains(url, Ole) {
		return new(olevod.OleVodCrawl)
	} else if gstr.Contains(url, TangRenJie) {
		return new(tangrenjie.TangRenJieCrawl)
	} else if gstr.Contains(url, QQ) {
		return new(qq.QQCrawl)
	} else if gregex.IsMatchString(NiVod, url) || gregex.IsMatchString(MudVod, url) {
		return new(nivod.NiVodCrawl)
	} else if gstr.Contains(url, Banan) {
		return new(banan.BananCrawl)
	} else if gstr.Contains(url, Iqiyi) {
		return new(iqiyi.IqiyiMobileCrawl)
	}

	return nil
}

func getCrawlVodTVStrategy(seed *entity.CmsCrawlVodConfig) taskdto.CrawlVodTVInterface {

	url := seed.SeedUrl

	if gstr.Contains(url, NiVod) || gstr.ContainsI(url, MudVod) {
		return new(nivod.NiVodTVTask)
	} else if gstr.Contains(url, Banan) {
		return new(banan.BananTvCrawl)
	}

	return nil
}

func getCrawlVodPadInfoStrategy(seed *entity.CmsCrawlVodTv) taskdto.CrawlVodTVInterface {

	url := seed.SeedUrl
	if gstr.Contains(url, NiVod) || gstr.ContainsI(url, MudVod) {
		return new(nivod.NiVodTVTask)
	} else if gstr.Contains(url, Banan) {
		return new(banan.BananTvCrawl)
	}

	return nil
}

func GetHostType(crawlSeedUrl string) int {
	if gstr.Contains(crawlSeedUrl, QQ) {
		return vodservice.HostTypeCrawlLogin
	} else if gstr.Contains(crawlSeedUrl, Iqiyi) {
		return vodservice.HostTypeCrawlLogin
	} else {
		return vodservice.HostTypeNormal
	}
}
