package lockservice

import (
	"fmt"
	"testing"
)

func TestTryLockSelenium(t *testing.T) {
	fmt.Println(TryLockSelenium())
	fmt.Println(TryLockSelenium())
}

func TestAddSeleniumExpireTime(t *testing.T) {
	TryLockSelenium()
	for i := 0; i < 10; i++ {
		AddSeleniumExpireTime(60)
	}
}

func TestIncreaseValue(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Println(IncreaseValue("key"))
	}

	for i := 0; i < 6; i++ {
		fmt.Println(DecreaseValue("key"))
	}

}
