package replayservice

import (
	"fmt"
	"testing"
)

func TestGetDetailById(t *testing.T) {
	do := GetDetailById(3)
	fmt.Println(do)
}
