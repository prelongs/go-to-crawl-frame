// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsCrawlVodTvItemDao is internal type for wrapping internal DAO implements.
type internalCmsCrawlVodTvItemDao = *internal.CmsCrawlVodTvItemDao

// cmsCrawlVodTvItemDao is the data access object for table cms_crawl_vod_tv_item.
// You can define custom methods on it to extend its functionality as you wish.
type cmsCrawlVodTvItemDao struct {
	internalCmsCrawlVodTvItemDao
}

var (
	// CmsCrawlVodTvItem is globally public accessible object for table cms_crawl_vod_tv_item operations.
	CmsCrawlVodTvItem = cmsCrawlVodTvItemDao{
		internal.NewCmsCrawlVodTvItemDao(),
	}
)

// Fill with you ideas below.