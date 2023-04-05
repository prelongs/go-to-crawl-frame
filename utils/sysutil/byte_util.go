package sysutil

import "fmt"

func FormatSize(bytes int) (size string) {
	if bytes < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(bytes)/float64(1))
	} else if bytes < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(1024))
	} else if bytes < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(1024*1024))
	} else if bytes < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(1024*1024*1024))
	} else if bytes < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(bytes)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fPB", float64(bytes)/float64(1024*1024*1024*1024*1024))
	}
}
