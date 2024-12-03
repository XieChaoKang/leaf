package common

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func TailFile(filename string, size int64) []byte {
	f, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fileSize := f.Size()
	if fileSize == 0 {
		return nil
	}

	if size > fileSize {
		size = fileSize
	}

	res := make([]byte, size)
	pos := fileSize - size

	var fi *os.File
	fi, err = os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer fi.Close()
	fi.Seek(pos, os.SEEK_SET)

	readSize, _ := fi.Read(res)
	if int64(readSize) < size {
		if readSize == 0 {
			return nil
		}

		return res[0:readSize]
	}

	return res
}

func GetDirPath(filename string) string {
	pos := strings.LastIndexByte(filename, filepath.Separator)
	if pos < 0 {
		pos = strings.LastIndex(filename, "/")
	}

	if pos > 0 {
		return filename[0:pos]
	}

	return ""
}

func GetIP(r *http.Request) string {
	//接入CF
	ip := r.Header.Get("cf-connecting-ip")
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip = r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip = r.Header.Get("X-Forwarded-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	if net.ParseIP(ip) != nil {
		return ip
	}
	return ""
}
