package uploadservice

import (
	"github.com/gogf/gf/v2/os/gctx"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
)

var (
	columns = dao.CmsUploadQueue.Columns()
)

func GetByVideoItemId(videoItemId int64, status int) *entity.CmsUploadQueue {
	var one *entity.CmsUploadQueue
	_ = dao.CmsUploadQueue.Ctx(gctx.GetInitCtx()).
		Where(columns.VideoItemId, videoItemId).
		Where(columns.UploadStatus, status).Scan(&one)
	return one
}

func GetById(id int64) *entity.CmsUploadQueue {
	var one *entity.CmsUploadQueue
	_ = dao.CmsUploadQueue.
		Ctx(gctx.GetInitCtx()).
		Where(columns.Id, id).
		Scan(&one)
	return one
}

func UpdateById(queue *entity.CmsUploadQueue, status int) {
	queue.UploadStatus = status
	where := dao.CmsUploadQueue.Ctx(gctx.GetInitCtx()).Data(queue).Where(columns.Id, queue.Id)
	where.Update()
}
