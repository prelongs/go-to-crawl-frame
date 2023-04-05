package browsermob

import (
	"fmt"
	"github.com/gogf/gf/v2/text/gregex"
	"testing"
)

func TestPattern(t *testing.T) {
	str := gregex.IsMatchString("/x/player/playurl", "https://api.bilibili.com/x/player/playurl?avid=504150088&bvid=BV1yg411T7Za")
	fmt.Println(str)
}
