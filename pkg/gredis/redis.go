package gredis

import (
	"encoding/json"
	"gin-blog/pkg/setting"
	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle : setting.RedisSetting.MaxIdle,
		MaxActive : setting.RedisSetting.MaxActive,
		IdleTimeout : setting.RedisSetting.IdleTimeout,

		Dial : func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil , err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err !=nil{
					c.Close()
					return nil ,err
				}
			}
			return c, err
		},
	}
	return nil
}

func Set(key string ,data interface{},time int) (bool, error){
	conn := RedisConn.Get()
	defer conn.Close()

	value ,err := json.Marshal(data)
	if err != nil {
		return false , err
	}
	reply , err := redis.Bool(conn.Do("SET",key, value))
	conn.Do("EXPIRE", key, time )
	return reply, err
}