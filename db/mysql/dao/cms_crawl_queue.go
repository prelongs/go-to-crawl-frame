// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsCrawlQueueDao is internal type for wrapping internal DAO implements.
type internalCmsCrawlQueueDao = *internal.CmsCrawlQueueDao

// cmsCrawlQueueDao is the data access object for table cms_crawl_queue.
// You can define custom methods on it to extend its functionality as you wish.
type cmsCrawlQueueDao struct {
	internalCmsCrawlQueueDao
}

var (
	// CmsCrawlQueue is globally public accessible object for table cms_crawl_queue operations.
	CmsCrawlQueue = cmsCrawlQueueDao{
		internal.NewCmsCrawlQueueDao(),
	}
)

// Fill with you ideas below.