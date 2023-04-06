package crawldto

import "github.com/JervisPG/go-to-crawl-frame/entity/cmsdto"

type CmsUploadQueueCreate struct {
	HostIp string `p:"hostIp" v:"required#主机IP为空"`
	//CreateUser  int    `p:"createUser" v:"required#登录信息过期"`
	CountryCode string `p:"countryCode" v:"required#国家编码为空"`
	VideoCollId int64  `p:"videoCollId" v:"required#剧ID为空"`
	VideoItemId int64  `p:"videoItemId" v:"required#剧集ID为空"`
}

type CmsUploadQueueQry struct {
	*cmsdto.CmsBasePageQry
	UploadStatus int `p:"uploadStatus"`
}

type CmsUploadQueueCmd struct {
	Id           int64 `p:"id" v:"required#任务ID为空"`
	UploadStatus int   `p:"uploadStatus" v:"required#任务状态"`
}
