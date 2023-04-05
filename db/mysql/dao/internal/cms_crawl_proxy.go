// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CmsCrawlProxyDao is the data access object for table cms_crawl_proxy.
type CmsCrawlProxyDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns CmsCrawlProxyColumns // columns contains all the column names of Table for convenient usage.
}

// CmsCrawlProxyColumns defines and stores column names for table cms_crawl_proxy.
type CmsCrawlProxyColumns struct {
	Id           string // 主键ID
	TargetDomain string // 目标网站顶级域名
	ProxyUrl     string // 代理地址
	ProxyStatus  string // 代理状态. 0-停用,1-使用中
	CreateUser   string // 添加人
	CreateTime   string // 添加时间
	UpdateUser   string // 更新人
	UpdateTime   string // 更新时间
}

// cmsCrawlProxyColumns holds the columns for table cms_crawl_proxy.
var cmsCrawlProxyColumns = CmsCrawlProxyColumns{
	Id:           "id",
	TargetDomain: "target_domain",
	ProxyUrl:     "proxy_url",
	ProxyStatus:  "proxy_status",
	CreateUser:   "create_user",
	CreateTime:   "create_time",
	UpdateUser:   "update_user",
	UpdateTime:   "update_time",
}

// NewCmsCrawlProxyDao creates and returns a new DAO object for table data access.
func NewCmsCrawlProxyDao() *CmsCrawlProxyDao {
	return &CmsCrawlProxyDao{
		group:   "default",
		table:   "cms_crawl_proxy",
		columns: cmsCrawlProxyColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CmsCrawlProxyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CmsCrawlProxyDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CmsCrawlProxyDao) Columns() CmsCrawlProxyColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CmsCrawlProxyDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CmsCrawlProxyDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CmsCrawlProxyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
