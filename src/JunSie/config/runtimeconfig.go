package config

import (
	"JsGo/JsConfig"
	"JsGo/JsHttp"
	"JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JunSie/constant"
	"log"
)

var config *JsConfig.RunTimeConfig

func InitConfig(Fun func()) {
	var err error
	config, err = getRunTimeConfig("WeChat")
	if err != nil {
		log.Fatalln("加载系统配置失败！")
		return
	}
	JsConfig.InitConfig(config)

	JsHttp.WhiteHttp("/updateruntimeconfig", UpdateRunTimeConfig) //更新配置
	JsHttp.WhiteHttp("/delruntimeconfig", DelRunTimeConfig)       //删除配置
	JsHttp.WhiteHttp("/getruntimeconfig", GetRunTimeConfig)       //获取配置

	if Fun != nil {
		Fun()
	}

}

//更新配置
func UpdateRunTimeConfig(session *JsHttp.Session) {
	type Para struct {
		BID    string
		Config JsConfig.RunTimeConfig
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if err := updateRunTimeConfig(st.BID, &st.Config); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "success", st.Config)
}

//删除配置
func DelRunTimeConfig(session *JsHttp.Session) {
	type Para struct {
		BID string
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.BID == "" {
		session.Forward("1", "DelRunTimeConfig:BID 为空\n", nil)
		return
	}
	if err := delRunTimeConfig(st.BID); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "DelRunTimeConfig success\n", nil)
}

//获取配置
func GetRunTimeConfig(session *JsHttp.Session) {
	type Para struct {
		BID string
	}
	st := &Para{}
	if err := session.GetPara(st); err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	if st.BID == "" {
		session.Forward("1", "GetRunTimeConfig:BID 为空\n", nil)
		return
	}
	config, err := getRunTimeConfig(st.BID)
	if err != nil {
		JsLogger.Error(err.Error())
		session.Forward("1", err.Error(), nil)
		return
	}
	session.Forward("0", "GetRunTimeConfig success\n", config)
}

//更新配置
func updateRunTimeConfig(BID string, config *JsConfig.RunTimeConfig) error {
	if BID == "" {
		return JsLogger.ErrorLog("updateCompanyConfig faild,商家ID 为空\n")
	}
	if err := checkConfig(config); err != nil {
		return err
	}
	return JsRedis.Redis_hset(constant.H_RunTimeConfig, BID, config)
}

//获取运行配置
func getRunTimeConfig(BID string) (*JsConfig.RunTimeConfig, error) {
	config := &JsConfig.RunTimeConfig{}
	e := JsRedis.Redis_hget(constant.H_RunTimeConfig, BID, config)
	return config, e
}

//删除运行配置
func delRunTimeConfig(BID string) error {
	return JsRedis.Redis_hdel(constant.H_RunTimeConfig, BID)
}

func checkConfig(config *JsConfig.RunTimeConfig) error {
	return nil
}
