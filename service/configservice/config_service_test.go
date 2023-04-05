package configservice

import (
	"fmt"
	"testing"
)

func TestGetCrawlHostIp(t *testing.T) {
	fmt.Println(GetCrawlHostIp())
}

func TestGetCrawlBrowser(t *testing.T) {
	browser := GetCrawlBrowser()
	fmt.Printf("%v", browser)
}

func TestGetCrawlBrowserInfo(t *testing.T) {
	info := GetCrawlBrowserInfo()
	fmt.Printf("%v", info)
}
