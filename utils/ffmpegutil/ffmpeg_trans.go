package ffmpegutil

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"os/exec"
	"path/filepath"
	"strings"
)

type mFFMPEG struct {
	//ffmpeg执行命令路径(方便不配置环境变量时候new全路径)
	FFMPEG_PATH string
}

func FmpegTrans(ffmpegPath string) *mFFMPEG {
	return &mFFMPEG{FFMPEG_PATH: ffmpegPath}
}

//过滤文件(操作文件夹,源文件路径,处理后的文件路径)
func (f *mFFMPEG) CheckFile(FileDir, Path_Source, Path_Target string) error {
	sizebyte := gfile.GetBytes(Path_Source)
	if len(sizebyte) < 1024*1024*2 {
		return fmt.Errorf("file too small")
	}
	_, str := f.Exec("-i", Path_Source)
	isTrans := f.IsTransAudio(str)
	if isTrans {
		//	需要转码 ->音频格式转为AAC标准音频
		g.Log().Infof(gctx.GetInitCtx(), "change file package to mp4 and audio to aac")
		f.Exec("-v", "error", "-y", "-i", Path_Source, "-acodec", "aac", "-vcodec", "copy", Path_Target)
		//g.Dump(o,s)
	} else {
		//	不需要转码 ->封装格式转为mp4
		//g.Log().Infof(gctx.GetInitCtx(), "change file package to mp4")
		f.Exec("-v", "error", "-y", "-i", Path_Source, "-c", "copy", Path_Target)
		//g.Dump(o,s)
	}
	if !gfile.IsFile(Path_Target) {
		//转码后目标文件不存在则判断为转码失败
		return fmt.Errorf("trans fail source:%v", Path_Source)
	}

	g.Log().Infof(gctx.GetInitCtx(), "change file to m3u8")
	f.Exec("-v", "error",
		"-y", "-i", Path_Target,
		"-c", "copy", "-f", "segment",
		"-segment_list", filepath.Join(FileDir, "index.m3u8"),
		"-segment_time", "5",
		filepath.Join(FileDir, "segment-%04d.ts"))

	// 不知道判断是否切片成功 暂定为判断index.m3u8存在则成功
	if gfile.IsFile(filepath.Join(FileDir, "index.m3u8")) {
		//转化成功删除源视频（删除失败也没关系 所以不做判断）
		_ = gfile.Remove(Path_Source)
		_ = gfile.Remove(Path_Target)
	} else {
		//	如果index.m3u8都不存在肯定无法播放判断为失败
		return fmt.Errorf("miss m3u8 file:%s", FileDir)
	}
	return nil
	//g.Dump(str,isTrans)
}

//判断是否需要转码音频(不是AAC则要转为AAC)
//参数为Exec返回的FFPMEG信息
func (f *mFFMPEG) IsTransAudio(findstr string) bool {

	arr := strings.Split(findstr, "\n")
	if len(arr) < 1 {
		return false
	}
	//记录音频流的数量(有双声道情况多个音频流的情况转码为默认aac音频流)
	//var AudioCount int=0
	var AudioLine []string = []string{}

	for i := 0; i < len(arr); i++ {
		//查询所有流信息所在行
		StreamLine := strings.Index(arr[i], "Stream")
		if StreamLine != -1 {
			//进入流信息行分支

			//	查找音频流所在行
			AudioLineIndex := strings.Index(arr[i], "Audio")
			if AudioLineIndex != -1 {
				//	音频流所在行
				AudioLine = append(AudioLine, arr[i])

			}
		}
	}

	switch len(AudioLine) {
	case 0:
		//	无音频流处理(无音频流不转码 本身没有音频流也转不了)
		return false
	case 1:
		//g.Dump("1111111111")
		itemLine := AudioLine[0]
		aacindex := strings.Index(itemLine, "Audio: aac")
		if aacindex == -1 {
			//	查找不到Audio: aac字符串暂时认为就不是aac格式(可能ffmpeg版本不一样输出不一样待定)

			//不是AAC格式则转码
			return true
		} else {
			//是aac格式不需要转码
			return false
		}
		//g.Dump(aacindex)
		//	只有一个音频流情况处理
		break
	default:
		//	大于1个音频流处理(直接返回true表示需要转码)
		return true
	}
	return false

}

//执行FFMPEG命令输出返回值
func (f *mFFMPEG) Exec(args ...string) (strout, strerr string) {
	cmd := exec.Command(f.FFMPEG_PATH, args...)
	g.Log().Infof(gctx.GetInitCtx(), "执行命令 do cmd:%s", cmd.String())
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		//log.Println("exec the cmd ", name, " failed")
		return "", ""
	}
	// 正常日志
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			//log.Println(logScan.Text())
			strout = logScan.Text()
		}
	}()
	// 错误日志
	errBuf := bytes.NewBufferString("")
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		s := scan.Text()

		//log.Println("build error: ", s)

		errBuf.WriteString(s)
		errBuf.WriteString("\n")
		strerr = errBuf.String()
	}
	// 等待命令执行完
	cmd.Wait()
	return
}
