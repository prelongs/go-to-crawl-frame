package replay

import (
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/task/taskdto"
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-video/task/replaytask/astro"
	"go-to-crawl-video/task/replaytask/autolist"
	"go-to-crawl-video/task/replaytask/cntv"
	"go-to-crawl-video/task/replaytask/epg51zmt"
	"go-to-crawl-video/task/replaytask/programme"
	"go-to-crawl-video/task/replaytask/programtable"
	"go-to-crawl-video/task/replaytask/sports8"
	"go-to-crawl-video/task/replaytask/tbc"
	"go-to-crawl-video/task/replaytask/tvkingdom"
	"go-to-crawl-video/task/replaytask/viu"
)

func getCrawlReplayStrategy(replayConfig *entity.CmsCrawlReplayConfig) taskdto.CrawlReplayInterface {
	//g.Dump(replayConfig)
	if replayConfig.Type == 0 {
		//	爬虫类型
		//g.Dump("replayConfig.SeedUrl:", replayConfig.SeedUrl)
		if gstr.Contains(replayConfig.SeedUrl, "astro") {
			return new(astro.AstroCrawl)
		} else if gstr.Contains(replayConfig.SeedUrl, "api.cntv.cn") {
			return new(cntv.Cntv)
		} else if gstr.Contains(replayConfig.SeedUrl, "programme.tvb") {
			return new(programme.ProgrammeCrawl)
		} else if gstr.Contains(replayConfig.SeedUrl, "tbc.net.tw") {
			return new(tbc.TbcCrawl)
		} else if gstr.Contains(replayConfig.SeedUrl, "sports8.net") {
			return new(sports8.Sports8Crawl)
		} else if gstr.Contains(replayConfig.SeedUrl, "epg.51zmt") {
			//海外免费epgxml格式解析
			return new(epg51zmt.Epg51zmt)
		} else if gstr.Contains(replayConfig.SeedUrl, "tvkingdom.jp") {
			return new(tvkingdom.TvkingdomCrawl)
		} else if gstr.Contains(replayConfig.SeedUrl, "api.viu.tv") {
			return new(viu.ViuTvCrawl)
		} else if gstr.Contains(replayConfig.SeedUrl, "節目表.tw") {
			return new(programtable.ProgramtableCrawl)
		}

	} else if replayConfig.Type == 1 {
		//自动生成

		return new(autolist.AutoList)
	}
	return nil
}
