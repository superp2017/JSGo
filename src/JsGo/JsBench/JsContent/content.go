package JsContent

type CFeature struct {
	Type  string //类型
	Text  string //文本
	Image string //图片
	ExID  string //外部id
	ProName string //关联产品名称
}

type Contents struct {
	UID       string     //内容创建者
	UserHead  string     //用户头像
	ID        string     //id
	Title     string     //标题
	SubTitle  string     //副标题
	Brief     string     //简介
	Author    string     //作者
	TimeStamp string     //时间戳
	Images    []string   //图片列表
	Thumbnail string     //缩略图
	Content   []CFeature //正文
	Products  []string   //关联的产品id列表
	DelTag    bool       //是否删除
	Status    string     //状态
	CreatDate string     //创建时间
}

type ContentsAbs struct { //abstract,摘要，拉取列表时使用
	UID       string   //内容创建者
	UserHead  string   //用户头像
	ID        string   //id
	Title     string   //标题
	Brief     string   //简介
	Author    string   //作者
	CreatDate string   //创建时间
	Images    []string //图片列表
	Thumbnail string   //缩略图
	DelTag    bool     //是否删除
	Status    string   //状态
}
