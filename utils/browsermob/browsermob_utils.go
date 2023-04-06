package browsermob

import (
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/service/browsermobservice"
	"github.com/JervisPG/go-to-crawl-frame/service/crawl/sysservice"
	"github.com/JervisPG/go-to-crawl-frame/utils/constant"
	"github.com/JervisPG/go-to-crawl-frame/utils/fileutil"
	"github.com/JervisPG/go-to-crawl-frame/utils/timeutil"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/tebeka/selenium"
	"io"
)

const ResponseBody = "responseBody" // 外部根据情况使用响应体

func GetHarRequestLocalRetry(proxy *browsermobservice.Client, patternUrl string, patternContent string) *gjson.Json {
	return GetHarRequest(proxy, patternUrl, patternContent, constant.LocalRetry)
}

func GetHarRequest(proxy *browsermobservice.Client, patternUrl string, patternContent string, retry int) *gjson.Json {
	array := GetHarRequestArray(proxy, patternUrl, patternContent, retry, true)
	if len(array) == 0 {
		return nil
	}
	return array[0]
}

func GetHarRequestArray(proxy *browsermobservice.Client, patternUrl string, patternContent string, retry int, returnOnGetOne bool) []*gjson.Json {
	log := g.Log().Line()
	var targetRequestArray []*gjson.Json

	result := proxy.Har()
	if result == nil {
		return nil
	}

	array, _ := result.Get("log").Get("entries").Array()
	log.Info(gctx.GetInitCtx(), "浏览器请求总数 = ", len(array))

	for idx, entry := range array {
		json := gjson.Json{}
		_ = gconv.Struct(entry, &json)
		reqJson := json.GetJson("request")
		reqUrl := reqJson.Get("url").String()
		rspContent := json.Get("response.content.text").String()

		var targetRequest *gjson.Json

		if gstr.ContainsI(reqUrl, patternUrl) {
			log.Infof(gctx.GetInitCtx(),
				"idx = %d, reqUrl = %s", idx, reqUrl)
			log.Infof(gctx.GetInitCtx(),
				"rsp body part: %s", gstr.SubStrRune(rspContent, 0, 1000))
			proxyUrl := sysservice.GetProxyByUrl(reqUrl)

			builder := fileutil.CreateBuilder().Url(reqUrl).Proxy(proxyUrl).Retry(fileutil.Retry)

			method := reqJson.Get("method").String()
			builder = builder.Method(method)
			builder = initHeaderFromHarItem(builder, reqJson)

			if fileutil.POST == method {
				// 使用原始header先对post开放，实际上get也适用
				builder = builder.Body(reqJson.Get("postData.text").String())
			}

			body := fileutil.DownloadToReaderByBuilder(builder)
			if body == nil {
				continue
			}
			bytes, err := io.ReadAll(body)
			if err != nil {
				continue
			}
			m3u8Content := string(bytes)

			if patternContent != "" {
				// 响应内容为非M3U8格式，需要各自抓取实现类再去提取
				if gstr.ContainsI(m3u8Content, patternContent) {
					targetRequest = reqJson
					_ = targetRequest.Append(ResponseBody, m3u8Content)
					targetRequestArray = append(targetRequestArray, targetRequest)
					if returnOnGetOne {
						break
					}
				}
			} else {
				// 响应内容直接为M3U8格式
				if gstr.ContainsI(m3u8Content, "EXT-X-") {
					targetRequest = reqJson
					_ = targetRequest.Append(ResponseBody, m3u8Content)
					targetRequestArray = append(targetRequestArray, targetRequest)
					if returnOnGetOne {
						break
					}
				}
			}
		} else if gregex.IsMatchString(patternUrl, reqUrl) {
			log.Infof(gctx.GetInitCtx(),
				"idx = %d, reqUrl = %s", idx, reqUrl)
			targetRequest = reqJson
			targetRequestArray = append(targetRequestArray, targetRequest)
			if returnOnGetOne {
				break
			}
		}
	}

	if returnOnGetOne {
		if len(targetRequestArray) == 0 && retry > 0 {
			timeutil.Sleep(5)
			targetRequestArray = GetHarRequestArray(proxy, patternUrl, patternContent, retry-1, returnOnGetOne)
		}
	} else {
		// 需要多个结果的
		// 场景一：len > 0，爱奇艺WEB UA下，可以抓到多个清晰度的M3U8，因此需要重试，补充更多后加载的M3U8
		if retry > 0 {
			// 还有retry次数的
			timeutil.Sleep(5)
			targetRequestArray = GetHarRequestArray(proxy, patternUrl, patternContent, retry-1, returnOnGetOne)
		}
	}

	printTargetUrl(targetRequestArray)
	return targetRequestArray
}

func GetHarResponseBody(harRequest *gjson.Json) string {
	if harRequest == nil {
		return ""
	}

	rspBody := harRequest.Get(ResponseBody).String()
	item := gjson.New(rspBody).Array()[0]
	return fmt.Sprintf("%s", item)
}

func NewHar(proxy *browsermobservice.Client) {
	proxy.NewHar("", map[string]string{
		"captureHeaders": "true",
		"captureContent": "true",
	})
}

func NewHarWait(wd selenium.WebDriver, proxy *browsermobservice.Client) {
	NewHar(proxy)
	_ = wd.SetPageLoadTimeout(gtime.S * 10)
}

func printTargetUrl(targetRequestArray []*gjson.Json) {
	for _, targetRequest := range targetRequestArray {
		targetUrl := ""
		if targetRequest != nil {
			targetUrl = targetRequest.Get("url").String()
		}
		g.Log().Line().Info(gctx.GetInitCtx(), "代理获取的目标 URL: ", targetUrl)
	}
}

func initHeaderFromHarItem(builder *fileutil.DownloadBuilder, reqJson *gjson.Json) *fileutil.DownloadBuilder {
	headers := reqJson.Get("headers").Map()
	for _, headerMap := range headers {
		headerJson := gjson.New(headerMap)
		key := headerJson.Get("name").String()
		value := headerJson.Get("value").String()

		if "Accept-Encoding" == key {
			// 解决响应乱码问题, 原为https的gzip格式
			value = ""
		}

		builder = builder.Header(key, value)
	}
	return builder
}
