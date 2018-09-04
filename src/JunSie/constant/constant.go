package constant

//表名
const (
	UNIONUID        = "H_UNIONUID"
	OPENUID         = "H_OPENUID"
	Key_Mobile      = "H_Mobile"
	USER            = "H_USER"
	ADMIN           = "H_ADMIN"
	TAG             = "H_TAG"
	TAGPRO          = "H_Tag_Pro" //标签产品关联
	TAGART          = "H_tag_Art" //标签文章关联
	COMMENT         = "H_COMMENT"
	Rela_COMMENT    = "Re_COMMENT"      //评论
	H_Order         = "H_Order"         //订单
	H_UserOrder     = "H_UserOrder"     //用户订单
	SHOPPINGCART    = "ShoppingCart"    //购物车
	H_RunTimeConfig = "RunTimeConfig"   //运行配置
	H_ApplyBusiness = "H_ApplyBusiness" //商家申请
	H_Business      = "H_Business"      //商家
)

///产品
const (
	DB_Product        = "Product"       //产品数据库名称
	DB_ProDraft       = "ProDraft"      //产品草稿
	KEY_Globalproduct = "GlobalProduct" //全局产品数据库名称
	Status_ON         = "Status_ON"     //上架状态
	Status_OFF        = "Status_OFF"    //下架状态
	ProDraft_SIZE     = 10              //草稿最大数量
	ProStatistics     = "ProStatistics" //产品统计表 -----
	Key_WaitPay       = "WaitPay"       //待支付的订单列表
)

const (
	DB_Content        = "Content"           // 内容
	DB_ContentDraft   = "ContentDraft"      // 内容草稿
	ContentDraft_SIZE = 10                  //草稿最大数量
	KEY_GlobalContent = "GlobalContent"     //全局产品数据库名称
	ContentStatistics = "ContentStatistics" //产品统计表
)

///首页
const (
	POSTER     = "Poster"     //首页海报
	FRESHPRO   = "FreshPro"   //新品首发（6）
	POPULARITY = "Popularity" //人气推荐（4）受大众欢迎
	FLASHSALE  = "FlashSale"  //限时抢购（3）
	GUESSSLIKE = "GuessLike"  //guessLike猜你喜欢
)

const (
	HomeModule       = "HomeModule"   //首页模块
	ModuleMap        = "ModuleMap"    //模块映射
	KEY_ModuleConfig = "ModuleConfig" //模块列表配置
)

//板块
const (
	SHOWARTONE = "ShowArticleOne" //展示板块一
)

const ( //个人中心记录（私人记录）表名字
	MYLikeArt       = "MyLikeArt"       //2.01我赞过的文章
	MYLikePro       = "MyLikePro"       //2.02我赞过的产品
	MYCollectionArt = "MyCollectionArt" //3.01我收藏的文章
	MYCollectionPro = "MyCollectionPro" //3.02我收藏的产品
	MYAttention     = "MyAttention"     //3.1记录关注（作者）
	MYComment       = "MyComment"       //4.记录(评论过的)(X)
	MYPageViewArt   = "MyPageViewArt"   //1.我浏览过的产品
	MYPageViewPro   = "MyPageViewPro"   //1.我浏览过的文章
)

const (
	OrderStatus_Creat  = "OrderStatus_Creat"  //订单创建
	OrderStatus_Submit = "OrderStatus_Submit" //订单提交

	OrderStatus_Paid = "OrderStatus_Paid" //订单支付
	OrderStatus_Send = "OrderStatus_Send" //订单发货

	OrderStatus_Receive = "OrderStatus_Receive" //订单收货

	OrderStatus_Success = "OrderStatusSuccess" //订单评价完成

	OrderStatus_Cancel = "OrderStatus_Cancel" //订单取消

)
const (
	MessageReadyBus_001 = "MessageReadyBus_001" //一号商家消息准备草稿
	MessageBus_001      = "MessageBus_001"      //一号商家推送消息
	Chat                = "H_Chat"              //留言，用户消息（客服）
	ChatList            = "S_ChartList"         //商家消息队列（ID列表）
	WaitOrder           = "S_WaitOrder"         //待处理订单
	RepertoryLess       = "S_RepertoryLess"     //库存不足产品
	Statistics          = "Statistics"          //网站统计数据
	HotProList          = "HotProList"          //热门产品
	SalesTrend          = "H_SalesTrend"        //销售趋势
	VisitNumber         = "S_VisitNumber"       //当日访问用户
)
