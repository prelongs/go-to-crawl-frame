// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsCrawlDictDao is internal type for wrapping internal DAO implements.
type internalCmsCrawlDictDao = *internal.CmsCrawlDictDao

// cmsCrawlDictDao is the data access object for table cms_crawl_dict.
// You can define custom methods on it to extend its functionality as you wish.
type cmsCrawlDictDao struct {
	internalCmsCrawlDictDao
}

var (
	// CmsCrawlDict is globally public accessible object for table cms_crawl_dict operations.
	CmsCrawlDict = cmsCrawlDictDao{
		internal.NewCmsCrawlDictDao(),
	}
)

// Fill with you ideas below.