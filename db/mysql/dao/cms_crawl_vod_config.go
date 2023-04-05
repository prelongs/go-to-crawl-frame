// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsCrawlVodConfigDao is internal type for wrapping internal DAO implements.
type internalCmsCrawlVodConfigDao = *internal.CmsCrawlVodConfigDao

// cmsCrawlVodConfigDao is the data access object for table cms_crawl_vod_config.
// You can define custom methods on it to extend its functionality as you wish.
type cmsCrawlVodConfigDao struct {
	internalCmsCrawlVodConfigDao
}

var (
	// CmsCrawlVodConfig is globally public accessible object for table cms_crawl_vod_config operations.
	CmsCrawlVodConfig = cmsCrawlVodConfigDao{
		internal.NewCmsCrawlVodConfigDao(),
	}
)

// Fill with you ideas below.