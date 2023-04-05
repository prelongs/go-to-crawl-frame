package servicedto

import "github.com/gogf/gf/v2/os/gtime"

type CmsCrawlReplayProgramTaskDO struct {
	Id               int         `json:"id"`               // 主键ID
	ManifestId       int         `json:"manifestId"`       // 队列ID
	ProgramNo        string      `json:"programNo"`        // 节目编号
	ProgramName      string      `json:"programName"`      // 节目名称
	ProgramStartTime *gtime.Time `json:"programStartTime"` // 节目开始时间
	ProgramEndTime   *gtime.Time `json:"programEndTime"`   // 节目结束时间
	CrawlStatus      int         `json:"crawlStatus"`      // 抓取状态.0-创建任务;1-录制中;2-录制失败;3-录制完成;
	ErrorMsg         string      `json:"errorMsg"`         // 错误信息
	CreateUser       int         `json:"createUser"`       // 添加人
	CreateTime       *gtime.Time `json:"createTime"`       // 添加时间
	UpdateUser       int         `json:"updateUser"`       // 更新人
	UpdateTime       *gtime.Time `json:"updateTime"`       // 更新时间

	ReplayDay string `json:"replayDay"` // 节目日

	ChannelNo   string `json:"channelNo"`   // 频道编号
	ChannelName string `json:"channelName"` // 频道名称
	PlayUrl     string `json:"playUrl"`     // 播放地址(流)
}
