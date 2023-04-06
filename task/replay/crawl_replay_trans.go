package replay

import (
	"context"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/configservice"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/replayservice"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/videoservice"
	"github.com/JervisPG/go-to-crawl-frame/utils/ffmpegutil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

var Transform = new(CrawlReplayTransform)

type CrawlReplayTransform struct {
	programTask *entity.CmsCrawlReplayProgramTask
}

func (receiver *CrawlReplayTransform) TransformTask(context context.Context) {
	programTask := replayservice.GetOneByIpStatus(configservice.GetCrawlHostIp(), replayservice.ProgramTaskCrawlFinish)
	if programTask == nil {
		return
	}
	replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskParsing)
	taskDO := replayservice.GetDetailById(programTask.Id)
	savePath := videoservice.GetReplayDir(taskDO)

	mp4File := ffmpegutil.GetGenericFilePath(savePath)
	err := ffmpegutil.RunFfmpegGenericSlice(savePath)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), err)
		replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskParseErr)
	} else {
		replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskParsed)
		_ = gfile.Remove(mp4File)
	}

}
