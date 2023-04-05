package browserutil

import (
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"testing"
)

func TestURand(t *testing.T) {
	for i := 0; i < 500; i++ {
		ua := GetRandomUA(true)
		if gstr.ContainsI(ua, "iPhone") {
			fmt.Println(ua)
		}
	}
	fmt.Println("end")
}
