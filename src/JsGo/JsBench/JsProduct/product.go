package JsProduct

import (
	"JsGo/JsBench/JsComments"
	"JsGo/JsBench/JsEvaluate"
)

type Color struct {
	Color string
	Text  string
}

type ProductFormat struct { //产品规格
	Format    string //规格说明
	Pic       string //规格对应的图片
	Price     int    //规格对应的价格
	Inventory int    //规格库存
	//
}

type Product struct {
	UID        string              //产品创建者
	ID         string              //产品ID
	Title      string              //标题
	SubTitle   string              //副标题
	Label      []string            //标签
	Color      []string            //颜色
	Thumbnail  string              //缩略图
	Images     []string            //海报
	Pic        []string            //详情介绍图
	ProFormat  []ProductFormat     //产品规格
	FormatList map[string][]string //规格列表
	Desc       string              //描述
	OriPrice   int                 //原件
	NowPrice   int                 //现价
	Freight    int                 //运费
	Type       string              //类型
	DelTag     bool                //是否删除
	Status     string              //状态
	CommentID  string              //评论ID
	ContentID  string              //文章id
	EvaluateId string              //评价id
	CreatDate  string              //创建时间
}

type ProductAbs struct { //abstract摘要
	ID        string   //产品ID
	Title     string   //标题
	Color     []string //颜色
	Thumbnail string   //缩略图
	Images    []string //海报
	Desc      string   //描述
	OriPrice  int      //原件
	NowPrice  int      //现价
	Type      string   //类型
	DelTag    bool     //是否删除
	Status    string   //状态
	CreatDate string   //创建时间
}

type IProduct interface {
	AddProduct(p *Product) error
	ModProduct(p *Product) error
	QueryProduct(PID string) (*Product, error)
	DelProduct(PID string) error
	VisitProduct(PID string) error
	EvaluateProduct(PID string, eva *JsEvaluate.Evaluate) error
	CommentProduct(PID string, c *JsComments.Comments) error
}
