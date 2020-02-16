package cache

import (
	"bookms/utils"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"os"
)

const zsetopFileName  = "cache/zsetop.lua"

var (
	pool *redis.Pool
	zsetopScript string
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
	file, err := os.Open(zsetopFileName)
	if err != nil {
		panic("open"+zsetopFileName+" "+err.Error())
	}
	bytes,err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}
	zsetopScript = utils.UnsafeBytesToString(bytes)
	logs.Debug(zsetopScript)
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

func ListPush(key string, value interface{}, expire ...int) (err error, reply interface{}) {
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
		reply, err = c.Do("LPUSH", key, v,"EX", expire[0])
	}else {
		reply, err = c.Do("LPUSH", key, v)
	}
	return
}

func ListPop(key string) (err error, reply interface{}) {
	c := pool.Get()
	reply, err = c.Do("RPOP", key)
	return
}

func ListLen(key string) (err error, reply interface{}) {
	c := pool.Get()
	reply,err = c.Do("LLEN", key)
	return
}

func ZaddWithCap2(key,member string, score float32, maxScore, cap int) (reply interface{}, err error) {
	src := "local key,cap,maxSetScore,newMemberScore,member = KEYS[1],ARGV[1],ARGV[2],ARGV[3],ARGV[4] "+
		"redis.log(redis.LOG_NOTICE, 'key=', key,',cap=', cap,',maxSetScore=', maxSetScore,',newMemberScore=', newMemberScore,',member=', member) "+
		"local len = redis.call('zcard', key) "+
		"if len then "+
		"if tonumber(len) >= tonumber(cap) then "+
		"local num = tonumber(len)-tonumber(cap)+1 "+
		"local list = redis.call('zrangebyscore',key,0,maxSetScore,'limit',0,num) "+
		"redis.log(redis.LOG_NOTICE,'key=',key,'maxSetScore=',maxSetScore, 'num=',num) "+
		"for k,lowestScoreMember in pairs(list) do "+
		"local lowestScore = redis.call('zscore', key,lowestScoreMember) "+
		"redis.log(redis.LOG_NOTICE, 'list: ', lowestScore, lowestScoreMember) "+
		"if tonumber(newMemberScore) > tonumber(lowestScore) then "+
		"local rank = redis.call('zrevrank',key,member) "+
		"if not rank then "+
		"local index = tonumber(len) - tonumber(cap) "+
		"redis.call('zremrangebyrank',key, 0, index) "+
		"end "+
		"redis.call('zadd', key, newMemberScore, member) "+
		"break "+
		"end "+
		"end "+
		"else "+
		"redis.call('zadd', key, newMemberScore, member) "+
		"end "+
		"end"
	logs.Debug(src)
	lua := redis.NewScript(1, src)
	c := pool.Get()
	logs.Debug("ZaddWithCap:",key, cap,maxScore,score,member)
	reply,err = lua.Do(c, key, cap, maxScore, score, member)
	return
}

func ZaddWithCap(key,member string, score float32, maxScore, cap int) (reply interface{}, err error) {
	logs.Debug("ZaddWithCap: ", key, cap,maxScore,score,member)
	c := pool.Get()
	//eval zsetop.lua mtest , 3 5 5 test1
	//reply, err = c.Do("eval",zsetopScript,1,key,cap,maxScore,score,member)
	lua :=redis.NewScript(1,zsetopScript)
	reply, err= lua.Do(c,key,cap,maxScore,score,member)
	return
}

func ZrevRangeByScore(key string, max,min, start, stop int) (res []string, err error) {
	c := pool.Get()
	reply,err := redis.Values(c.Do("zrevrangebyscore", key, max, min,"limit", start, stop))
	//reply,err := redis.Values(c.Do("zrangebyscore", key, min, max,"limit", start, stop))
	if err != nil {
		return
	} else {
		for _, r := range reply {
			logs.Debug("ZrevRangeByScore: ", utils.UnsafeBytesToString(r.([]byte)))
			res = append(res, utils.UnsafeBytesToString(r.([]byte)))
		}
		logs.Debug("ZrevRangeByScore: ", res)
	}
	return
}