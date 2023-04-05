package fileutil

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"go-to-crawl-frame/db/mysql/model/entity"
	"go-to-crawl-frame/service/crawl/sysservice"
	"strings"
	"testing"
)

var (
	testUrl = "https://xlzycdn1.sy-precise.com:65/20220712/Ioe91OwL/2637kb/hls/index.m3u8"
	path    = "/app/docker/go-fastdfs/data/files/video/CN/11/22/org_index.m3u8"
)

func TestDownloadFile(t *testing.T) {
	proxyUrl := sysservice.GetProxyByUrl(testUrl)
	err := DownloadFile(testUrl, proxyUrl, path, 2)
	fmt.Println(err)
}

func TestDownloadSeed(t *testing.T) {

	seed := new(entity.CmsCrawlQueue)
	seed.CrawlSeedUrl = "https://glmrtothemoon.xyz/tv247/footnw-135835.png"

	builder := CreateBuilder()
	builder.Url(seed.CrawlSeedUrl)

	path1 := "D:/cache2/"
	_ = gfile.Mkdir(path1)
	builder.SaveFile(path1 + "footnw-135835.png")
	builder.Header("origin", "https://tv247.us")
	builder.Header("referer", "https://tv247.us/")
	err2 := DownloadFileByBuilder(builder)

	fmt.Println(err2)
	fmt.Println("结束")
}

func TestSubStr(t *testing.T) {
	name := "111.jpg"
	str := gstr.SubStr(name, 0, strings.Index(name, "."))
	fmt.Println(str)
}
