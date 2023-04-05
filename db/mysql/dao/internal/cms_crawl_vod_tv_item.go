// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CmsCrawlVodTvItemDao is the data access object for table cms_crawl_vod_tv_item.
type CmsCrawlVodTvItemDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns CmsCrawlVodTvItemColumns // columns contains all the column names of Table for convenient usage.
}

// CmsCrawlVodTvItemColumns defines and stores column names for table cms_crawl_vod_tv_item.
type CmsCrawlVodTvItemColumns struct {
	Id          string //
	TvId        string //
	TvItemMd5   string //
	CrawlStatus string //
	SeedUrl     string //
	SeedParams  string //
	ErrorMsg    string //
	VideoCollId string //
	VideoItemId string //
	Episodes    string //
	CreateUser  string //
	CreateTime  string //
	UpdateUser  string //
	UpdateTime  string //
}

// cmsCrawlVodTvItemColumns holds the columns for table cms_crawl_vod_tv_item.
var cmsCrawlVodTvItemColumns = CmsCrawlVodTvItemColumns{
	Id:          "id",
	TvId:        "tv_id",
	TvItemMd5:   "tv_item_md5",
	CrawlStatus: "crawl_status",
	SeedUrl:     "seed_url",
	SeedParams:  "seed_params",
	ErrorMsg:    "error_msg",
	VideoCollId: "video_coll_id",
	VideoItemId: "video_item_id",
	Episodes:    "episodes",
	CreateUser:  "create_user",
	CreateTime:  "create_time",
	UpdateUser:  "update_user",
	UpdateTime:  "update_time",
}

// NewCmsCrawlVodTvItemDao creates and returns a new DAO object for table data access.
func NewCmsCrawlVodTvItemDao() *CmsCrawlVodTvItemDao {
	return &CmsCrawlVodTvItemDao{
		group:   "default",
		table:   "cms_crawl_vod_tv_item",
		columns: cmsCrawlVodTvItemColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CmsCrawlVodTvItemDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CmsCrawlVodTvItemDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CmsCrawlVodTvItemDao) Columns() CmsCrawlVodTvItemColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CmsCrawlVodTvItemDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CmsCrawlVodTvItemDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CmsCrawlVodTvItemDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
