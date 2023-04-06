package browsermobservice

import (
	"fmt"
	"github.com/JervisPG/go-to-crawl-frame/utils/processutil"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	_LOG_FILE = "log/browsermob-proxy.out.log"
	PORT      = 8080
)

type Server struct {
	Path    string      `json:"path"`
	Host    string      `json:"host"`
	Port    int         `json:"port"`
	Process *os.Process `json:"process"`
	Command string      `json:"command"`
	Url     string      `json:"url"`
}

//initialises a server object
func NewServer(path string) *Server {
	server := new(Server)
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(path, ".bat") {
			path += ".bat"
		}
	}
	server.Path = path
	server.Host = "localhost"
	server.Port = PORT
	server.Url = fmt.Sprintf("http://%s:%d", server.Host, server.Port)
	return server
}

//启动
func (s *Server) Start() {
	StopBrowserMobProxy(true)
	stdOut, _ := os.OpenFile(_LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	//r, w, err := os.Pipe()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer r.Close()
	procAttr := &os.ProcAttr{
		Files: []*os.File{nil, stdOut, stdOut},
	}
	process, err := os.StartProcess(s.Path, nil, procAttr) //运行脚本

	if err != nil {
		fmt.Println("look path error:", err)
		os.Exit(1)
	}

	//buf := make([]byte,1024)
	//for{
	//	n,err := r.Read(buf)
	//	if err != nil && err != io.EOF{panic(err)}
	//	if 0 ==n {break}
	//	fmt.Println(string(buf[:n]))
	//}
	//w.Close()
	s.Process = process
	time.Sleep(2 * time.Second)

}

//判断是否启动
func (s *Server) isListen() bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", s.Port))
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return false
	}
	defer conn.Close()
	return true
}

//停止
func (s *Server) Stop() {
	StopBrowserMobProxy(false)
}

// hangingTooLong=true表示挂起太久
func StopBrowserMobProxy(hangingTooLong bool) {
	// 停止proxy-server进程
	pid, _ := processutil.CheckProcessRunning("browsermob-proxy")
	if pid != "" {
		processutil.KillPid(pid)
		// 停止子进程
		time.Sleep(time.Millisecond * 500)
		g.Log().Line().Infof(gctx.GetInitCtx(), "StopBrowserMobProxy. hangingTooLong = %v, pid = %s", hangingTooLong, pid)
		StopBrowserMobProxy(hangingTooLong)
	}
}

//创建代理
func (s *Server) CreateProxy(param Params) *Client {
	client := NewClient(s.Url[7:], param, nil)
	return client
}
