package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"leaf/common"
	"leaf/common/uuid"
	"os"
	"sort"
	"strings"
	"sync"
)

func WAFmtHandle(body map[string]interface{}) []byte {
	buffer := new(bytes.Buffer)
	firstOne := true
	sortedKeys := make([]string, 0)
	for k, _ := range body {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		val := body[k]
		if firstOne {
			firstOne = false
		} else {
			buffer.WriteString("`")
		}
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(common.ToString(val))
	}

	buffer.WriteString("\n")
	return buffer.Bytes()
}

func JSONFmtHandle(body map[string]interface{}) []byte {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(&body)
	return buf.Bytes()
}

func JSONStructHandle(body interface{}) []byte {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(&body)
	return buf.Bytes()
}

type AppLog struct {
	file             *common.TimeRotateFile
	pathFmt          string
	extFiles         map[string]*common.TimeRotateFile
	lock             *sync.Mutex
	WriteHandler     func([]byte)
	maxFiles         int
	FmtHandler       func(map[string]interface{}) []byte
	StructFmtHandler func(interface{}) []byte
}

func AppLogNew(pathFmt string, maxFiles int) *AppLog {
	this := &AppLog{
		FmtHandler:       WAFmtHandle,
		StructFmtHandler: JSONStructHandle,
		extFiles:         map[string]*common.TimeRotateFile{},
		lock:             new(sync.Mutex),
		pathFmt:          pathFmt,
		maxFiles:         maxFiles,
	}
	if pathFmt != "" {
		this.file = common.TimeRotateFileNew(pathFmt, maxFiles)
		this.WriteHandler = func(b []byte) {
			this.file.Write(b)
		}
	} else {
		this.WriteHandler = func(b []byte) {
			fmt.Println(string(b))
		}
	}
	return this
}

func (this *AppLog) UpdateMaxFileCnt(cnt int) {
	this.file.MaxFileCnt = cnt
}

func (this *AppLog) Flush() {
	if this.file != nil {
		this.file.Flush()
	}
}

func (this *AppLog) PeekEnd(size int64) []byte {
	if this.file != nil {
		return this.file.PeekEnd(size)
	}
	return nil
}
func (this *AppLog) CleanUp() {
	if this.file != nil {
		this.file.CleanUp()
	}
}

func (this *AppLog) LogBytes(b []byte) {
	this.WriteHandler(b)
}

func (this *AppLog) LogString(line string) {
	this.WriteHandler([]byte(line + "\n"))
}

func (this *AppLog) LogMap(minfo map[string]string) {
	this.WriteHandler(common.ToBytes(common.Map2Str(minfo, "`", "=") + "\n"))
}

func (this *AppLog) getFile(filePrefix string) *common.TimeRotateFile {
	this.lock.Lock()
	defer this.lock.Unlock()
	if f, ok := this.extFiles[filePrefix]; ok {
		return f
	}

	fields := strings.Split(this.pathFmt, "/")
	fields[len(fields)-1] = filePrefix + fields[len(fields)-1]
	f := common.TimeRotateFileNew(strings.Join(fields, "/"), this.maxFiles)
	this.extFiles[filePrefix] = f
	return f
}

func (this *AppLog) LogBytesSpec(filePrefix string, b []byte) {
	this.getFile(filePrefix).Write(b)
}

func (this *AppLog) LogMapSpec(filePrefix string, minfo map[string]string) {
	this.getFile(filePrefix).Write(common.ToBytes(common.Map2Str(minfo, "`", "=") + "\n"))
}

func (this *AppLog) Log(body map[string]interface{}) {
	if body == nil || len(body) == 0 {
		return
	}

	this.WriteHandler(this.FmtHandler(body))
}

func (this *AppLog) LogStruct(body interface{}) {
	if body == nil {
		return
	}

	this.WriteHandler(this.StructFmtHandler(body))
}

func (this *AppLog) GameRecord(body *GameRecordLog) {
	if body == nil {
		return
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(&body)
	this.WriteHandler(buf.Bytes())
}

func (this *AppLog) EventRecord(body *EventLog) {
	if body == nil {
		return
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(&body)
	this.WriteHandler(buf.Bytes())
}

func appLogDefaultConfig() (string, int) {
	maxFiles := 24
	pathFmt := "~/logs/apps/game/"
	servName := ""
	if servNameTmp, ok := os.LookupEnv("ServName"); ok {
		servName = servNameTmp
	} else {
		appPaths, _ := os.Getwd()
		appPath := strings.Split(appPaths, "/")
		servName = appPath[len(appPath)-1]
	}

	pathFmt = pathFmt + servName + "_%Y%m%d%H.log"
	return pathFmt, maxFiles
}

var __appLogger *AppLog

func GameLogger() *AppLog {
	if __appLogger == nil {
		__appLogger = AppLogNew(appLogDefaultConfig())
	}

	return __appLogger
}

func AppLogExport(logger *AppLog) {
	if logger != nil {
		__appLogger = logger
	}
}

func GameRecord(body *GameRecordLog) {
	if __appLogger == nil {
		__appLogger = AppLogNew(appLogDefaultConfig())
	}

	_uuid, _ := uuid.NewV4()
	body.LogId = _uuid.String()
	__appLogger.GameRecord(body)
}
