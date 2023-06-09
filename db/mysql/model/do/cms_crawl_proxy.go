// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsCrawlProxy is the golang structure of table cms_crawl_proxy for DAO operations like Where/Data.
type CmsCrawlProxy struct {
	g.Meta       `orm:"table:cms_crawl_proxy, do:true"`
	Id           interface{} // 主键ID
	TargetDomain interface{} // 目标网站顶级域名
	ProxyUrl     interface{} // 代理地址
	ProxyStatus  interface{} // 代理状态. 0-停用,1-使用中
	CreateUser   interface{} // 添加人
	CreateTime   *gtime.Time // 添加时间
	UpdateUser   interface{} // 更新人
	UpdateTime   *gtime.Time // 更新时间
}
