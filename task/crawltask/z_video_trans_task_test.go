package crawltask

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"go-to-crawl-frame/utils/ffmpegutil"
	"testing"
)

// 解析命令逻辑有点费解，运行不了
func TestFfmpegGproc(t *testing.T) {
	cmd := "ffmpeg -i 'D:/刘星/video/CN/11/22/公信宝授权视频.mp4' -c:v h264 -flags +cgop -g 30 -hls_time 20 -hls_list_size 0 -hls_segment_filename D:/刘星/video/CN/11/22/index%3d.ts D:/刘星/video/CN/11/22/index.m3u8"
	r, err := gproc.ShellExec(gctx.GetInitCtx(), cmd)
	fmt.Println("result:", r)
	fmt.Println(err)
}

// 测试通过
func TestFfmpegExec(t *testing.T) {
	basePath := "D:\\cache2\\replay\\test"
	err := ffmpegutil.RunFfmpegGenericSlice(basePath)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGproc(t *testing.T) {
	r, err := gproc.ShellExec(gctx.GetInitCtx(), `dir g:\`)
	fmt.Println("result:", r)
	fmt.Println(err)
}

func TestTransformTask(t *testing.T) {
	TransformTask(gctx.GetInitCtx())
}
