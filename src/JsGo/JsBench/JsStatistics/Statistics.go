package JsStatistics

import (
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
)

type Statistics struct {
	Type           string   // 类型 0:产品 1:内容
	VisitNum       int      // 访问量
	PraiseNum      int      // 点赞数量
	AttentionNums  int      // 关注数量
	CompositeScore float64  // 综合评分（所有评分的平均值）
	CommentNum     int      // 评论的人数
	SalesVolume    int      //销量
	PraiseUser     []string // 点赞人前十个小头像
	AttentiUser    []string // 关注人前十个小头像
}

func (this *Statistics) Init() {
	this.Type = ""
	this.VisitNum = 0
	this.PraiseNum = 0
	this.AttentionNums = 0
	this.CompositeScore = 5
	this.CommentNum = 0
	this.SalesVolume = 0
	this.PraiseUser = []string{}
	this.AttentiUser = []string{}
}

//新增访问量
func (this *Statistics) NewVisit() {
	this.VisitNum++
}

//新增加点赞量
func (this *Statistics) NewPraise(uid string) {
	ok := false
	for _, v := range this.PraiseUser {
		if v == uid {
			ok = true
			break
		}
	}
	if !ok {
		this.PraiseUser = append(this.PraiseUser, uid)
		this.PraiseNum++
	}
}

//新增关注
func (this *Statistics) NewAttention(uid string) {
	ok := false
	for _, v := range this.AttentiUser {
		if v == uid {
			ok = true
			break
		}
	}
	if !ok {
		this.AttentiUser = append(this.AttentiUser, uid)
		this.AttentionNums++
	}
}

//新增评论
func (this *Statistics) NewComment(score float64) {
	if score < 0 {
		score = 0
	}
	d := this.CompositeScore * (float64)(this.CommentNum*1.0)
	this.CommentNum++
	this.CompositeScore = (d + score) / (float64)(this.CommentNum*1.0)
}

//取消关注
func (this *Statistics) RemoveAttention(uid string) {
	index := -1
	for i, v := range this.AttentiUser {
		if v == uid {
			index = i
			break
		}
	}
	if index != -1 {
		this.AttentiUser = append(this.AttentiUser[:index], this.AttentiUser[index+1:]...)
		this.AttentionNums--
	}
}

//取消点赞
func (this *Statistics) RemovePraise(uid string) {
	index := -1
	for i, v := range this.PraiseUser {
		if v == uid {
			index = i
			break
		}
	}
	if index != -1 {
		this.PraiseUser = append(this.PraiseUser[:index], this.PraiseUser[index+1:]...)
		this.PraiseNum--
	}
}

//新增访问量
func NewVisit(db, ID, Type string) (*Statistics, error) {
	data := &Statistics{}
	if err := JsRedis.Redis_hget(db, ID, data); err != nil {
		JsLogger.Info(err.Error())
		data.Init()
	}
	data.Type = Type
	data.NewVisit()
	return data, JsRedis.Redis_hset(db, ID, data)
}

//新增加点赞量
func NewPraise(db, ID, Type, uid string) (*Statistics, error) {
	data := &Statistics{}
	if err := JsRedis.Redis_hget(db, ID, data); err != nil {
		JsLogger.Info(err.Error())
		data.Init()
	}
	data.Type = Type
	data.NewPraise(uid)
	return data, JsRedis.Redis_hset(db, ID, data)
}

//新增关注
func NewAttention(db, ID, Type, uid string) (*Statistics, error) {
	data := &Statistics{}
	if err := JsRedis.Redis_hget(db, ID, data); err != nil {
		JsLogger.Info(err.Error())
		data.Init()
	}
	data.Type = Type
	data.NewAttention(uid)
	return data, JsRedis.Redis_hset(db, ID, data)
}

//新增评论
func NewComment(db, ID, Type string, score float64) (*Statistics, error) {
	data := &Statistics{}
	if err := JsRedis.Redis_hget(db, ID, data); err != nil {
		JsLogger.Info(err.Error())
		data.Init()
	}
	data.Type = Type
	data.NewComment(score)
	return data, JsRedis.Redis_hset(db, ID, data)
}

//取消关注
func RemoveAttention(db, ID, Type, uid string) (*Statistics, error) {
	data := &Statistics{}
	if err := JsRedis.Redis_hget(db, ID, data); err != nil {
		JsLogger.Info(err.Error())
		data.Init()
	}
	data.Type = Type
	data.RemoveAttention(uid)
	return data, JsRedis.Redis_hset(db, ID, data)
}

//取消点赞
func RemovePraise(db, ID, Type, uid string) (*Statistics, error) {
	data := &Statistics{}
	if err := JsRedis.Redis_hget(db, ID, data); err != nil {
		JsLogger.Info(err.Error())
		data.Init()
	}
	data.Type = Type
	data.RemovePraise(uid)
	return data, JsRedis.Redis_hset(db, ID, data)
}

//增加销售量SalesVolume
func NewSales(db, ID, Type string, num int) (*Statistics, error) {
	data := &Statistics{}
	if err := JsRedis.Redis_hget(db, ID, data); err != nil {
		JsLogger.Info(err.Error())
		data.Init()
	}
	data.Type = Type
	data.SalesVolume += num
	return data, JsRedis.Redis_hset(db, ID, data)
}
