package crawltask

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/configservice"
	"go-to-crawl-frame/service/crawl/uploadservice"
)

// 把任务队列里的视频资源转换成M3U8格式视频资源
func PostCmsTask(context context.Context) {
	log := g.Log().Line()

	//os.Exit(1)
	//查找数据库上传完毕状态的数据进行处理
	var queue *entity.CmsUploadQueue
	err := dao.CmsUploadQueue.Ctx(gctx.GetInitCtx()).Scan(&queue,
		g.Map{
			columns.UploadStatus: uploadservice.Transformed,
			columns.HostIp:       configservice.GetCrawlCfg("hostIp"),
		})
	if err != nil {
		return
	}
	if queue == nil {
		return
	}
	//视频处理完毕通知CMS处理
	info, _ := g.NewVar(queue).MarshalJSON()

	log.Infof(gctx.GetInitCtx(),
		"DoCallBack:%v->Param:%v", configservice.GetCrawlCfg("callback_url"), string(info))
	postRes := g.Client().PostContent(gctx.GetInitCtx(), configservice.GetCrawlCfg("callback_url"), g.Map{"info": info})
	if !g.NewVar(postRes).IsEmpty() {
		if !g.NewVar(g.NewVar(postRes).Map()["code"]).IsEmpty() {
			if g.NewVar(g.NewVar(postRes).Map()["code"]).Int() == 0 {
				//通知cms返回成功则更新状态为CmsPostSuccess（方便后续做新定时任务检测状态为upload.Transformed但是请求cms不成功的 重新发送请求）
				queue.UploadStatus = uploadservice.CmsPostSuccess
				queue.UpdateTime = gtime.Now()
				_, _ = dao.CmsUploadQueue.Ctx(gctx.GetInitCtx()).Data(queue).Where(columns.Id, queue.Id).Update()
			}
		}
	}

}
