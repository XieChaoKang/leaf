package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type TimeRotateFile struct {
	MaxFileCnt  int
	Locale      *time.Location
	PathFmt     string
	buf         bytes.Buffer
	lock        *sync.Mutex
	curFilePath string
	curFile     *os.File
	stop        bool
	cntlines    int
}

func TimeRotateFileNew(pathFmt string, maxFileCnt int) *TimeRotateFile {
	pathFmt = EnvString(pathFmt)
	this := &TimeRotateFile{
		PathFmt:    pathFmt,
		MaxFileCnt: maxFileCnt,
		buf:        bytes.Buffer{},
		lock:       new(sync.Mutex),
		stop:       false,
		cntlines:   0,
	}
	this.Locale, _ = time.LoadLocation("Asia/Shanghai")
	this.buf.Grow(1024 * 1024 * 4)
	this.Flush()
	go this.rountine()
	return this
}

func (this *TimeRotateFile) checkFileNum() {
	//check and delete files
	timeFields := strings.Split("[0-9][0-9][0-9][0-9],[0-9][0-9],[0-9][0-9],[0-9][0-9],[0-9][0-9]", ",")
	pathTimeFields := strings.Split("%Y-%m-%d-%H-%M", "-")
	path := this.PathFmt
	for idx, timeField := range timeFields {
		path = strings.Replace(path, pathTimeFields[idx], timeField, -1)
	}
	matches, _ := filepath.Glob(path)

	if matches != nil && len(matches) > this.MaxFileCnt+1 {
		for idx, filepath := range matches {
			var modifyTime int64 = 0
			finfo, err := os.Stat(filepath)
			if err == nil {
				modifyTime = finfo.ModTime().Unix()
			}
			matches[idx] = fmt.Sprintf("%10x-", modifyTime) + filepath
		}
		sort.Strings(matches)
		count := len(matches) - (this.MaxFileCnt + 1)
		for idx, filepath := range matches {
			if idx > count {
				break
			}
			os.Remove(strings.SplitN(filepath, "-", 2)[1])
		}
	}
}
func (this *TimeRotateFile) rountine() {
	idx := 0
	for {
		if this.stop {
			break
		}
		time.Sleep(time.Millisecond * 100)
		if this.buf.Len() == 0 {
			idx += 1
			if idx%10 == 0 {
				this.chkFile()
			}
			continue
		}
		this.Flush()
	}
}
func (this *TimeRotateFile) WriteString(s string) {
	this.lock.Lock()
	this.buf.WriteString(s)
	this.cntlines++
	this.lock.Unlock()
	if this.cntlines > 100 || (this.cntlines > 4 && this.buf.Len() > 1048576) {
		go this.Flush()
	}
}
func (this *TimeRotateFile) Write(b []byte) {
	this.lock.Lock()
	this.buf.Write(b)
	this.cntlines++
	this.lock.Unlock()
	if this.cntlines > 100 || (this.cntlines > 4 && this.buf.Len() > 1048576) {
		go this.Flush()
	}
}
func (this *TimeRotateFile) getDatePath() string {
	path := this.PathFmt
	return GetFilePath(path, this.Locale)
}

func (this *TimeRotateFile) chkAndCleanBuf() {
	if this.buf.Len() > 1024*1024*100 {
		this.buf.Reset()
	}
}

func (this *TimeRotateFile) PeekEnd(size int64) []byte {
	if this.curFilePath != "" {
		return TailFile(this.curFilePath, size)
	}
	return nil
}

func (this *TimeRotateFile) chkFile() {

	this.lock.Lock()
	defer this.lock.Unlock()
	filename := this.getDatePath()
	if filename != this.curFilePath {
		if this.curFile != nil {
			this.curFile.Close()
		}
		if !PathExists(filename) {
			dirpath := GetDirPath(filename)
			if dirpath != "" {
				if !PathExists(dirpath) {
					err := os.MkdirAll(dirpath, os.ModePerm)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
		}
		this.curFile = file
		this.curFilePath = filename
		this.checkFileNum()
	}
}

func (this *TimeRotateFile) Flush() bool {
	if this.buf.Len() == 0 && this.curFilePath != "" {
		return true
	}
	this.cntlines = 0
	this.chkFile()
	this.lock.Lock()
	defer this.lock.Unlock()
	this.cntlines = 0

	_, err := this.curFile.Write(this.buf.Bytes())
	if err != nil {
		fmt.Println(time.Now(), "TimeRotateFile flush fail", err)
	}
	//this.curFile.Sync()
	this.buf.Reset()
	return true
}

func (this *TimeRotateFile) CleanUp() {

	this.Flush()
	this.stop = true
	if this.curFile != nil {
		this.curFile.Close()
	}
}
