package cache

import (
	"inochat/server/config"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Instance().Redis.Host)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}
			if config.Instance().Redis.Passwd != "" {
				if _, err := c.Do("AUTH", config.Instance().Redis.Passwd); err != nil {
					c.Close()
					return nil, errors.Wrap(err, "")
				}
			}
			return c, errors.Wrap(err, "")
		},
	}
}

func GetConn() redis.Conn {
	conn := pool.Get()
	return conn
}

func Get(key string) (string, error) {
	conn := GetConn()
	defer conn.Close()
	cache, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return cache, nil
}

func Del(key string) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func Set(key, val string, expire ...int) error {
	conn := GetConn()
	defer conn.Close()

	err := conn.Send("SET", key, val)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if len(expire) > 0 {
		err = conn.Send("EXPIRE", key, expire[0])
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	err = conn.Flush()
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func Lindex(key string, index int) (string, error) {
	conn := GetConn()
	defer conn.Close()
	cache, err := redis.String(conn.Do("LINDEX", key, index))
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return cache, nil
}

func Len(key string) int {
	conn := GetConn()
	defer conn.Close()
	cache, err := redis.Int(conn.Do("LLEN", key))
	if err != nil {
		return 0
	}
	return cache
}

func LPush(key string, val ...string) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("LPUSH", key, val)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func RPop(key string) (string, error) {
	conn := GetConn()
	defer conn.Close()
	cache, err := redis.String(conn.Do("RPOP", key))
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return cache, nil
}

// 再加一个HASH函数
