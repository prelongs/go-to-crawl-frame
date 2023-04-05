// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsCrawlVodConfigTaskDao is internal type for wrapping internal DAO implements.
type internalCmsCrawlVodConfigTaskDao = *internal.CmsCrawlVodConfigTaskDao

// cmsCrawlVodConfigTaskDao is the data access object for table cms_crawl_vod_config_task.
// You can define custom methods on it to extend its functionality as you wish.
type cmsCrawlVodConfigTaskDao struct {
	internalCmsCrawlVodConfigTaskDao
}

var (
	// CmsCrawlVodConfigTask is globally public accessible object for table cms_crawl_vod_config_task operations.
	CmsCrawlVodConfigTask = cmsCrawlVodConfigTaskDao{
		internal.NewCmsCrawlVodConfigTaskDao(),
	}
)

// Fill with you ideas below.