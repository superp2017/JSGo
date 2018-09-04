package JsUser

//用户
import (
	. "JsGo/JsLogger"
	"JsGo/JsStore/JsRedis"
	"JsGo/JsUuid"
	"errors"
	"fmt"
)

type Wechat struct {
	Openid       string //公众号
	Openid_small string //小程序
	Openid_app   string //App
	Unionid      string //unionID
}

type RecAddr struct {
	Alias    string //收货信息别名
	Name     string //收货人姓名
	Cell     string //收货人电话
	Province string //省
	City     string //市
	Area     string //区
	Addr     string //收货地址
}

//用户摘要信息ABS//（用于商家获取用户列表）
type UserABS struct { //abstract摘要
	ID       string //uid
	Nickname string //昵称
	Header   string //头像
	Sex      int    //性别
	Mobile   string //手机号
	Address  string //地址
}

type User struct {
	ID             string    //uid
	Nickname       string    //昵称
	Header         string    //头像
	Sex            int       //性别
	Birthday       string    //生日
	Mobile         string    //手机号
	Role           string    //角色
	Address        string    //地址
	City           string    //城市
	Country        string    //国家
	Profession     string    //职业
	Tag            string    //标签
	Email          string    //邮箱
	Account        string    //账户
	Password       string    //密码
	CreatTime      string    //创建时间
	ReceivingAddrs []RecAddr //收货地址列表
	Wechat                   //微信
	MessTime       string    //（查看）消息时间
}

//interface
type INewUser interface {
	NewUser(user *User) error
	ModUser(user *User) error
	DelUser(ID string) error
	QueryUser(ID string) (*User, error)
}

func (user *User) NewUser(db string) error {
	id := JsUuid.NewV4().String()
	for b, e := JsRedis.Redis_hexists(db, id); b; {
		id = JsUuid.NewV4().String()
		if e != nil {
			break
		}
	}
	user.ID = id
	return JsRedis.Redis_hset(db, user.ID, user)
}

func (user *User) ModUser(db string) error {
	b, e := JsRedis.Redis_hexists(db, user.ID)
	if e != nil {
		Error(e.Error())
		return e
	}

	if b {
		JsRedis.Redis_hset(db, user.ID, user)
		return nil
	} else {
		return errors.New(fmt.Sprintf("user[%s] not exist", user.ID))
	}
}

func (user *User) DelUser(db string) error {
	return errors.New(fmt.Sprintf("delete user[%s] not support", user.ID))
}

func (user *User) QueryUser(db string) error {
	e := JsRedis.Redis_hget(db, user.ID, user)
	if e != nil {
		Error(e.Error())
		return e
	}
	return e
}
