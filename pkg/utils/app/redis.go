package app

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func RSetup(config AppConfig) error {
	RedisConn = &redis.Pool{
		MaxIdle:     config.REDIS_MAX_IDLE,
		MaxActive:   config.REDIS_MAX_ACTIVE,
		IdleTimeout: config.REDIS_IDLE_TIMEOUT,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.REDIS_HOST)
			if err != nil {
				return nil, err
			}
			if config.REDIS_PASSWORD != "" {
				if _, err := c.Do("AUTH", config.REDIS_PASSWORD); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, _ time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set a key/value
func RSet(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
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

// Exists check a key
func RExists(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Get get a key
func RGet(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a key
func RDelete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func RLikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = RDelete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get value
func RGetValue[T any](key string, defaultValue T) (T, error) {
	var value T

	data, err := RGet(key)

	if err != nil {
		return defaultValue, err
	}

	err = json.Unmarshal(data, &value)

	if err != nil {
		return defaultValue, err
	}

	return value, nil
}
