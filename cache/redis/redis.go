package redis

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var Pool *redis.Pool

// Setup Initialize the Redis instance
func Setup(host,maxIdle,maxActive,idleTimeout,passWord interface{}) error {
	conn, err := redis.Dial("tcp", host.(string))
	if err != nil {
		log.Panic("Connect to redis error", err)
		return err
	}
	_ = conn.Close()
	Pool = &redis.Pool{
		MaxIdle:     maxIdle.(int),
		MaxActive:   maxActive.(int),
		IdleTimeout: idleTimeout.(time.Duration),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host.(string))
			if err != nil {
				return nil, err
			}
			if passWord.(string) != "" {
				if _, err := c.Do("AUTH", passWord.(string)); err != nil {
					err := c.Close()
					if err != nil {
						return nil, err
					}
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			conn, err := c.Do("PING")
			if conn == nil {
				return errors.New("gan")
			}
			return err
		},
	}
	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := Pool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// SetString a key/value
func SetString(key, data string, time int) error {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := Pool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := Pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
