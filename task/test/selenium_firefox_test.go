package test

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"net"
	"os"
	"testing"
)

const geckoDriverPath = "D:\\ApplicationsPro\\BrowserDriver\\geckodriver.exe"

func pickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func TestStartCrawl(t *testing.T) {
	port, err := pickUnusedPort()
	fmt.Println("port: ", port)

	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath),
		selenium.Output(os.Stderr),
	}
	selenium.SetDebug(false)

	service, err := selenium.NewGeckoDriverService(geckoDriverPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()
	fmt.Println("Here 1")
	fireCap := firefox.Capabilities{}
	fireCap.Args = append(fireCap.Args, "--profile", "D:\\ApplicationsPro\\BrowserDriver\\firefoxProfile")

	caps := selenium.Capabilities{"browserName": "firefox"}
	caps.AddFirefox(fireCap)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()
	if err = wd.Get("https://www.baidu.com"); err != nil {
		panic(err)
	}

	fmt.Println("Here 2")
}
