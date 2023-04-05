package crawltask

import (
	"fmt"
	"testing"
)

func TestGetCrawlVodFlowStrategy(t *testing.T) {
	strategy := doGetCrawlVodFlowStrategy("https://www.nivod3.tv/filter.html?x=1&channelId=7&showTypeId=145")
	fmt.Println(strategy)
}
