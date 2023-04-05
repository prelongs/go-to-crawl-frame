package httputil

import (
	"fmt"
	"testing"
)

func TestGetBaseUrlBySchema(t *testing.T) {
	fmt.Println(GetBaseUrlBySchema("https://t12.cdn2020.com:12337/video/m3u8/2021/10/28/12abec8c/index.m3u8"))
}

func TestGetBaseUrlByBackslash(t *testing.T) {
	fmt.Println(GetBaseUrlByBackslash("https://t12.cdn2020.com:12337/video/m3u8/2021/10/28/12abec8c/index.m3u8"))
}
