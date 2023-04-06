package vodservice

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/dao"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/configservice"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// 转换到爬取队列走标准化抓取
func TransToCrawlQueue(vodTv *entity.CmsCrawlVodTv, vodTvItem *entity.CmsCrawlVodTvItem) {

	vodConfig := GetVodConfigById(vodTv.VodConfigId)
	hostType := HostTypeNiVod
	if vodConfig != nil && vodConfig.HostType != 0 {
		hostType = vodConfig.HostType
	}
	crawlQueue := new(entity.CmsCrawlQueue)
	crawlQueue.HostIp = configservice.GetString("crawl.down_ip")
	crawlQueue.HostType = hostType

	crawlQueue.VideoYear = gconv.Int(vodTv.VideoYear)
	crawlQueue.VideoCollId = vodTv.VideoCollId
	crawlQueue.VideoItemId = vodTvItem.VideoItemId
	crawlQueue.CrawlType = TypePageUrl
	crawlQueue.CrawlStatus = Init
	crawlQueue.CrawlSeedUrl = vodTvItem.SeedUrl
	crawlQueue.CountryCode = getCountryCodeByString(vodTv.VideoCountry)
	crawlQueue.CrawlSeedParams = ""
	crawlQueue.CreateTime = gtime.Now()

	_, _ = dao.CmsCrawlQueue.Ctx(gctx.GetInitCtx()).Save(crawlQueue)
}

func getCountryCodeByString(country string) string {
	//大陆 香港 台湾 日本 韩国 欧美 泰国 新马 其它
	defaultCry := "OTHER"

	if country == "" {
		return defaultCry
	}

	var cryMap = make(map[string]string)
	cryMap["大陆"] = "CN"
	cryMap["香港"] = "HK"
	cryMap["台湾"] = "TW"
	cryMap["日本"] = "JP"
	cryMap["韩国"] = "KR"
	cryMap["欧美"] = "US"
	cryMap["泰国"] = "TH"
	cryMap["新马"] = "MY"
	cryMap["其它"] = defaultCry
	if v, ok := cryMap[country]; ok {
		return v
	} else {
		return defaultCry
	}
}
