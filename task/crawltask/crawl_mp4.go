package crawltask

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-frame/db/mysql/dao"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/configservice"
	"go-to-crawl-frame/service/crawl/sysservice"
	"go-to-crawl-frame/service/crawl/uploadservice"
	"go-to-crawl-frame/service/crawl/videoservice"
	"go-to-crawl-frame/service/crawl/vodservice"
	"go-to-crawl-frame/service/lockservice"
	"go-to-crawl-frame/utils/ffmpegutil"
	"go-to-crawl-frame/utils/fileutil"
)

func DownloadMp4Type1Task(context context.Context) {
	doDownloadMp4(vodservice.HostTypeCrawlLogin)
}

func DownloadMp4Type2Task(context context.Context) {
	doDownloadMp4(vodservice.HostTypeNiVod)
}

func DownloadMp4Type3Task(context context.Context) {
	if lockservice.IncreaseValue(lockservice.DownloadMp4Type3) {
		defer lockservice.DecreaseValue(lockservice.DownloadMp4Type3)
		doDownloadMp4(vodservice.HostTypeBananTV)
	}
}

func DownloadMp4Task(context context.Context) {
	if !lockservice.TryLockSelenium() {
		return
	}
	defer lockservice.ReleaseLockSelenium()
	doDownloadMp4(vodservice.HostTypeNormal)
}

func doDownloadMp4(hostType int) {
	seed := vodservice.GetSeed(vodservice.CrawlFinish, configservice.GetCrawlHostIp(), hostType)

	if seed == nil {
		return
	}
	log := g.Log().Line()
	log.Infof(gctx.GetInitCtx(),
		"doDownloadMp4: seed = %v", gjson.New(seed).MustToJsonString())

	vodservice.UpdateStatus(seed, vodservice.M3U8Parsing)
	// 创建最终目录
	videoDir := videoservice.GetVideoDir(seed.CountryCode, seed.VideoYear, seed.VideoCollId, seed.VideoItemId)
	_ = gfile.Mkdir(videoDir)

	// 下载完M3U8后，后续操作都只能当前主机处理
	seed.HostIp = configservice.GetCrawlHostIp()

	if vodservice.TypeMP4Url == seed.CrawlType {
		var downloadErr error
		// 直接下载MP4
		if vodservice.IsVideoAudioSeparated(seed.CrawlSeedParams) {
			// 音视频分离场景
			downloadErr = doDownloadMp4Separated(seed, videoDir)
		} else {
			// 音视频未分离场景
			downloadErr = doDownloadMp4Direct(seed, videoDir)
		}

		if downloadErr != nil {
			videoservice.UpdateDownloadStatus(seed, errors.New("MP4下载失败"))
			return
		} else {
			videoservice.UpdateDownloadStatus(seed, nil)
		}

	} else {

		// 下载M3U8文件
		orgM3U8File := videoDir + ffmpegutil.OrgM3U8Name
		proxyUrl := sysservice.GetProxyByUrl(seed.CrawlM3U8Url)

		err := fileutil.DownloadM3U8File(seed.CrawlM3U8Url, proxyUrl, orgM3U8File, fileutil.Retry, seed.CrawlM3U8Text)
		if err != nil {
			log.Info(gctx.GetInitCtx(), err)
			seed.ErrorMsg = "Download M3U8 Error"
			vodservice.UpdateStatus(seed, vodservice.M3U8Err)
			return
		}

		strategy := getCrawlVodFlowStrategy(seed)
		m3u8DO, err2 := strategy.ConvertM3U8(seed, orgM3U8File)
		if err2 != nil {
			log.Info(gctx.GetInitCtx(), err2)
			seed.ErrorMsg = "标准化M3U8文件出错"
			vodservice.UpdateStatus(seed, vodservice.M3U8Err)
			return
		}

		m3u8DO.NotDownloadAll = configservice.GetCrawlDebugBool("notDownloadAll")
		m3u8DO.FromUrlProxy = sysservice.GetProxyByUrl(m3u8DO.FromUrl)
		err2 = strategy.DownLoadToMp4(m3u8DO)

		if err2 != nil {
			videoservice.UpdateDownloadStatus(seed, errors.New("M3U8转MP4出错"))
			return
		} else {
			videoservice.UpdateDownloadStatus(seed, nil)
		}
		//更新成功後刪除原m3u8文件
		_ = gfile.Remove(orgM3U8File)

		if vodservice.HostTypeCrawlLogin == hostType {
			// 国内指定机器下载的，需要上传到国外点播服务器
			//file.UpLoadToFastDFS(m3u8DO.MP4File, seed)
		}
	}

	// 添加到转换队列
	upQueue := new(entity.CmsUploadQueue)
	_ = gconv.Struct(seed, upQueue)
	upQueue.Id = 0
	upQueue.FileName = ffmpegutil.OrgMp4Name
	upQueue.UploadStatus = uploadservice.Uploaded
	upQueue.CreateTime = gtime.Now()
	_, _ = dao.CmsUploadQueue.Ctx(gctx.GetInitCtx()).Insert(upQueue)

}

func doDownloadMp4Direct(seed *entity.CmsCrawlQueue, videoDir string) error {
	builder := fileutil.CreateBuilder()
	builder.Url(getMp4Url(seed))
	builder.SaveFile(fmt.Sprintf("%s%s", videoDir, ffmpegutil.OrgMp4Name))
	useSeedParams(seed, builder)
	downloadErr := fileutil.DownloadFileByBuilder(builder)
	return downloadErr
}

func doDownloadMp4Separated(seed *entity.CmsCrawlQueue, videoDir string) error {
	params := new(vodservice.CrawlSeedParams)
	params.InitSeedParams(seed.CrawlSeedParams)

	// 下载视频
	videoFile := fmt.Sprintf("%s%s", videoDir, ffmpegutil.OrgVideoName)
	videoBuilder := fileutil.CreateBuilder().Url(params.VideoUrl).SaveFile(videoFile)
	useSeedParams(seed, videoBuilder)
	videoErr := fileutil.DownloadFileByBuilder(videoBuilder)

	// 下载音频
	audioFile := fmt.Sprintf("%s%s", videoDir, ffmpegutil.OrgAudioName)
	audioBuilder := fileutil.CreateBuilder().Url(params.AudioUrl).SaveFile(audioFile)
	useSeedParams(seed, audioBuilder)
	audioErr := fileutil.DownloadFileByBuilder(audioBuilder)

	if videoErr != nil || audioErr != nil {
		return errors.New("视频或者音频下载失败")
	} else {
		mergeErr := ffmpegutil.MergeVideoAudio(videoDir, videoFile, audioFile)
		if mergeErr != nil {
			return mergeErr
		}
		_ = gfile.Remove(videoFile)
		_ = gfile.Remove(audioFile)
		return nil
	}
}

func useSeedParams(seed *entity.CmsCrawlQueue, builder *fileutil.DownloadBuilder) {
	if seed.CrawlSeedParams != "" {
		params := new(vodservice.CrawlSeedParams)
		params.InitSeedParams(seed.CrawlSeedParams)
		if params.M3u8UrlReqHeader != nil {
			builder.Headers(params.M3u8UrlReqHeader)
		}
	}
}

// CrawlType=3可能是直接录入的MP4地址导seedUrl字段，也可能是根据seedUrl抓到的M3U8Url字段(如:B站)
func getMp4Url(seed *entity.CmsCrawlQueue) string {
	if gstr.LenRune(seed.CrawlM3U8Url) != 0 {
		return seed.CrawlM3U8Url
	}
	return seed.CrawlSeedUrl
}
