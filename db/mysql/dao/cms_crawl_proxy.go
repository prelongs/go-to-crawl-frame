// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsCrawlProxyDao is internal type for wrapping internal DAO implements.
type internalCmsCrawlProxyDao = *internal.CmsCrawlProxyDao

// cmsCrawlProxyDao is the data access object for table cms_crawl_proxy.
// You can define custom methods on it to extend its functionality as you wish.
type cmsCrawlProxyDao struct {
	internalCmsCrawlProxyDao
}

var (
	// CmsCrawlProxy is globally public accessible object for table cms_crawl_proxy operations.
	CmsCrawlProxy = cmsCrawlProxyDao{
		internal.NewCmsCrawlProxyDao(),
	}
)

// Fill with you ideas below.