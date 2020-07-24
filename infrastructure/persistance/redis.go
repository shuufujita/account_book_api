package persistance

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var cachePool *redis.Pool

// GetRedisPool return redis connection pool
func GetRedisPool() (*redis.Pool, error) {
	maxConnSize, err := strconv.Atoi(os.Getenv("REDIS_MAX_CONN"))
	if err != nil {
		log.Println("error: [redis] convert int " + err.Error())
	}

	maxPoolSize, err := strconv.Atoi(os.Getenv("REDIS_MAX_CONN_POOL"))
	if err != nil {
		log.Println("error: [redis] convert int " + err.Error())
	}

	cachePool = &redis.Pool{
		Wait:        true,
		IdleTimeout: 30 * time.Second,
		MaxActive:   maxConnSize,
		MaxIdle:     maxPoolSize,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp",
				os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	return cachePool, nil
}

// GetConn return redis pool connection
func GetConn() redis.Conn {
	return cachePool.Get()
}

// RedisGet get key/value
func RedisGet(key string) (string, error) {
	conn := GetConn()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return "", err
	}

	return data, nil
}

// RedisSetJSON redis set with convert to json
func RedisSetJSON(key string, data interface{}, ttl int) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	conn := GetConn()
	defer conn.Close()

	_, err = conn.Do("SET", key, string(b))
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, ttl)
	if err != nil {
		return err
	}

	return nil
}

// RedisSet save key/value
func RedisSet(key string, data string, ttl int) error {
	conn := GetConn()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, ttl)
	if err != nil {
		return err
	}

	return err
}
