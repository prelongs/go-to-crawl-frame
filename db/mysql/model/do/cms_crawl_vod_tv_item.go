// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsCrawlVodTvItem is the golang structure of table cms_crawl_vod_tv_item for DAO operations like Where/Data.
type CmsCrawlVodTvItem struct {
	g.Meta      `orm:"table:cms_crawl_vod_tv_item, do:true"`
	Id          interface{} //
	TvId        interface{} //
	TvItemMd5   interface{} //
	CrawlStatus interface{} //
	SeedUrl     interface{} //
	SeedParams  interface{} //
	ErrorMsg    interface{} //
	VideoCollId interface{} //
	VideoItemId interface{} //
	Episodes    interface{} //
	CreateUser  interface{} //
	CreateTime  *gtime.Time //
	UpdateUser  interface{} //
	UpdateTime  *gtime.Time //
}