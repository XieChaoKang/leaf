package log

import (
	"testing"
	"time"
)

func TestAppLog(t *testing.T) {
	dir := "./"
	filename := dir + "app.log.%Y%m%d-%H%M"
	appLog := AppLogNew(filename, 0)
	// kv format
	appLog.FmtHandler = WAFmtHandle
	obj := map[string]interface{}{
		"logfile": filename,
		"map2": map[string]interface{}{
			"logfile2": filename + "2",
			"int":      123,
		},
		"int": 123,
	}
	appLog.Log(obj)
	appLog.FmtHandler = JSONFmtHandle
	appLog.Log(obj)
	appLog.LogBytes([]byte("byte test"))
	appLog.LogString("string test")
	appLog.LogMap(map[string]string{"k": "v", "1": "2"})
	appLog.CleanUp()
	//files, _ := vio.ReadDir(dir)
	//for _, fileinfo := range files {
	//	b, _ := vio.Get(fileinfo.Name())
	//	fmt.Println(string(b))
	//}
}

// 使用默认配置直接打印
func TestGameRecordLog(t *testing.T) {
	GameRecord(&GameRecordLog{
		ID:       1,
		SubAnm:   "vgame",
		Did:      "111111111",
		PkgName:  "com.demo",
		RoomType: "user",
	})

	time.Sleep(1 * time.Second)
}
