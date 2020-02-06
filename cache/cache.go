package cache

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func init() {
	server := beego.AppConfig.String("server")
	password := beego.AppConfig.String("password")
	db := beego.AppConfig.String("db")
	max_active_connection, err :=beego.AppConfig.Int("max_active_connection")
	if err != nil {
		panic(err.Error())
	}
	max_idle_connection,err  :=beego.AppConfig.Int("max_idle_connection")
	if err != nil {
		panic(err.Error())
	}
	pool = &redis.Pool{
		MaxActive: max_active_connection,
		MaxIdle: max_idle_connection,
		// Other pool configuration not shown in this example.
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func ClosePool() (err error) {
	err = pool.Close()
	return
}

func SetString(key, value string, expire ...int) (err error, reply interface{}) {
	if len(expire) > 1 {
		err = errors.New("Too many arguments")
		return
	}
	c := pool.Get()
	if len(expire) == 1 {
		reply, err = c.Do("SET", key, value,"EX", expire[0])
	} else {
		reply,err = c.Do("SET", key, value)
	}
	return
}

func GetString(key string) (err error, res string) {
	c := pool.Get()
	res,err = redis.String(c.Do("GET", key))
	return
}

func SetInt(key string, value int, expire ...int) (err error, reply interface{}) {
	if len(expire) > 1 {
		err = errors.New("Too many arguments")
	}
	c := pool.Get()
	if 1 == len(expire) {
		reply, err = c.Do("SET",key,value,"EX",expire[0])
	}else {
		reply,err = c.Do("SET",key,value)
	}
	return
}

func GetInt(key string) (err error, res int) {
	c := pool.Get()
	res, err = redis.Int(c.Do("GET", key))
	return
}

func SetInterface(key string, value interface{}, expire ...int) (err error, reply interface{}) {
	if len(expire) > 1 {
		err = errors.New("Too many arguments")
		return
	}
	c := pool.Get()
	v, err := json.Marshal(value)
	if err != nil {
		return
	}
	if 1 == len(expire) {
		reply, err = c.Do("SET", key, v,"EX", expire[0])
	}else {
		reply, err = c.Do("SET", key, v)
	}
	return
}

func GetInterface(key string, res interface{}) (err error) {
	c := pool.Get()
	reply, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		return
	}
	err = json.Unmarshal(reply, res)
	return
}

func SetExpire(key string, expire int) (err error, reply interface{}) {
	c := pool.Get()
	reply, err = c.Do("EXPIRE", key, expire)
	return
}