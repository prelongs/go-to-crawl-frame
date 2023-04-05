// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsCrawlVodTvItem is the golang structure for table cms_crawl_vod_tv_item.
type CmsCrawlVodTvItem struct {
	Id          int         `json:"id"          ` //
	TvId        int         `json:"tvId"        ` //
	TvItemMd5   string      `json:"tvItemMd5"   ` //
	CrawlStatus int         `json:"crawlStatus" ` //
	SeedUrl     string      `json:"seedUrl"     ` //
	SeedParams  string      `json:"seedParams"  ` //
	ErrorMsg    string      `json:"errorMsg"    ` //
	VideoCollId int64       `json:"videoCollId" ` //
	VideoItemId int64       `json:"videoItemId" ` //
	Episodes    string      `json:"episodes"    ` //
	CreateUser  int         `json:"createUser"  ` //
	CreateTime  *gtime.Time `json:"createTime"  ` //
	UpdateUser  int         `json:"updateUser"  ` //
	UpdateTime  *gtime.Time `json:"updateTime"  ` //
}