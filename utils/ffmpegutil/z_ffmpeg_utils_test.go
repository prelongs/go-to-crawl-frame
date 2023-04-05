package ffmpegutil

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
	"testing"
)

var (
	url  = "https://xlzycdn1.sy-precise.com:65/20220712/Ioe91OwL/2637kb/hls/index.m3u8"
	path = "D:\\cache2\\video\\CN\\2022\\12\\36\\org_index_bak2.m3u8"
)

func TestConvertM3U8(t *testing.T) {
	log := g.Log().Line()
	m3U8, _ := ConvertM3U8(url, "", path)
	log.Info(gctx.GetInitCtx(), m3U8)
}

func TestPaddingKey(t *testing.T) {
	ret := paddingKeyUrl("https://v.v1kd.com", "#EXT-X-KEY:METHOD=AES-128,URI=\"/20220510/EAlZ3rpV/2000kb/hls/key.key\",IV=0x9180b4da3f0c7e80975fad685f7f134e #EXTINF:6.416667")
	fmt.Println(ret)
}

func TestFormatExtName(t *testing.T) {
	imgMode := false
	name := formatExtName("http://xxx.com/a.jpg?sign=khi34h", &imgMode)
	fmt.Println(name)
	fmt.Println(imgMode)
}

func TestTruncateTS(t *testing.T) {
	truncateTS("D:\\cache2\\tnt-626299.png", SrcTypeImg, 308)
}

func TestIsPngType(t *testing.T) {
	pngType := IsPngType("D:\\cache2\\tnt-626299.png")
	fmt.Println(pngType)
}

func TestRunFfmpegGenericMerge(t *testing.T) {
	err := RunFfmpegGenericMerge("D:\\cache2\\replay")
	fmt.Println(err)
}

func TestDeleteLiveTmpResource(t *testing.T) {
	m3u8 := new(M3u8DO)
	m3u8.FromDir = "D:\\app\\cms\\video\\live\\CNN"

	DeleteLiveTmpResource(m3u8, tsFilePattern)
}

func TestGetPlaySecond(t *testing.T) {
	playSecondMatch, _ := gregex.MatchString("\\d+\\.?\\d*", "#EXTINF:4.004,")
	fmt.Println(gconv.Float32(playSecondMatch[0]))
}
