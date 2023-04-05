// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsCrawlReplayConfig is the golang structure for table cms_crawl_replay_config.
type CmsCrawlReplayConfig struct {
	Id            int         `json:"id"            ` // 主键ID
	ChannelNo     string      `json:"channelNo"     ` // 频道编号
	ChannelName   string      `json:"channelName"   ` // 频道名称
	ChannelNameEn string      `json:"channelNameEn" ` //
	ChannelIcon   string      `json:"channelIcon"   ` //
	Type          int         `json:"type"          ` // 节目单生成方式 0 按爬虫地址 1 自动生成
	SeedUrl       string      `json:"seedUrl"       ` // 种子URL
	SeedParams    string      `json:"seedParams"    ` // 种子URL携带的参数。保存Json串
	Host          string      `json:"host"          ` // 配置录制的机器
	Domain        string      `json:"domain"        ` //
	PlayUrl       string      `json:"playUrl"       ` // 播放地址(流)
	ErrorMsg      string      `json:"errorMsg"      ` // 错误信息
	CreateUser    int         `json:"createUser"    ` // 添加人
	CreateTime    *gtime.Time `json:"createTime"    ` // 添加时间
	UpdateUser    int         `json:"updateUser"    ` // 更新人
	UpdateTime    *gtime.Time `json:"updateTime"    ` // 更新时间
	Status        int         `json:"status"        ` // 状态：1在用 2停用
	Sort          int         `json:"sort"          ` // 排序
	Note          string      `json:"note"          ` // 备注
	Mark          int         `json:"mark"          ` // 有效标识(1正常 0删除)
}
