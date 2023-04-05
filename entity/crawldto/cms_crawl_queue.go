package crawldto

type CmsCrawlQueueCreate struct {
	//CreateUser      int    `p:"createUser" v:"required#登录信息过期"`
	HostIp          string `p:"host_ip" v:"required#下载IP为空"`
	CountryCode     string `p:"countryCode" v:"required#国家编码为空"`
	VideoYear       int    `p:"videoYear" v:"required#年份为空"`
	VideoCollId     int64  `p:"videoCollId" v:"required#剧ID为空"`
	VideoItemId     int64  `p:"videoItemId" v:"required#剧集ID为空"`
	CrawlType       int    `p:"crawlType" v:"required#抓取类型.1-页面URL;2-文件m3u8"`
	CrawlSeedUrl    string `p:"crawlSeedUrl"`
	CrawlSeedParams string `p:"crawlSeedParams"`
	CrawlM3U8Url    string `p:"crawlM3U8Url"`
}

type CmsCrawlQueueQry struct {
	VideoItemId int64 `p:"videoItemId" v:"required#剧集ID为空"`
}
