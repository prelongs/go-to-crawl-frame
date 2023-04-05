package middleware

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"testing"
)

func TestGetWhiteList(t *testing.T) {
	list := getWhiteList(g.Cfg())
	for _, uri := range list {
		fmt.Println(uri)
	}
}
