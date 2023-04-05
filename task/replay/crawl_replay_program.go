package replay

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/configservice"
	"go-to-crawl-frame/service/crawl/replayservice"
	"go-to-crawl-frame/service/crawl/videoservice"
	"go-to-crawl-frame/utils/ffmpegutil"

	"net/http"
	"os"
	"time"
)

var Program = new(CrawlReplayProgram)

type CrawlReplayProgram struct {
	programTask        *entity.CmsCrawlReplayProgramTask
	startRecordingTime *gtime.Time
}

func doCreate(receiver *CrawlReplayProgram, programTask *entity.CmsCrawlReplayProgramTask, sTime *gtime.Time, fixMode bool) {
	taskName := fmt.Sprintf("RecordingTask-%d", programTask.Id)
	task := gcron.Search(taskName)

	if task == nil {
		r := &CrawlReplayProgram{
			programTask:        programTask,
			startRecordingTime: sTime,
		}
		doFinalCreate(r, sTime, programTask, taskName, fixMode)
	}
}

// fixMode: 修复模式
func doFinalCreate(receiver *CrawlReplayProgram, sTime *gtime.Time, programTask *entity.CmsCrawlReplayProgramTask, taskName string, fixMode bool) {

	job := receiver.StartRecording
	if fixMode {
		job = receiver.StartFixRecording
		g.Log().Info(gctx.GetInitCtx(), "修复录制: taskId = ", programTask.Id, ", programName = ", programTask.ProgramName)
	}

	pattern := fmt.Sprintf("%d %d %d * * %s", sTime.Second(), sTime.Minute(), sTime.Hour(), sTime.Format("D"))
	receiver.programTask = programTask
	receiver.startRecordingTime = sTime
	task, _ := gcron.AddOnce(gctx.GetInitCtx(), pattern, job, taskName)
	g.Log().Infof(gctx.GetInitCtx(), "动态新增录制任务: cron = %s, taskName = %s, programName = %s", pattern, task.Name, programTask.ProgramName)
	gcron.Start(taskName)
}

// 创建定时任务列表-未来需要录制的
func (receiver *CrawlReplayProgram) CreateAllFutureRecordingTask(context context.Context) {
	programTaskList := replayservice.GetAllByStatus(replayservice.ProgramTaskInit)
	//glog.Infof(gctx.GetInitCtx(),"cronList:%v", len(programTaskList))
	for _, programTask := range programTaskList {
		taskStatTime := programTask.ProgramStartTime
		doCreate(receiver, programTask, taskStatTime, false)
	}
}
func (receiver *CrawlReplayProgram) Test() {

	//task,_:=dao.CmsCrawlReplayProgramTask.Where(g.Map{"crawl_status":3}).FindOne()
	//do:=replay.GetDetailById(task.Id)
	//res:=videoservice.GetServerUrl(do)
	//res2:=videoservice.GetReplayDir(do)
	//g.Dump(res)
	//g.Dump(res2)
	res := `{"code":0,"data":{"code":100,"data":"{\"uuid\": \"487895cae7a9c9110f68c90a82d43a8f\"}","error":""}}`
	var replayRes *replayservice.ReplayRes
	err := g.NewVar(res).Structs(&replayRes)
	if err != nil {
		g.Dump(err.Error())
	}
	g.Dump(replayRes.Data.Data)
	var uuid replayservice.ReplayUUID
	g.NewVar(replayRes.Data.Data).Struct(&uuid)
	g.Dump(uuid)

}

// 创建定时任务列表-当前需要录制的
func (receiver *CrawlReplayProgram) CreateAllCurrentRecordingTask(context context.Context) {
	programTaskList := replayservice.GetAllCurrentByStatus(replayservice.ProgramTaskInit)
	for _, programTask := range programTaskList {
		taskStatTime := gtime.Now().Add(gtime.S * 3) // 给一个3秒延迟，防止pattern表达式过期
		doCreate(receiver, programTask, taskStatTime, false)
	}
}

func (receiver *CrawlReplayProgram) DeleCache(context context.Context) {
	deleDay := configservice.GetInt("replay.day", 7)
	receiver.DelBeForDay(deleDay)
}
func (receiver *CrawlReplayProgram) DeleCacheSql(context context.Context) {
	deleDay := configservice.GetInt("replay.day", 7)
	receiver.DelBeForDaySql(deleDay)
}
func (receiver *CrawlReplayProgram) DelBeForDaySql(day int) {
	dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Where(g.Map{
		"program_start_time < DATE_SUB(now(), INTERVAL ? DAY)": day,
	}).Delete()
}
func (receiver *CrawlReplayProgram) PostSuccessReplay() {
	//g.Dump("start")
	var programTask *entity.CmsCrawlReplayProgramTask
	err := dao.CmsCrawlReplayProgramTask.Ctx(gctx.GetInitCtx()).Where(g.Map{
		"crawl_status = ?": 2,
		"error_msg =?":     "",
	}).Order("program_start_time desc").Scan(&programTask)
	if err != nil {
		return
	}
	status := replayservice.ManifestTaskCrawlErr
	ApiUrl := configservice.GetString("replay.callback_url")
	//g.Dump(ApiUrl, programTask)
	res := g.Client().PostContent(gctx.GetInitCtx(), ApiUrl, g.Map{"info": programTask})
	//g.Dump(res)
	var replayRes *replayservice.ReplayRes
	err = g.NewVar(res).Structs(&replayRes)
	var isPostSucess bool = false
	var uuid string = ""
	//g.Dump(replayRes, err)
	if err == nil {
		if replayRes != nil {
			if replayRes.Code == 0 {
				if replayRes.Data.Code == 200 {
					var ReplayUUID replayservice.ReplayUUID
					er2 := g.NewVar(replayRes.Data.Data).Struct(&ReplayUUID)
					if er2 == nil {
						uuid = ReplayUUID.UUID
						isPostSucess = true
					}
					//uuid=replayRes.Data.Data.Uuid
				}
				if replayRes.Data.Code == 500 {
					uuidRe, e := gregex.MatchString(`Duplicate entry \'(.*)\' for key`, res)
					if e == nil {
						if len(uuidRe) > 1 {
							uuid = uuidRe[1]
							isPostSucess = true
						}

					}
				}
			}
		}
	}
	if isPostSucess {
		status = replayservice.ManifestTaskCrawlPostSuccess
		cdn := configservice.GetString("replay.play_cdn")
		programTask.ProgramPlayUrl = fmt.Sprintf("mtv://%s/%s", cdn, uuid)
	}
	replayservice.UpdateProgramTaskStatus(programTask, status)

}

func (receiver *CrawlReplayProgram) DelBeForDay(day int) {
	delDay := gtime.Now().Add(-time.Hour * 24 * g.NewVar(day).Duration()).Format("Y-m-d 01:00:00")
	path := configservice.GetString("dfs.rootPath")
	S := gfile.Separator
	path_rp := fmt.Sprintf("%s%sreplay", path, S)
	dir, _ := gfile.ScanDir(path_rp, "*", false)
	for _, v := range dir {
		fileDate := gtime.New(gfile.Basename(v))
		if gtime.New(delDay).Unix() > fileDate.Unix() {
			delPath := fmt.Sprintf("%s%sreplay%s%s", path, S, S, gfile.Basename(v))
			//g.Dump(delPath)
			gfile.Remove(delPath)
			g.Log().Infof(gctx.GetInitCtx(), "删除文件夹 del dir:%v", delPath)
		}
	}
}

// 创建定时任务列表-当前需要修复录制的(状态为抓取当中 且 Mp4文件大小3秒内无变化)
func (receiver *CrawlReplayProgram) CreateAllFixRecordingTask(context context.Context) {
	programTaskList := replayservice.GetAllByIpStatus(configservice.GetCrawlHostIp(), replayservice.ProgramTaskCrawling)
	for _, programTask := range programTaskList {

		go func(programTask *entity.CmsCrawlReplayProgramTask) {
			taskDO := replayservice.GetDetailById(programTask.Id)
			savePath := videoservice.GetReplayDir(taskDO)
			filePath := ffmpegutil.GetGenericFilePath(savePath)
			if isFixedSize(filePath) {
				taskStartTime := gtime.Now().Add(gtime.S * 3) // 给一个3秒延迟，防止pattern表达式过期
				doCreate(receiver, programTask, taskStartTime, true)
			}

		}(programTask)

	}
}

func (receiver *CrawlReplayProgram) StartRecording(ctx context.Context) {
	//programTask := replay.GetById(receiver.programTask.Id)
	programTask := receiver.programTask
	// 关闭数据库状态的判断和更新状态，加快调试速度
	disableDB := configservice.GetCrawlDebugBool("disableDB")
	if programTask.CrawlStatus != replayservice.ProgramTaskInit && !disableDB {
		// 在等待定时任务执行这个过程中，被别的工作节点抢先执行了
		return
	}
	// IP锁定后，后续正常流程后只会用这台机器去操作
	//programTask.HostIp = config.GetCrawlHostIp()

	// 确保多个工作节点不会录制同一个节目
	if !disableDB {
		replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskCrawling)
	}

	taskDO := replayservice.GetDetailById(programTask.Id)
	recordingSeconds := (taskDO.ProgramEndTime.UnixMilli() - taskDO.ProgramStartTime.UnixMilli()) / 1000
	g.Log().Infof(gctx.GetInitCtx(), "recordTime:%v", recordingSeconds)
	if disableDB {
		// 调试模式只录几秒
		recordingSeconds = 5
	}
	savePath := videoservice.GetReplayDir(taskDO)
	serverUrl := videoservice.GetServerUrl(taskDO)
	receiver.programTask.ProgramFilePath = savePath
	receiver.programTask.ProgramServerUrl = serverUrl
	_ = gfile.Mkdir(savePath)
	if recordingSeconds > 0 {
		//开始录制
		err := receiver.DoRecord(taskDO.PlayUrl, savePath, recordingSeconds)
		if err != nil {
			return
		}
		//录制成功后处理
		successreceiver := &CrawlReplayProgram{programTask: receiver.programTask}
		successreceiver.SuccessRecord()

		//receiver.SuccessRecord(savePath)

		//err := ffmpegutil.RunFfmpegGenericRecording(taskDO.PlayUrl, savePath, recordingSeconds)
		//if err != nil {
		//	g.Log().error(gctx.GetInitCtx(), "录制出错", err)
		//	if !disableDB {
		//		programTask.ErrorMsg = err.Error()
		//		replay.UpdateProgramTaskStatus(programTask, replay.ProgramTaskCrawlErr)
		//	}
		//} else {
		//	replay.UpdateProgramTaskStatus(programTask, replay.ProgramTaskCrawlFinish)
		//}
	} else {
		g.Log().Infof(gctx.GetInitCtx(), "已无需录制: programId = %d, programName = %s", programTask.Id, programTask.ProgramName)
		replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskCrawlFinish)
	}
}
func (receiver *CrawlReplayProgram) CheckManifestTaskCrawlFinish(context context.Context) {
	r := &CrawlReplayProgram{}
	r.SuccessRecord()
}
func (receiver *CrawlReplayProgram) SuccessRecord() {
	g.Log().Line().Debugf(gctx.GetInitCtx(), "开始转码检查 tanser check->%v", receiver.programTask)
	//超过最大转码数量则返回
	if !replayservice.NeedToTran(replayservice.ManifestTaskCrawlTraning) {
		g.Log().Line().Debugf(gctx.GetInitCtx(), "转码数量过多等待下次检查 tanser to more")
		return
	}
	//g.Dump(11)
	//查询最近一条录制完成的记录
	var err error
	//g.Dump(receiver.programTask)
	if receiver.programTask == nil {
		receiver.programTask, err = replayservice.GetOneByStatus(replayservice.ManifestTaskCrawlFinish)
		g.Log().Line().Debugf(gctx.GetInitCtx(), "执行定时检查切片任务 gron start")
	} else {
		g.Log().Line().Debugf(gctx.GetInitCtx(), "录制完成时转码 tanser start")
	}
	if err != nil || receiver.programTask == nil {
		g.Log().Line().Debugf(gctx.GetInitCtx(), "无等待转码任务 tanser not")
		return
	}

	//转换为正在转码
	replayservice.UpdateProgramTaskStatus(receiver.programTask, replayservice.ManifestTaskCrawlTraning)

	//开始转码
	g.Log().Line().Infof(gctx.GetInitCtx(), "开始执行转码 tanser start->%v", g.NewVar(receiver.programTask).Map())
	savePath := receiver.programTask.ProgramFilePath
	saveFile := fmt.Sprintf("%s%s%s", savePath, gfile.Separator, "org.ts")
	endFile := fmt.Sprintf("%s%s%s", savePath, gfile.Separator, "org.mp4")
	programTask := receiver.programTask
	ffmpegObj := ffmpegutil.FmpegTrans("ffmpeg")
	err = ffmpegObj.CheckFile(savePath, saveFile, endFile)
	if err != nil {
		programTask.ErrorMsg = err.Error()
		replayservice.UpdateProgramTaskStatus(programTask, replayservice.ManifestTaskCrawlErr)
		return
	}
	//转码完毕

	m3u8File := fmt.Sprintf("%s%s%s", savePath, gfile.Separator, "index.m3u8")
	receiver.programTask.ProgramFilePath = m3u8File
	needToPost := configservice.GetBool("replay.need_post")
	if needToPost {
		status := replayservice.ManifestTaskCrawlErr
		ApiUrl := configservice.GetString("replay.callback_url")
		res := g.Client().PostContent(gctx.GetInitCtx(), ApiUrl, g.Map{"info": receiver.programTask})
		//g.Dump(res)
		var replayRes *replayservice.ReplayRes
		err = g.NewVar(res).Structs(&replayRes)
		var isPostSucess bool = false
		var uuid string = ""
		//g.Dump(replayRes, err)
		if err == nil {
			if replayRes.Code == 0 {
				if replayRes.Data.Code == 200 {
					var ReplayUUID replayservice.ReplayUUID
					er2 := g.NewVar(replayRes.Data.Data).Struct(&ReplayUUID)
					if er2 == nil {
						uuid = ReplayUUID.UUID
						isPostSucess = true
					}
					//uuid=replayRes.Data.Data.Uuid
				}

			}
		}
		if isPostSucess {
			status = replayservice.ManifestTaskCrawlPostSuccess
			cdn := configservice.GetString("replay.play_cdn")
			receiver.programTask.ProgramPlayUrl = fmt.Sprintf("mtv://%s/%s", cdn, uuid)
		}
		replayservice.UpdateProgramTaskStatus(receiver.programTask, status)
	} else {
		replayservice.UpdateProgramTaskStatus(receiver.programTask, replayservice.ManifestTaskCrawlPostSuccess)
	}

}
func (receiver *CrawlReplayProgram) DoRecord(PlayUrl, savePath string, recordingSeconds int64) (err error) {
	//发起HTTP握手动作
	saveFile := fmt.Sprintf("%s%s%s", savePath, gfile.Separator, "org.ts")
	//endFile := fmt.Sprintf("%s%s%s", savePath, gfile.Separator, "org.mp4")
	g.Log().Infof(gctx.GetInitCtx(), "saveFile:%v", saveFile)
	request, err := http.NewRequest("GET", PlayUrl, nil)
	if err != nil {
		g.Log().Infof(gctx.GetInitCtx(), "err1:%v", err)
		return
	}
	client := &http.Client{}
	//建立HTTP请求 获取返回内容
	response, err := client.Do(request)
	if err != nil {
		g.Log().Infof(gctx.GetInitCtx(), "err2:%v", err)
		return
	}
	file, _ := os.Create(saveFile)
	//方法结束后关闭HTPP链接 关闭文件链接
	defer response.Body.Close()
	//成功建立连接后 因为直播流会不停推送数据所以挂起死循环等待连接的返回
	var isOver bool = false
	go func() {
		//录制时间定义 例子为10秒
		<-time.After(time.Second * g.NewVar(recordingSeconds).Duration())
		isOver = true
	}()
	for {
		if isOver {
			//停止录制后需要用ffmpeg -i 录制文件 -c copy 输出文件.mp4 转换视频封装格式才能观看（因为其他来源视频封装不一定统一 只转封装格式不转码 1秒就完事了）
			break
		}
		//新建4096字节缓冲区 （因为一次http请求最大返回4096）
		buf := make([]byte, 4096)
		n, ee := response.Body.Read(buf)
		if ee != nil {
			//	源掉线或者本地网络不好 这种原因断线之后的处理
			break
		}
		//获取到的实际大小buf[:n] 实际内容数据写入文件
		//可以每个buf[:n]直接go里处理好二进制文件 但是实现起来比较复杂 有参考代码。但是看不懂
		file.Write(buf[:n])
	}
	file.Close()
	receiver.programTask.ProgramFilePath = savePath
	replayservice.UpdateProgramTaskStatus(receiver.programTask, replayservice.ManifestTaskCrawlFinish)

	//programTask := receiver.programTask
	//ffmpeg := ffmpegutil.FmpegTrans("ffmpeg")
	//err = ffmpegutil.CheckFile(savePath, saveFile, endFile)
	//if err != nil {
	//	programTask.ErrorMsg = err.Error()
	//	replay.UpdateProgramTaskStatus(programTask, replay.ManifestTaskCrawlErr)
	//	return
	//}
	return

}

func (receiver *CrawlReplayProgram) StartFixRecording(ctx context.Context) {
	programTask := replayservice.GetById(receiver.programTask.Id)

	taskDO := replayservice.GetDetailById(programTask.Id)
	recordingSeconds := (taskDO.ProgramEndTime.UnixMilli() - receiver.startRecordingTime.UnixMilli()) / 1000

	if recordingSeconds > 0 {
		savePath := videoservice.GetReplayDir(taskDO)
		_ = gfile.Mkdir(savePath)
		filePath := ffmpegutil.GetGenericFilePath(savePath)

		// 之前录了上半段,需要合并的
		needMerge := gfile.Exists(filePath)

		// 断开多少次就有多少orgFiles
		orgFiles, _ := gfile.ScanDir(savePath, "org_*.mp4")
		g.Log().Info(gctx.GetInitCtx(), "等待合并片段数：%d", len(orgFiles))
		waitMergeMp4 := ffmpegutil.GetFilePath(savePath, fmt.Sprintf("org_%d.mp4", len(orgFiles)))
		if needMerge {
			_ = gfile.Rename(filePath, waitMergeMp4)
		}

		err := ffmpegutil.RunFfmpegGenericRecording(taskDO.PlayUrl, savePath, recordingSeconds)
		if err != nil {
			g.Log().Error(gctx.GetInitCtx(), "录制出错", err)
			programTask.ErrorMsg = err.Error()
			replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskCrawlErr)
		} else {
			if needMerge {
				err = ffmpegutil.RunFfmpegGenericMerge(savePath)
				g.Log().Error(gctx.GetInitCtx(), "视频合并出错", err)
			}
			replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskCrawlFinish)
		}
	} else {
		g.Log().Infof(gctx.GetInitCtx(), "已无需修复: programId = %d, programName = %s", programTask.Id, programTask.ProgramName)
		replayservice.UpdateProgramTaskStatus(programTask, replayservice.ProgramTaskCrawlFinish)
	}
}

// 是否是固定大小（即：3秒后体积不变）
func isFixedSize(filePath string) bool {
	beforeSize := gfile.Size(filePath)
	time.Sleep(gtime.S * 3)
	afterSize := gfile.Size(filePath)
	return afterSize == beforeSize
}
