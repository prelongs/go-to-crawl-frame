package ffmpegutil

import (
	"errors"
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/utils/fileutil"
	"github.com/JervisPG/go-to-crawl-frame/utils/httputil"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	ProtocolWhitelist = "concat,file,http,https,tcp,tls,crypto"
	OrgMp4Name        = "org.mp4"
	OrgVideoName      = "video.mp4"
	OrgAudioName      = "audio.wav"
	OrgTmpMp4Name     = "org_tmp.mp4" // 合并后临时的MP4
	OrgM3U8Name       = "org_index.m3u8"
	KeyLine           = "METHOD="
	KeyName           = "key.m3u8"
	ExtMapLine        = "EXT-X-MAP"
	ExtinfoLine       = "EXTINF"
	EndLine           = "EXT-X-ENDLIST"
)

const (
	LineTypeComment = 1 // 注解行
	LineTypeKey     = 2 // 解密KEY行
	LineTypeSrc     = 3 // TS资源行
	LineTypeXMap    = 4 // #EXT-X-MAP行（所有ts统一的文件头，压缩体积所需）
	LineTypeExtInfo = 5 // #EXTINF行（当前TS播放时间）
	LineTypeEnd     = 6 // #EXT-X-ENDLIST行
)

const (
	SrcTypeNormal   = 1 // 普通资源类型 eg: .ts后缀
	SrcTypeImg      = 2 // 图片资源类型 eg: .jpg后缀
	SrcTypeNoSuffix = 3 // 无后缀 eg: http://xxx.com/abc
)

const (
	tsFilePattern     = "*.ts"
	orgMp4FilePattern = "org_*.mp4"
)

func DownloadToMp4(m3u8DO *M3u8DO) error {
	log := g.Log().Line()
	log.Infof(gctx.GetInitCtx(),
		"final path : %s", m3u8DO.FromDir)

	DiscardTsWhenDebug(m3u8DO)
	err := downloadConcurrent(m3u8DO, 10)
	if err != nil {
		return err
	}

	transToLocalM3U8(m3u8DO)
	err = MergeTsFile(m3u8DO)
	if err != nil {
		return err
	}
	DeleteTmpResource(m3u8DO, tsFilePattern)

	return nil
}

func DownloadToLiveStream(m3u8DO *M3u8DO) error {
	log := g.Log().Line()
	log.Infof(gctx.GetInitCtx(),
		"final path : %s", m3u8DO.FromDir)

	m3u8CacheFile := m3u8DO.FromFile
	m3u8File := gstr.Replace(m3u8CacheFile, ".cache", "")

	m3u8ContentOld := gfile.GetContents(m3u8File)
	m3u8ContentNew := getLocalM3U8(m3u8DO, false, false)

	if gfile.Exists(m3u8File) {
		// 追加还未下载的TS为本地模式且添加到原有index.m3u8末尾
		m3u8ContentTail := getLocalM3U8Body(m3u8DO, false, m3u8ContentOld)
		_ = gfile.PutContentsAppend(m3u8File, m3u8ContentTail)
	} else {
		// 首次执行，index.m3u8还不存在
		_ = gfile.PutContentsAppend(m3u8File, m3u8ContentNew)
	}

	// 考虑本地下载速度和目标网的限速，线程数要支持可配置(暂定为2看效果)
	err := downloadConcurrent(m3u8DO, 2)
	if err != nil {
		return err
	}

	// 重写index.m3u8.cache
	_ = gfile.Remove(m3u8CacheFile)
	_ = gfile.PutContentsAppend(m3u8CacheFile, m3u8ContentNew)

	// cache转为非cache对外提供服务
	_ = os.Rename(m3u8CacheFile, m3u8File)

	// 删除无需对外提供服务的临时TS资源
	DeleteLiveTmpResource(m3u8DO, tsFilePattern)

	return nil
}

func downloadConcurrent(m3u8DO *M3u8DO, poolSize int) error {

	proxyUrl := m3u8DO.FromUrlProxy
	err := DownloadDependenceFile(m3u8DO, proxyUrl)
	if err != nil {
		return err
	}

	// 开启多线程下载
	maxChan := make(chan bool, poolSize)
	var failCount int64 = 0
	wg := sync.WaitGroup{}
	for _, m3u8LineDO := range m3u8DO.StreamLineList {
		wg.Add(1)
		if atomic.LoadInt64(&failCount) > 0 {
			return errors.New("")
		}
		maxChan <- true
		go func(lineDO StreamLineDO) {
			if lineDO.LineType != LineTypeSrc {
				<-maxChan
				wg.Done()
				return
			}

			err2 := downloadTsLine(m3u8DO, lineDO, proxyUrl)

			<-maxChan
			wg.Done()
			if err2 != nil {
				atomic.AddInt64(&failCount, 1)
				return
			}
		}(m3u8LineDO)

	}
	wg.Wait()

	return nil
}

func DiscardTsWhenDebug(m3u8DO *M3u8DO) {
	if m3u8DO.NotDownloadAll {
		srcCnt := 0
		var newLineList []StreamLineDO
		for _, lineDO := range m3u8DO.StreamLineList {
			newLineList = append(newLineList, lineDO)
			if lineDO.LineType == LineTypeSrc {
				srcCnt += 1
			}
			if srcCnt == 20 {
				break
			}
		}

		// 保存结束行标志，否则FFMPEG会一直等待结束而卡住
		newLineList = append(newLineList, m3u8DO.StreamLineList[len(m3u8DO.StreamLineList)-1])
		m3u8DO.StreamLineList = newLineList
	}
}

func ConvertM3U8(m3u8Url, baseUrl, filePath string) (*M3u8DO, error) {
	return doConvertM3U8(m3u8Url, baseUrl, filePath, false)
}

func ConvertM3U8WithOriginTsName(m3u8Url, baseUrl, filePath string) (*M3u8DO, error) {
	return doConvertM3U8(m3u8Url, baseUrl, filePath, true)
}

// 解析M3U8文件为对象
func doConvertM3U8(m3u8Url, baseUrl, filePath string, withOriginTsName bool) (*M3u8DO, error) {

	log := g.Log().Line()
	log.Info(gctx.GetInitCtx(), "标准化M3U8：", filePath)

	m3u8DO := new(M3u8DO)
	m3u8DO.FromUrl = m3u8Url
	m3u8DO.FromBaseUrl = baseUrl
	m3u8DO.FromFile = filePath
	m3u8DO.FromDir = gfile.Dir(filePath)

	contentNew := ""
	idx := 0
	_ = gfile.ReadLines(filePath, func(line string) error {

		if gstr.LenRune(line) == 0 {
			// 过滤空行
			return nil
		}

		m3u8LineDO := new(StreamLineDO)
		m3u8LineDO.LineType = -1
		m3u8LineDO.SrcType = -1
		m3u8LineDO.OriginLine = line
		m3u8LineDO.TransformedLine = line

		if gstr.HasPrefix(line, "#") {
			if gstr.ContainsI(line, ExtMapLine) {
				m3u8LineDO.LineType = LineTypeXMap
			} else if gstr.ContainsI(line, KeyLine) {
				m3u8LineDO.LineType = LineTypeKey
				m3u8LineDO.TransformedLine = paddingKeyUrl(baseUrl, line)
			} else if gstr.ContainsI(line, ExtinfoLine) {
				m3u8LineDO.LineType = LineTypeExtInfo
				playSecondMatch, _ := gregex.MatchString("\\d+\\.?\\d*", line)
				m3u8LineDO.PlaySecond = gconv.Float32(playSecondMatch[0])
				m3u8DO.PlaySecondTotal += m3u8LineDO.PlaySecond
			} else if gstr.ContainsI(line, EndLine) {
				m3u8LineDO.LineType = LineTypeEnd
			} else {
				m3u8LineDO.LineType = LineTypeComment
			}
		} else {
			m3u8LineDO.LineType = LineTypeSrc
			nameIdx := strings.LastIndex(line, "/")
			orgName := gstr.SubStr(line, nameIdx+1)
			m3u8LineDO.OriginTsSrcName = orgName

			if withOriginTsName {
				dotIdx := strings.LastIndex(orgName, ".")
				orgNameNoExt := gstr.SubStr(orgName, 0, dotIdx)
				m3u8LineDO.TransformedSrcName = fmt.Sprintf("%s.ts", orgNameNoExt)
			} else {
				m3u8LineDO.TransformedSrcName = fmt.Sprintf("%d.ts", idx)
			}
			idx += 1
			m3u8DO.SrcLineCount += 1

			if gstr.ContainsI(line, ".ts") {
				m3u8LineDO.SrcType = SrcTypeNormal
			} else if gstr.ContainsI(line, ".jpg") || gstr.ContainsI(line, ".png") {
				m3u8LineDO.SrcType = SrcTypeImg
			} else {
				m3u8LineDO.SrcType = SrcTypeNoSuffix
			}

			if !gstr.HasPrefix(line, "http") {
				if gstr.HasPrefix(line, "/") {
					m3u8LineDO.TransformedLine = baseUrl + line
				} else {
					m3u8LineDO.TransformedLine = httputil.GetBaseUrlByBackslash(m3u8Url) + line
				}
			}

		}

		m3u8DO.StreamLineList = append(m3u8DO.StreamLineList, *m3u8LineDO)
		contentNew += m3u8LineDO.TransformedLine + "\n"
		return nil
	})

	_ = gfile.Remove(filePath)
	_ = ioutil.WriteFile(filePath, []byte(contentNew), fs.ModeAppend)

	return m3u8DO, nil
}

func transToLocalM3U8(m3u8DO *M3u8DO) {
	contentNew := getLocalM3U8(m3u8DO, true, true)
	// 重写m3u8文件为本地模式
	_ = gfile.Remove(m3u8DO.FromFile)
	_ = ioutil.WriteFile(m3u8DO.FromFile, []byte(contentNew), fs.ModeAppend)
	_ = os.Chmod(m3u8DO.FromFile, 0666)
}

func getLocalM3U8(m3u8DO *M3u8DO, srcNameWithParentPath bool, withEndLine bool) string {
	content := ""
	// 转化为新的TS内容
	content += getLocalM3U8Head(m3u8DO)
	content += getLocalM3U8Body(m3u8DO, srcNameWithParentPath, "")
	if withEndLine {
		content += "#" + EndLine
	}
	return content
}

func getLocalM3U8Head(m3u8DO *M3u8DO) string {
	contentTmp := ""
	for _, lineDO := range m3u8DO.StreamLineList {
		if lineDO.LineType == LineTypeExtInfo {
			break
		}
		if lineDO.LineType == LineTypeXMap {
			// 预留
		} else if lineDO.LineType == LineTypeKey {
			keyLine := gstr.Replace(lineDO.OriginLine, m3u8DO.KeyOriginName, KeyName)
			contentTmp += keyLine
		} else {
			contentTmp += lineDO.TransformedLine
		}
		contentTmp += "\n"
	}
	return contentTmp
}

func getLocalM3U8Body(m3u8DO *M3u8DO, srcNameWithParentPath bool, oldM3U8Content string) string {
	contentTmp := ""
	for idx := range m3u8DO.StreamLineList {
		if idx == 0 {
			continue
		}

		lineExtInfo := m3u8DO.StreamLineList[idx-1]
		lineCurrent := m3u8DO.StreamLineList[idx]
		if lineCurrent.LineType == LineTypeSrc {
			srcLine := getLocalM3U8TypeSrc(m3u8DO, lineCurrent, srcNameWithParentPath)
			// 去重方案
			if oldM3U8Content == "" {
				contentTmp += lineExtInfo.TransformedLine + "\n"
				contentTmp += srcLine + "\n"
			} else {
				if !gstr.ContainsI(oldM3U8Content, srcLine) {
					contentTmp += lineExtInfo.TransformedLine + "\n"
					contentTmp += srcLine + "\n"
				}
			}
		}
	}
	return contentTmp
}

func getLocalM3U8TypeSrc(m3u8DO *M3u8DO, lineDO StreamLineDO, srcNameWithParentPath bool) string {

	if srcNameWithParentPath {
		// TS合成MP4的时候TS需要全路径，MP4切为最终TS的时候M3U8里的TS无需全路径
		return m3u8DO.FromDir + gfile.Separator + lineDO.TransformedSrcName
	} else {
		// 直播的时候直接就是最终的TS名称
		return lineDO.TransformedSrcName
	}
}

func downloadTsLine(m3u8DO *M3u8DO, lineDO StreamLineDO, proxyUrl string) error {
	// 随机等几百毫秒，一定程度防止把对方服务弄垮，也防止把自己机器CPU跑太高
	rand := grand.Intn(1000)
	time.Sleep(time.Nanosecond * time.Duration(rand))
	tsFilePath := m3u8DO.FromDir + gfile.Separator + lineDO.TransformedSrcName

	if gfile.Exists(tsFilePath) && gfile.Size(tsFilePath) > 100*1024 {
		// 如果TS片段存在 且 大于100K就认定为是已经存在的正确资源（eg：偶尔网络或者授权问题，下载的ts只有1K，那么下次就能继续下载）
		return nil
	}

	builder := fileutil.CreateBuilder().Url(lineDO.TransformedLine).Proxy(proxyUrl).SaveFile(tsFilePath).Retry(fileutil.Retry)
	if m3u8DO.Headers != nil {
		builder.Headers(*m3u8DO.Headers)
	}

	err := fileutil.DownloadFileByBuilder(builder)
	truncateTS(tsFilePath, lineDO.SrcType, m3u8DO.PngHeaderSize)
	return err
}

func formatExtName(line string, imgMode *bool) string {
	if gstr.ContainsI(line, ".jpg") {
		line = gstr.ReplaceI(line, ".jpg", ".ts")
		*imgMode = true
	}

	// eg：https://dweb.link/ipfs/bafybeihzekyi4nrinzktgwthvhrv47stbeylvyprdjlicngf65hssgj2n4
	if gstr.HasPrefix(line, "http") && !gstr.ContainsI(line, ".ts") {
		line += ".ts"
	}
	return line
}

func paddingKeyUrl(baseUrl string, line string) string {
	log := g.Log().Line()
	if gstr.Contains(line, KeyLine) {
		log.Info(gctx.GetInitCtx(), "填充前key Url = ", line)
		ret, err := gregex.ReplaceStringFunc("\".*\"", line, func(matchedStr string) string {
			url := gstr.Trim(matchedStr, "\"")
			if gstr.ContainsI(url, "http") {
				// 已经有Schema
				return matchedStr
			}
			// 没有Schema
			return fmt.Sprintf("\"%s%s\"", baseUrl, url)
		})

		if err == nil {
			line = ret
		}

		log.Info(gctx.GetInitCtx(), "填充后key Url = ", line)
	}
	return line
}

// 去掉PNG头
func truncateTS(tsPath string, srcType int, pngHeaderSize int64) {
	if srcType == SrcTypeImg || IsPngType(tsPath) {
		size := gfile.Size(tsPath)
		bytes := gfile.GetBytesByTwoOffsetsByPath(tsPath, pngHeaderSize, size)
		_ = ioutil.WriteFile(tsPath, bytes, fs.ModeAppend)
	}
}

func IsPngType(tsPath string) bool {
	bytes := gfile.GetBytesByTwoOffsetsByPath(tsPath, 0, 512)
	contentType := http.DetectContentType(bytes)
	return gstr.ContainsI(contentType, "png")
}

// 运行FFMPEG进行通用型切片(通用型定义：1、原视频名称为固定的；2、FFMPEG切片参数为固定的)
func RunFfmpegGenericSlice(basePath string) error {
	mp4File := GetGenericFilePath(basePath)
	segmentFile := fmt.Sprintf("%s%s%s", basePath, gfile.Separator, "index%3d.ts")
	m3u8File := fmt.Sprintf("%s%s%s", basePath, gfile.Separator, "index.m3u8")
	args := []string{"-i", mp4File, "-hls_time", "3", "-hls_list_size", "0", "-hls_segment_filename", segmentFile, m3u8File}
	return RunFfmpegCommand(args...)
}

// 运行FFMPEG进行通用型切片(通用型定义：1、原视频名称为固定的；2、FFMPEG切片参数为固定的)
func RunFfmpegGenericRecording(recordingUrl, basePath string, recordingSeconds int64) error {
	filePath := GetGenericFilePath(basePath)
	args := []string{"-t", strconv.Itoa(int(recordingSeconds)), "-v", "error", "-y", "-re", "-i", recordingUrl, "-c", "copy", filePath}
	return RunFfmpegCommand(args...)
}

// 运行FFMPEG进行通用合并（通用型定义：OrgWaitMp4Name拼接到OrgMp4Name的前面）
func RunFfmpegGenericMerge(basePath string) error {

	// 1、生成list_file内容
	listFileContent := ""
	orgFiles, _ := gfile.ScanDir(basePath, orgMp4FilePattern)
	g.Log().Infof(gctx.GetInitCtx(), "待合并文件数: %d", len(orgFiles))
	for idx := range orgFiles {
		listFileContent += fmt.Sprintf("file 'org_%d.mp4'\n", idx)
	}
	listFileContent += fmt.Sprintf("file '%s'\n", OrgMp4Name)

	// 2、list_file内容写入filelist.txt
	listFilePath := basePath + gfile.Separator + "filelist.txt"
	_ = gfile.Remove(listFilePath)
	_ = ioutil.WriteFile(listFilePath, []byte(listFileContent), fs.ModeAppend)

	// 3、执行合并
	orgMp4FilePath := GetGenericFilePath(basePath)
	orgTmpMp4FilePath := basePath + gfile.Separator + OrgTmpMp4Name
	args := []string{"-f", "concat", "-safe", "0", "-i", listFilePath, "-y", orgTmpMp4FilePath}
	err := RunFfmpegCommand(args...)
	if err != nil {
		g.Log().Errorf(gctx.GetInitCtx(), "合并出错：%s", err.Error())
		return err
	}

	// 4、删除临时文件
	_ = gfile.Remove(listFilePath)
	_ = gfile.Remove(orgMp4FilePath)
	for _, orgFile := range orgFiles {
		g.Log().Infof(gctx.GetInitCtx(), "删除mp4片段：%s", orgFile)
		_ = gfile.Remove(orgFile)
	}

	// 5、重命名为项目标准mp4命令文件
	_ = gfile.Rename(orgTmpMp4FilePath, orgMp4FilePath)

	return nil
}

func GetGenericFilePath(basePath string) string {
	return GetFilePath(basePath, OrgMp4Name)
}

func GetFilePath(basePath string, fileName string) string {
	return fmt.Sprintf("%s%s%s", basePath, gfile.Separator, fileName)
}

func RunFfmpegCommand(arg ...string) error {
	cmd := exec.Command("ffmpeg", arg...)
	g.Log().Info(gctx.GetInitCtx(), "执行命令: ", cmd.String())
	return cmd.Run()
}

func DownloadDependenceFile(m3u8DO *M3u8DO, proxyUrl string) error {
	err := downloadMapFile(m3u8DO, proxyUrl)
	if err != nil {
		return err
	}
	err = downloadKeyFile(m3u8DO, proxyUrl)
	if err != nil {
		return err
	}
	return nil
}

func downloadMapFile(m3u8DO *M3u8DO, proxyUrl string) error {
	for _, m3u8LineDO := range m3u8DO.StreamLineList {
		if m3u8LineDO.LineType == LineTypeXMap {
			err := downloadResourceFile(m3u8DO, m3u8LineDO, proxyUrl)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func downloadKeyFile(m3u8DO *M3u8DO, proxyUrl string) error {
	for _, m3u8LineDO := range m3u8DO.StreamLineList {
		if m3u8LineDO.LineType == LineTypeKey {
			err := downloadResourceFile(m3u8DO, m3u8LineDO, proxyUrl)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func downloadResourceFile(m3u8DO *M3u8DO, m3u8LineDO StreamLineDO, proxyUrl string) error {
	URIS, _ := gregex.MatchString(`URI="(.*)"`, m3u8LineDO.OriginLine)
	fileName := URIS[1]
	idx := strings.LastIndex(fileName, "/")
	fileNameSimple := gstr.SubStr(fileName, idx+1)

	downloadUrl := m3u8DO.FromBaseUrl + fileName
	saveFile := m3u8DO.FromDir + gfile.Separator + fileNameSimple

	builder := fileutil.CreateBuilder().Url(downloadUrl).Proxy(proxyUrl).SaveFile(saveFile).Retry(fileutil.Retry)
	if m3u8DO.Headers != nil {
		builder.Headers(*m3u8DO.Headers)
	}
	err := fileutil.DownloadFileByBuilder(builder)

	if m3u8LineDO.LineType == LineTypeKey {
		m3u8DO.KeyOriginName = fileName
		m3u8DO.KeySaveFile = m3u8DO.FromDir + gfile.Separator + KeyName
		_ = gfile.Rename(saveFile, m3u8DO.KeySaveFile)
	} else if m3u8LineDO.LineType == LineTypeXMap {
		m3u8DO.MapOriginName = fileName
		m3u8DO.MapSaveFile = saveFile
	}

	if err != nil {
		return err
	}
	return nil
}

func MergeTsFile(m3u8DO *M3u8DO) error {
	log := g.Log().Line()
	// 合并
	mp4File := fmt.Sprintf("%s%s%s", m3u8DO.FromDir, gfile.Separator, OrgMp4Name)
	args := []string{"-protocol_whitelist", ProtocolWhitelist, "-v", "error", "-y", "-i", m3u8DO.FromFile, "-c", "copy", mp4File}
	err := RunFfmpegCommand(args...)
	if err != nil {
		log.Infof(gctx.GetInitCtx(),
			"转换出错->%v", err)
		return err
	}
	log.Info(gctx.GetInitCtx(), "转换成MP4成功：", mp4File)
	return nil
}

func MergeVideoAudio(videoDir, videoFile, audioFile string) error {
	log := g.Log().Line()
	mp4File := fmt.Sprintf("%s%s", videoDir, OrgMp4Name)
	// 合并
	args := []string{"-i", videoFile, "-i", audioFile, "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", mp4File}
	err := RunFfmpegCommand(args...)
	if err != nil {
		log.Infof(gctx.GetInitCtx(),
			"转换出错->%v", err)
		return err
	}
	log.Info(gctx.GetInitCtx(), "转换成MP4成功：", videoFile)
	return nil
}

func DeleteTmpResource(m3u8DO *M3u8DO, tsFilePattern string) {
	log := g.Log().Line()

	// 删除依赖资源文件
	log.Infof(gctx.GetInitCtx(),
		"删除Map文件: %s", m3u8DO.MapSaveFile)
	log.Infof(gctx.GetInitCtx(),
		"删除Key文件: %s", m3u8DO.KeySaveFile)
	_ = gfile.Remove(m3u8DO.MapSaveFile)
	_ = gfile.Remove(m3u8DO.KeySaveFile)

	// 删除流片段文件
	streamFiles, _ := gfile.ScanDir(m3u8DO.FromDir, tsFilePattern)
	log.Info(gctx.GetInitCtx(), "删除流片段个数：", len(streamFiles))
	for _, tsFile := range streamFiles {
		_ = gfile.Remove(tsFile)
	}
}

func DeleteLiveTmpResource(m3u8DO *M3u8DO, tsFilePattern string) {
	tsFiles, _ := gfile.ScanDir(m3u8DO.FromDir, tsFilePattern)

	sortedTsFiles := garray.NewSortedStrArrayFrom(tsFiles)
	// 保留个数
	retainLen := 50
	if sortedTsFiles.Len() <= retainLen {
		return
	}

	for _, tsFile := range sortedTsFiles.SubSlice(0, sortedTsFiles.Len()-retainLen) {
		_ = gfile.Remove(tsFile)
	}
}
