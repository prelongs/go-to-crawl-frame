package vodservice

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

func TestGetVodConfigTaskDO(t *testing.T) {
	g.Dump(GetVodConfigTaskDO())
}

func TestGetPreparedVodTvItem(t *testing.T) {
	item := GetPreparedVodTvItem()
	g.Dump(item)
}

func TestToJsonString(t *testing.T) {
	params := new(CrawlSeedParams)
	params.AddM3U8UrlReqHeader("k1", "v1")
	params.VideoUrl = "hv"
	params.AudioUrl = "ha"

	json := params.ToJsonString()

	p1 := new(CrawlSeedParams)
	p1.InitSeedParams(json)
	fmt.Println(p1)
}
