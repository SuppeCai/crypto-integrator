package exchange

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
	"encoding/json"
)

var (
	host      = "localhostd:6379"
	pwd       = "password"
	Separator = ":"
)

type Cache struct {
	pool *redis.Pool
}

func (cache *Cache) Init() {
	cache.pool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				LogErr.Error("redis conn error:" + err.Error())
				return nil, err
			}
			if _, err := c.Do("AUTH", pwd); err != nil {
				LogErr.Error("redis auth error:" + err.Error())
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (cache *Cache) Get(key string, value interface{}) {
	conn := cache.pool.Get()
	defer conn.Close()

	v, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		LogErr.Error("redis get error:" + err.Error())
		return
	}

	err = json.Unmarshal(v, &value)
	if err != nil {
		LogErr.Error("redis json error:" + err.Error())
		return
	}
	return
}

func (cache *Cache) Put(key string, value interface{}, ttl int) string{
	conn := cache.pool.Get()
	defer conn.Close()

	v, err := json.Marshal(value)
	if err != nil {
		LogErr.Error("redis marshal error:" + err.Error())
		return ""
	}

	reply, err := conn.Do("SET", key, v)
	if err != nil {
		fmt.Println(reply)
		LogErr.Error("redis set error:" + err.Error())
		return ""
	}

	reply, err = conn.Do("EXPIRE", key, ttl)
	if err != nil {
		fmt.Println(reply)
		LogErr.Error("redis marshal error:" + err.Error())
		return ""
	}
	return string(v)
}

func (cache *Cache) MGetKline(prefix string, values *[]Kline) {
	conn := cache.pool.Get()
	defer conn.Close()

	arr, err := redis.Values(conn.Do("scan", 0, "MATCH", prefix, "COUNT", 10000))
	if err != nil {
		LogErr.Error("redis mget kline error:" + err.Error())
		return
	}

	keys, _ := redis.Strings(arr[1], nil)
	var args []interface{}
	for _, k := range keys {
		args = append(args, k)
	}

	v, _ := redis.Strings(conn.Do("MGET", args...))
	j := "["
	last := len(v) - 1
	for i, s := range v {
		if i < last {
			j += s + ","
		} else {
			j += s
		}
	}
	j += "]"

	err = json.Unmarshal([]byte(j), &values)
	if err != nil {
		LogErr.Error("redis unmarshal error:" + err.Error())
		return
	}
}

func (cache *Cache) Puts() {

}
