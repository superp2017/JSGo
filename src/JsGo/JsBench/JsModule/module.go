package JsModule

type Timelimit struct {
	Title     string //标题
	StartTime int64  //开始时间
	StartDate string //开始时间
	EndTime   int64  //结束时间
	EndDate   string //结束时间
	Pic       string //宣传图片
	//timer     *time.Timer //定时器
}

type Para struct {
	PageID string //版面号
}

type ModuleDetail struct {
	UID           string   //用户或作者的UID
	UserHead      string   //用户或作者的头像
	UserName      string   //用户或作者的姓名
	ID            string   //产品或者内容id
	Type          string   //0:产品 1: 内容
	Price         int      //价格
	OrigPrice     int      //原价
	Name          string   //标题
	SubName       string   //副标题
	Pic           string   //展示的图片
	Tags          []string //标签
	VisitNum      int      //访问量
	PraiseNum     int      //点赞数量
	AttentionNums int      //关注数量
	CommentNum    int      //评论的人数
	RelationID    string   //关联的id（关联的产品或者关联的内容）
	ExData        string   //扩展字段
}

type Module struct {
	ModuleID  string         //模块id
	Name      string         //模块名字
	Title     string         //模块主标题
	SubTitle  string         //模块副标题
	PageID    string         //版面号
	Data      []ModuleDetail //关联的内容或者产品
	Islimit   bool           //是否是限时活动
	Limit     *Timelimit     //限时结构
	Status    string         //状态
	CreatTime string         //创建时间
}
