// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsCrawlReplayManifestTask is the golang structure for table cms_crawl_replay_manifest_task.
type CmsCrawlReplayManifestTask struct {
	Id             int         `json:"id"             ` // 主键ID
	ReplayConfigId int         `json:"replayConfigId" ` // 任务ID
	ReplayDay      string      `json:"replayDay"      ` // 节目日
	CrawlStatus    int         `json:"crawlStatus"    ` // 抓取状态.0-创建任务;1-抓取中;2-抓取失败;3-抓取完成;
	ErrorMsg       string      `json:"errorMsg"       ` // 错误信息
	CreateUser     int         `json:"createUser"     ` // 添加人
	CreateTime     *gtime.Time `json:"createTime"     ` // 添加时间
	UpdateUser     int         `json:"updateUser"     ` // 更新人
	UpdateTime     *gtime.Time `json:"updateTime"     ` // 更新时间
	Type           int         `json:"type"           ` // 节目单生成方式 0 按爬虫地址 1 自动生成
}
