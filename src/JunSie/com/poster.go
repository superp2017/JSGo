package com

import (
	"JsGo/JsNet"
)

//海报创建、修改、上下线、删除、查询
type Poster struct {
	PosterID  string
	Status    string
	Pic       string
	Title     string
	ArticleID string
}

func InitPoster() {
	//JsNet.Https("/createposter", createPoster)
}

func createPoster(s *JsNet.StSession) {

}

func modifyPoster(s *JsNet.StSession) {

}

func deletePoster(s *JsNet.StSession) {

}

func queryPoster(s *JsNet.StSession) {

}
