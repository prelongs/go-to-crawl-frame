// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsUploadQueue is the golang structure for table cms_upload_queue.
type CmsUploadQueue struct {
	Id           int         `json:"id"           ` // 主键ID
	HostIp       string      `json:"hostIp"       ` //
	CountryCode  string      `json:"countryCode"  ` // 国家二字码.(eg: CN,US,SG等)
	VideoYear    int         `json:"videoYear"    ` // 视频发布年份
	VideoCollId  int64       `json:"videoCollId"  ` // 视频集ID（视频集ID，不限于电视剧,-1代表单集视频，或者说电影）
	VideoItemId  int64       `json:"videoItemId"  ` // 视频集对应视频项ID（不限于电视剧的剧集）
	FileName     string      `json:"fileName"     ` // 文件标题
	FileType     int         `json:"fileType"     ` // 文件类型. 1-视频；2-大体积资源；（小文件无需用队列，直接用上传接口）
	FileSize     int64       `json:"fileSize"     ` // 文件大小. 单位KB
	Msg          string      `json:"msg"          ` //
	UploadStatus int         `json:"uploadStatus" ` // 上传状态.0-创建任务;1-上传中;2-上传完成;3-流媒体处理中;4-流媒体处理结束;5-流媒体处理异常
	CreateUser   int         `json:"createUser"   ` // 添加人
	CreateTime   *gtime.Time `json:"createTime"   ` // 添加时间
	UpdateUser   int         `json:"updateUser"   ` // 更新人
	UpdateTime   *gtime.Time `json:"updateTime"   ` // 更新时间
}
