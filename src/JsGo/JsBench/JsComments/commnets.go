package JsComments

//评论
type Comments struct {
	ID        string   //评论id
	UID       string   //用户ID
	UserHead  string   //用户头像
	UserName  string   //用户姓名
	Type      string   //类型
	RealID    string   //被评论的id
	Text      string   //文字
	TimeStamp int64    //时间
	Score     int      //星星等级（1-5）
	CreatTmie string   //评论时间
	Images    []string //图片
	Video     string   //录像
	ThumbsUp  int      //点赞数量
	Next      string   //下一个评论id
}

type IComments interface {
	NewComment() error
	PosiComment() error
	NegaComment() error
}

//评论分类功能（好评、差评、中评、有图、）
//五角星数量
//用户头像以及昵称加***
//给评论点赞数量,其他人平论数量。
//其他人的评论（其他人提问或者发表观点）
//商家的回复
//
