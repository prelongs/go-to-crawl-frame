package ftp

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"time"
)

func ftpUpload() {
	c, _ := ftp.Dial("107.167.22.34:22", ftp.DialWithTimeout(5*time.Second))
	defer c.Quit()
	//登陆
	_ = c.Login("root", "Quuiyai6s")
	_ = c.ChangeDir("/home/ftptest")
	dir, _ := c.CurrentDir()
	fmt.Print("current dir ", dir)
	//上传文件a.txt到远程ftp服务器/home/ftptest/111
	_ = c.MakeDir("111")
	_ = c.ChangeDir("/home/ftptest/111")
	file, _ := os.Open("a.txt")
	defer file.Close()
	_ = c.Stor("b.txt", file)
	//从ftp服务器下载文件b.txt到本地目录c:\\WORK\\src\\learn\\ftptest\\c.txt
	f, _ := os.OpenFile("D:\\tmp\\c.txt", os.O_RDWR|os.O_CREATE, 0777)
	defer f.Close()
	r, _ := c.Retr("/home/ftptest/111/b.txt")
	defer r.Close()
	io.Copy(f, r)
	c.Logout()
}
