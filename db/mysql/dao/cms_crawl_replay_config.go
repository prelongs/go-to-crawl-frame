// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"go-to-crawl-frame/db/mysql/dao/internal"
)

// internalCmsCrawlReplayConfigDao is internal type for wrapping internal DAO implements.
type internalCmsCrawlReplayConfigDao = *internal.CmsCrawlReplayConfigDao

// cmsCrawlReplayConfigDao is the data access object for table cms_crawl_replay_config.
// You can define custom methods on it to extend its functionality as you wish.
type cmsCrawlReplayConfigDao struct {
	internalCmsCrawlReplayConfigDao
}

var (
	// CmsCrawlReplayConfig is globally public accessible object for table cms_crawl_replay_config operations.
	CmsCrawlReplayConfig = cmsCrawlReplayConfigDao{
		internal.NewCmsCrawlReplayConfigDao(),
	}
)

// Fill with you ideas below.