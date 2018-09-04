package JsRedis

//哈希表，自选选DB块(0---15) //db
import (
	. "JsGo/JsLogger"
	"encoding/json"
	"errors"

	"github.com/gomodule/redigo/redis"
)

//读取数据库
func Redis_hdbget(db int8, t, k string, v interface{}) error {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return errors.New("connect is nil")
	}
	c.Send("SELECT", db)
	c.Send("HGET", t, k)
	c.Flush()
	c.Receive()
	b, e := redis.Bytes(c.Receive())
	if e != nil {
		Error("hdbget hdbsize DB【%s】 Table【%s】,Key【%s】，error:%v", db, t, k, e.Error())
		return e
	}
	e = json.Unmarshal(b, v)
	if e != nil {
		Error(e.Error())
		return e
	}
	return nil
}

//写入数据库
func Redis_hdbset(db int8, t, k string, v interface{}) error {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return errors.New("connect is nil")
	}

	b, e := json.Marshal(v)
	if e != nil {
		Error(e.Error())
		return e
	}

	c.Send("SELECT", db)
	c.Send("HSET", t, k, b)
	c.Flush()
	c.Receive()
	_, e = c.Receive()
	if e != nil {
		Error("hdbset hdbsize DB【%s】 Table【%s】,Key【%s】，error:%v", db, t, k, e.Error())
		return e
	}
	return nil
}

//从数据库查找是否存在
func Redis_hdbexists(db int8, t, k string) (bool, error) {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return false, errors.New("connect is nil")
	}

	c.Send("SELECT", db)
	c.Send("HEXISTS", t, k)
	c.Flush()
	c.Receive()
	b, e := redis.Bool(c.Receive())
	if e != nil {
		Error("hdbexists hdbsize DB【%s】 Table【%s】,Key【%s】，error:%v", db, t, k, e.Error())
		return false, e
	}

	return b, nil
}

func Redis_hdbsize(db int8, t string) (int64, error) {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return -1, errors.New("connect is nil")
	}

	c.Send("SELECT", db)
	c.Send("HLEN", t)
	c.Flush()
	c.Receive()
	b, e := redis.Int64(c.Receive())
	if e != nil {
		Error("hdbsize DB【%s】 Table【%s】,error:%v", db, t, e.Error())
		return -1, e
	}

	return b, nil
}

func Redis_hdbkeys(db int8, t string) ([]string, error) {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return nil, errors.New("connect is nil")
	}

	c.Send("SELECT", db)
	c.Send("HKEYS", t)
	c.Flush()
	c.Receive()
	b, e := redis.Strings(c.Receive())
	if e != nil {
		Error("hdbkeys DB【%s】 Table【%s】,error:%v", db, t, e.Error())
		return nil, e
	}

	return b, nil
}

func Redis_hdbdel(db int8, t, k string) error {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return errors.New("connect is nil")
	}

	c.Send("SELECT", db)
	c.Send("HDEL", t, k)
	c.Flush()
	c.Receive()
	_, e := c.Receive()
	if e != nil {
		Error("hdbdel DB【%s】 Table【%s】,error:%v", db, t, e.Error())
		return e
	}

	return nil
}

func Redis_hdbmset(db int8, t string, data map[string]interface{}) error {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return errors.New("connect is nil")
	}

	c.Send("SELECT", db)
	for k, v := range data {
		b, e := json.Marshal(v)
		if e != nil {
			Error(e.Error())
			return e
		}
		// fmt.Printf("b = %v, %s\n", v, string(b))
		c.Send("HSET", t, k, b)
	}

	c.Flush()
	c.Receive()
	for k, _ := range data {
		_, e := c.Receive()
		if e != nil {
			Error("k[%s],%s", k, e.Error())
			return e
		}
	}

	return nil
}

func Redis_hdbmget(db int8, t string, ret *map[string]interface{}) error {
	c := g_pool.Get()
	if c == nil {
		Error("connect is nil")
		return errors.New("connect is nil")
	}

	c.Send("SELECT", db)
	keys := make([]string, len(*ret))
	i := 0
	for k, _ := range *ret {
		c.Send("HGET", t, k)
		keys[i] = k
		i++
	}

	c.Flush()
	c.Receive()
	for _, k := range keys {

		b, e := redis.Bytes(c.Receive())
		if e != nil {
			Error(e.Error())
			return e
		}

		e = json.Unmarshal(b, (*ret)[k])
		if e != nil {
			Error(e.Error())
			return e
		}
	}

	return nil
}
