package JsHotword

////////////////////
//目前最多收录前3000个
//算法描述：按搜索次数排名，搜索次数相同的按最新时间戳排名，在2000名之外的，只出现一次的的热词，根据时间戳最久的最先删除，如：超过一周的即删除

type HotWord struct {
	Hotword string
	Count   int
	Stamp   int64
}

//追加热词
func AppendHotword(w string) error {
	return nil
}

//查询前N个热词
func QueryNHotword(n int) []*HotWord {
	return nil
}
