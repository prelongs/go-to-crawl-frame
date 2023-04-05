package liveservice

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"testing"
)

func TestGetLiveConfigList(t *testing.T) {
	for _, liveConfig := range GetLiveConfigList() {
		fmt.Println(gjson.New(liveConfig).MustToJsonString())
	}
}
