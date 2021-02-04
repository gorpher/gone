package gone

import (
	"strconv"
	"strings"
	"time"
)

const timeTpl string = "2006-01-02 15:04:05"

// https://github.com/dxvgef/gommon/blob/master/datatime/datetime.go
// const timeShortTpl string = "2006-01-02 15:04"

// TimeToStr 返回时间的字符串格式
func TimeToStr(t time.Time, format ...string) string {
	if len(format) == 0 {
		return t.Format(timeTpl)
	}
	return t.Format(format[0])
}

// Timestamp将unix时间转为时间字符串
func TimestampToStr(t int64, format ...string) string {
	if len(format) == 0 {
		return time.Unix(t, 0).Format(timeTpl)
	}
	return time.Unix(t, 0).Format(format[0])
}

// FormatByStr 将字符串中的时间变量（y年/m月/d日/h时/i分/s秒）转换成时间字符串
func FormatByStr(tpl string, t int64) string {
	tpl = strings.Replace(tpl, "y", "2006", -1)
	tpl = strings.Replace(tpl, "m", "01", -1)
	tpl = strings.Replace(tpl, "d", "02", -1)
	tpl = strings.Replace(tpl, "h", "15", -1)
	tpl = strings.Replace(tpl, "i", "04", -1)
	tpl = strings.Replace(tpl, "s", "05", -1)
	return time.Unix(t, 0).Format(tpl)
}

// GetMonthRange 获得指定年份和月份的起始unix时间和截止unix时间
func GetMonthRange(year int, month int) (beginTime, endTime int64, err error) {
	// 获得当前时间
	t := time.Now()

	if year == 0 {
		year = t.Year()
	}

	if month == 0 {
		month = int(t.Month())
	}
	yearStr := strconv.Itoa(year)
	monthStr := strconv.Itoa(month)

	// 拼接当时时间的字符串格式
	str := yearStr + "-" + monthStr + "-1 00:00:00"

	// 起始时间
	begin, err := time.ParseInLocation(timeTpl, str, time.Local)
	if err != nil {
		return
	}
	beginTime = begin.Unix()
	month = int(begin.Month())
	day := 30
	if month == 2 {
		day = 28
	} else if month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12 {
		day = 31
	}

	// 截止时间
	str = yearStr + "-" + monthStr + "-" + strconv.Itoa(day) + " 23:59:59"
	end, err := time.ParseInLocation(timeTpl, str, time.Local)
	if err != nil {
		return
	}
	endTime = end.Unix()
	return
}

// GetWeek 获得星期的数字
func GetWeek(t time.Time) int {
	week := t.Weekday().String()
	switch week {
	case "Monday":
		return 1
	case "Tuesday":
		return 2
	case "Wednesday":
		return 3
	case "Thursday":
		return 4
	case "Friday":
		return 5
	case "Saturday":
		return 6
	case "Sunday":
		return 7
	}
	return 0
}
