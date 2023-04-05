package service

import (
	"github.com/gogf/gf/v2/frame/g"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/entity/cmsdto"
	"go-to-crawl-frame/service/common"
)

type sNewCms struct {
}

func NewCms() *sNewCms {
	return &sNewCms{}
}

func (cms *sNewCms) ListByParam(option *cmsdto.ListReq) (list interface{}, count int, err error) {
	return common.DbNewCms().ReadByParam(option)
}

func (cms *sNewCms) UpByParam(option *cmsdto.UpdateReq, UserId int) (rows int64, err error) {
	return common.DbNewCms().UpByParam(option, UserId)
}

func (cms *sNewCms) AddByParam(option *cmsdto.AddReq, UserId int) (rows int64, add_id int64, err error) {
	return common.DbNewCms().AddByParam(option, UserId)
}
func (cms *sNewCms) DelByParam(option *cmsdto.DelReq) (rows int64, err error) {
	return common.DbNewCms().DelByParam(option)
}
func (cms *sNewCms) VodUpSuccess(queue entity.CmsUploadQueue) (err error) {
	where := g.Map{
		"video_id": queue.VideoCollId,
		"id":       queue.VideoItemId,
	}
	return common.DbNewCms().UpVideoStatus(where)
}
