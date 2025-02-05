package util

import "time"

// 所有都是本地时间

const timeFormat = "2006-01-02 00:00:00"

func ConvertStrTime(strTime string) (time.Time, bool) {
	t, err := time.ParseInLocation(timeFormat, strTime, time.Local)
	return t, err == nil // if error, time will return "0001-01-01 00:00:00"
}

func FormatTime(t time.Time) string {
	return t.Format(timeFormat)
}

func NowTimeString() string {
	return time.Now().Format(timeFormat)
}

func UnixTimeString(unix int64) string {
	return time.Unix(unix, 0).Format(timeFormat)
}

func NowTime() time.Time {
	return time.Now()
}

func NowUnixTime() int64 {
	return time.Now().Unix()
}

func NowUnixNanoTime() int64 {
	return time.Now().UnixNano()
}

func IsToday(t time.Time) bool {
	nowTime := NowTime()
	return t.Year() == nowTime.Year() && t.Month() == nowTime.Month() && t.Day() == nowTime.Day()
}

func AddTime(dueTime time.Time, dur time.Duration) time.Time {
	nowTime := NowTime()
	if dueTime.Before(nowTime) {
		dueTime = nowTime
	}
	return dueTime.Add(dur)
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
