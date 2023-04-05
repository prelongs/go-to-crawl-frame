package taskdto

import (
	"fmt"
	"testing"
)

func TestGetInitedCrawlContext(t *testing.T) {
	context := GetInitedCrawlContext()
	fmt.Println(context)
}
