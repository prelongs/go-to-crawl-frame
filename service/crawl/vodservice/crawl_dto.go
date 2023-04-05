package vodservice

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"net/http"
)

type CrawlSeedParams struct {
	SeedUrlReqHeader http.Header `json:"seedUrlReqHeader"`
	M3u8UrlReqHeader http.Header `json:"m3U8UrlReqHeader"`

	M3u8CompUrl `json:"m3U8CompUrl"`
}

type M3u8CompUrl struct {
	VideoUrl string `json:"videoUrl"`
	AudioUrl string `json:"audioUrl"`
}

func (r *CrawlSeedParams) AddSeedUrlReqHeader(key, value string) {
	if r.SeedUrlReqHeader == nil {
		r.M3u8UrlReqHeader = make(http.Header)
	}
	r.M3u8UrlReqHeader.Add(key, value)
}

func (r *CrawlSeedParams) AddM3U8UrlReqHeader(key, value string) {
	if r.M3u8UrlReqHeader == nil {
		r.M3u8UrlReqHeader = make(http.Header)
	}
	r.M3u8UrlReqHeader.Add(key, value)
}

func (r *CrawlSeedParams) ToJsonString() string {
	return gjson.New(r).MustToJsonString()
}

func (r *CrawlSeedParams) InitSeedParams(json string) {
	_ = gjson.DecodeTo(json, r)
}

func IsVideoAudioSeparated(crawlSeedParams string) bool {
	if crawlSeedParams == "" {
		return false
	}

	params := new(CrawlSeedParams)
	params.InitSeedParams(crawlSeedParams)

	return params.AudioUrl != ""
}
