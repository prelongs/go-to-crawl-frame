package fileutil

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	Retry = 2
)

func DownloadToReaderByBuilder(downloadBuilder *DownloadBuilder) io.ReadCloser {
	log := g.Log().Line()
	client := new(http.Client)

	downloadDO := downloadBuilder.Build()

	// 设置超时
	if downloadDO.timeout > 0 {
		client.Timeout = time.Duration(downloadDO.timeout) * time.Millisecond
	}

	if downloadDO.proxy != "" {
		// 设置代理
		proxyUrl, err := url.Parse(downloadDO.proxy)
		if err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		} else {
			log.Error(gctx.GetInitCtx(), err)
		}
	}

	start := time.Now().UnixMilli()
	// 获取资源
	log.Infof(gctx.GetInitCtx(),
		"DownloadToReaderByBuilder: method = %v, url = %v, body = %v", downloadDO.method, downloadDO.url, downloadDO.body)

	body := strings.NewReader(downloadDO.body)
	request, err := http.NewRequest(downloadDO.method, downloadDO.url, body)
	if err != nil || request == nil {
		log.Errorf(gctx.GetInitCtx(), "err = %v, request = %v", err, request)
	}

	// Headers
	for key, values := range downloadDO.headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	resp, err := client.Do(request)
	spend := time.Now().UnixMilli() - start

	if err != nil {
		log.Error(gctx.GetInitCtx(), err)
		if downloadDO.retry <= 0 {
			return nil
		} else {
			downloadDO.retry -= 1
			return DownloadToReaderByBuilder(downloadBuilder)
		}
	} else {
		log.Infof(gctx.GetInitCtx(),
			"DownloadToReaderByBuilder: %v|%v ms|%v bytes|%s", resp.StatusCode, spend, resp.ContentLength, downloadDO.url)
		return resp.Body
	}
}

func DownloadM3U8File(url, proxy, saveFile string, retry int, preRspBody string) error {
	builder := CreateBuilder().Url(url).Proxy(proxy).SaveFile(saveFile).Retry(retry).PreRspBody(preRspBody)
	return DownloadFileByBuilder(builder)
}

func DownloadFile(url string, proxy string, saveFile string, retry int) error {
	url = gstr.Replace(url, "ali.", "") // eg: https://ali.static.yximgs.com/bs2/adcarsku/sku9fcdd221-94ee-4f64-a20d-ad085a78d424.png
	builder := CreateBuilder().Url(url).Proxy(proxy).SaveFile(saveFile).Retry(retry)
	return DownloadFileByBuilder(builder)
}

func DownloadFileByBuilder(downloadBuilder *DownloadBuilder) error {
	downloadDO := downloadBuilder.Build()

	if gfile.Exists(downloadDO.saveFile) {
		gfile.Remove(downloadDO.saveFile)
	}

	if downloadDO.preRspBody == "" {
		body := DownloadToReaderByBuilder(downloadBuilder)
		if body == nil {
			return errors.New("下载失败. " + downloadDO.url)
		}
		defer body.Close()

		file, err := gfile.Create(downloadDO.saveFile)
		if err != nil {
			return err
		}
		defer file.Close()

		// 获得文件的writer对象
		writer := bufio.NewWriter(file)
		// 获得get请求响应的reader对象
		reader := bufio.NewReaderSize(body, 32*1024)

		io.Copy(writer, reader)
	} else {
		file, err := os.Create(downloadDO.saveFile)
		if err != nil {
			return err
		}
		defer file.Close()

		// 获得文件的writer对象
		writer := bufio.NewWriter(file)
		reader := bytes.NewReader([]byte(downloadDO.preRspBody))
		io.Copy(writer, reader)
	}

	return nil
}
