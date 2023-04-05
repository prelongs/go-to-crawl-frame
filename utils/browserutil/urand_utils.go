package browserutil

import (
	"github.com/corpix/uarand"
	"github.com/gogf/gf/v2/text/gstr"
)

// Mozilla/5.0 (iPad; CPU OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Flipboard/4.2.142
const (
	Android   = "android"
	Iphone    = "iphone"
	SymbianOS = "SymbianOS"
	IPad      = "iPad"
)

const (
	IE = "MSIE"
)

func GetRandomUA(includeMobile bool) string {
	ua := uarand.GetRandom()

	// 排除指定类型UA
	exList := [...]string{IE}
	for _, item := range exList {
		if gstr.ContainsI(ua, item) {
			return GetRandomUA(includeMobile)
		}
	}

	if includeMobile {
		return ua
	}

	// 排除手机端UA
	UAFilterList := [...]string{Android, Iphone, SymbianOS, IPad}
	for _, item := range UAFilterList {
		if gstr.ContainsI(ua, item) {
			return GetRandomUA(includeMobile)
		}
	}

	return ua
}

func GetRandomUAByPlatform(platform string) string {
	ua := uarand.GetRandom()
	if gstr.ContainsI(ua, platform) {
		return ua
	} else {
		return GetRandomUAByPlatform(platform)
	}
}
