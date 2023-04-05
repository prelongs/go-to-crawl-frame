package ftp

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"

	"os"
	"path"
	"time"
)

func sftpUpload() {
	connect := getConnect()
	defer connect.Close()
	uploadFile(connect, "D:\\cache\\飞刀问情_01.ts", "/home/sftptest")
}

func getConnect() *sftp.Client {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)

	// 创建ssh连接
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password("Quuiyai6s")) // 抱歉，这是我电脑密码
	clientConfig = &ssh.ClientConfig{
		User:            "root", // windows账户
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr = fmt.Sprintf("%s:%d", "107.167.22.34", 22) // SSH的22端口
	sshClient, err = ssh.Dial("tcp", addr, clientConfig)
	if nil != err {
		fmt.Println("ssh.Dial error", err)
	} else {
		fmt.Println("ssh.Dial 成功")
	}

	// 创建sftp连接
	sftpClient, err = sftp.NewClient(sshClient)
	if nil != err {
		fmt.Println("sftp.NewClient error", err)
	} else {
		fmt.Println("sftp.NewClient 成功")
	}

	return sftpClient
}

func uploadFile(sftpClient *sftp.Client, localFile, remotePath string) {
	file, err := os.Open(localFile)
	if nil != err {
		fmt.Println("os.Open error", err)
		return
	}
	fi, _ := file.Stat()
	fmt.Printf("文件名: %s, 文件大小: %d\n", fi.Name(), fi.Size())
	defer file.Close()

	ftpFile, err := sftpClient.Create(path.Join(remotePath, fi.Name())) // 这里的remotePath是sftp根目录下的目录，是目录不是文件名
	if nil != err {
		fmt.Println("sftpClient.Create error", err)
		return
	}
	defer ftpFile.Close()

	buf := &WriteCounter{
		Size: fi.Size(),
		File: ftpFile,
	}

	ioutil.ReadAll(io.TeeReader(file, buf))
}

type WriteCounter struct {
	*sftp.File
	Size    int64
	HasRead int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.HasRead += int64(n)
	wc.File.Write(p)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	var p = gconv.Float32(wc.HasRead) / gconv.Float32(wc.Size)
	fmt.Printf("completed %f.\n", p*100)
}
