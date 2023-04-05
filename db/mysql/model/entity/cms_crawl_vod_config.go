// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsCrawlVodConfig is the golang structure for table cms_crawl_vod_config.
type CmsCrawlVodConfig struct {
	Id            int         `json:"id"            ` // 主键ID
	HostType      int         `json:"hostType"      ` // 传递给crawlQueue的hostType字段。2-nivod网；3-BananTV
	VodType       int         `json:"vodType"       ` // 点播类型.0-电影；1-剧集（标志给展示逻辑，爬虫统一按剧集逻辑走）
	DomainKeyPart string      `json:"domainKeyPart" ` // 域名关键部分.用于配置策略
	ProgramNo     string      `json:"programNo"     ` // 栏目编号
	ProgramName   string      `json:"programName"   ` // 栏目名称
	ProgramIcon   string      `json:"programIcon"   ` // 栏目图标
	CategoryNo    string      `json:"categoryNo"    ` // 分类编号
	CategoryName  string      `json:"categoryName"  ` // 分类名称
	SeedUrl       string      `json:"seedUrl"       ` // 种子URL
	SeedParams    string      `json:"seedParams"    ` // 种子URL携带的参数。保存Json串
	SeedStatus    int         `json:"seedStatus"    ` // 状态：1在用 2停用
	PageSize      int         `json:"pageSize"      ` // 翻页次数
	SeedDesc      string      `json:"seedDesc"      ` // 描述
	ErrorMsg      string      `json:"errorMsg"      ` // 错误信息
	CreateUser    int         `json:"createUser"    ` // 添加人
	CreateTime    *gtime.Time `json:"createTime"    ` // 添加时间
	UpdateUser    int         `json:"updateUser"    ` // 更新人
	UpdateTime    *gtime.Time `json:"updateTime"    ` // 更新时间
}