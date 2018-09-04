package JsHttp

import "time"

type ServerTime struct {
	Year      int    //年
	Month     int    //月
	Day       int    //日
	Hour      int    //时
	Minute    int    //分
	Second    int    //秒
	Timestamp int64  //时间戳
	Date      string //日期：如2018-07-30
}

//获取系统时间
func GetServerDateTime(s *Session) {
	data := ServerTime{}
	time := time.Now()
	data.Year = time.Year()
	data.Month = (int)(time.Month())
	data.Day = time.Day()
	data.Hour = time.Hour()
	data.Minute = time.Minute()
	data.Second = time.Second()
	data.Timestamp = time.Unix()
	data.Date = time.Format("2006-01-02")
	s.Forward("0", "GetServerDateTime：Success\n", data)
}
