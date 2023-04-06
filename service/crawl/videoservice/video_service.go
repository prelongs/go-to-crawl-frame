package videoservice

import (
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/db/mysql/model/entity"
	"github.com/JervisPG/go-to-crawl-frame/service/configservice"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/servicedto"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/vodservice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

var (
	s           = gfile.Separator
	dfsRootPath = configservice.GetString("dfs.rootPath")
)

func GetVideoDir(countryCode string, videoYear int, videoCollId, videoItemId int64) string {
	videoRootPath := fmt.Sprintf("%s%s%s%s", dfsRootPath, s, "video", s)
	finalFilePath := fmt.Sprintf("%s%s%s%d%s%d%s%d%s", videoRootPath, countryCode, s, videoYear, s, videoCollId, s, videoItemId, s)
	return finalFilePath
}

func GetReplayDir(taskDO *servicedto.CmsCrawlReplayProgramTaskDO) string {
	replayRootPath := fmt.Sprintf("%s%s%s", dfsRootPath, s, "replay")

	channelPath := fmt.Sprintf("%s%s%s%s%s", replayRootPath, s, gtime.New(taskDO.ProgramStartTime).Format("Y-m-d"), s, g.NewVar(taskDO.Id).String())
	fileDir := fmt.Sprintf("%s%s%s", channelPath, s, taskDO.ProgramNo)
	return fileDir
}

func GetLiveDir(programName string) string {
	return fmt.Sprintf("%s%s%s%s%s", dfsRootPath, s, "live", s, programName)
}

func GetLiveM3U8CacheFile(programName string) string {
	return fmt.Sprintf("%s%s%s", GetLiveDir(programName), gfile.Separator, "index.m3u8.cache")
}

func GetServerUrl(taskDO *servicedto.CmsCrawlReplayProgramTaskDO) string {
	serverUrl := configservice.GetString("replay.play_domain")
	return fmt.Sprintf("http://%s/%s/%s/%s/index.m3u8", serverUrl, gtime.New(taskDO.ProgramStartTime).Format("Y-m-d"), g.NewVar(taskDO.Id).String(), taskDO.ProgramNo)
}

func UpdateDownloadStatus(seed *entity.CmsCrawlQueue, err error) {
	if err != nil {
		g.Log().Line().Error(gctx.GetInitCtx(), err)
		seed.ErrorMsg = err.Error()
		vodservice.UpdateStatus(seed, vodservice.M3U8Err)
	} else {
		vodservice.UpdateStatus(seed, vodservice.M3U8Parsed)
	}
}
