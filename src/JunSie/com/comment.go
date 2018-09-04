package com

//评论（统计评论数量）
import (
	"JsGo/JsBench/JsComments"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"JunSie/util"
	"fmt"
)

type XMComment struct {
	JsComments.Comments
}

func Init_comment() {
	JsHttp.WhiteHttps("/newcomment", NewComment)              //新建评论
	JsHttp.WhiteHttp("/getrelatedcomment", GetRelatedComment) //获取评论
}

func Init_commentMall() {
	JsHttp.WhiteHttps("/newcomment", NewComment)               //新建评论
	JsHttp.WhiteHttps("/getrelatedcomment", GetRelatedComment) //获取评论
}

//新建评论
func NewComment(s *JsHttp.Session) {
	type Para struct {
		RelaID string    //相关id，产品id或者内容id
		Com    XMComment //评论结构
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.RelaID == "" || st.Com.UID == "" {
		info := fmt.Sprintf("newComment failed,UID=%s,RelaID=%s\n", st.Com.UID, st.RelaID)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}
	st.Com.ID = util.IDer(constant.COMMENT)
	st.Com.TimeStamp = util.CurStamp()
	st.Com.CreatTmie = util.CurTime()
	if st.Com.Score < 0 {
		st.Com.Score = 0
	}
	if st.Com.Score > 5 {
		st.Com.Score = 5
	}
	if err := JsRedis.Redis_hget(constant.Rela_COMMENT, st.RelaID, &st.Com.Next); err != nil {
		JsLogger.Info(err.Error()) //找到产品或文章关联的最新的放在新评论关联的下一个评论
	}
	st.Com.RealID = st.RelaID
	if err := JsRedis.Redis_hset(constant.COMMENT, st.Com.ID, &st.Com); err != nil {
		JsLogger.Error(err.Error()) //用新生成的评论ID为K保存到评论表中。
		s.Forward("1", err.Error(), nil)
		return
	}

	go JsRedis.Redis_hset(constant.Rela_COMMENT, st.RelaID, st.Com.ID)
	//将产品或内容ID为K，保存当前最新评论ID
	if st.Com.Type == "0" {
		go NewProComment(st.RelaID, float64(st.Com.Score*1.0))
	} else {
		go NewContentComment(st.RelaID, float64(st.Com.Score*1.0))
	}
	s.Forward("0", "success", st.Com)
}

//获取多个关联的评论
func GetRelatedComment(s *JsHttp.Session) {
	type Para struct {
		RelaID   string //关联的id，比如产品id或者内容id
		CurComID string //当前最后一个评论id
		Nums     int    //需要的数量
	}
	st := &Para{}
	if err := s.GetPara(st); err != nil {
		s.Forward("1", err.Error(), nil)
		return
	}
	if st.Nums <= 0 || (st.RelaID == "" && st.CurComID == "") {
		info := fmt.Sprintf("GetRelatedComment,RelaID=%s,CurComment=%s,Nums=%d\n", st.RelaID, st.CurComID, st.Nums)
		JsLogger.Error(info)
		s.Forward("1", info, nil)
		return
	}

	if st.CurComID == "" {
		if err := JsRedis.Redis_hget(constant.Rela_COMMENT, st.RelaID, &st.CurComID); err != nil {
			s.Forward("1", err.Error(), nil)
			return
		}
		if st.CurComID == "" {
			info := fmt.Sprintf("GetRelatedComment,RelaID=%s has no commnet", st.RelaID)
			JsLogger.Error(info)
			s.Forward("0", info, nil)
			return
		}
	} else {
		d := &XMComment{}
		if err := JsRedis.Redis_hget(constant.COMMENT, st.CurComID, d); err != nil {
			JsLogger.Error(err.Error())
			s.Forward("0", err.Error(), nil)
			return
		}
		if d.Next == "" {
			s.Forward("0", "not  more comments", nil)
			return
		}
		st.CurComID = d.Next
	}

	data := []*XMComment{}
	for i := 0; i < st.Nums; i++ {
		d := &XMComment{}
		if err := JsRedis.Redis_hget(constant.COMMENT, st.CurComID, d); err != nil {
			JsLogger.Error(err.Error())
			break
		}
		data = append(data, d)
		if d.Next == "" {
			break
		}
		st.CurComID = d.Next
	}

	s.Forward("0", "success", data)
}
