package storage

import (
	"github.com/Yangshuting/golang_model/config"
	"fmt"
	"strconv"
	"time"
	"github.com/go-redis/redis"
)

func RedisConn() *redis.Client {
	fmt.Printf("address_%+v_%+v \n", config.RawEnv("REDIS_HOST"), config.RawEnv("REDIS_PORT"))
	var redisClient *redis.Client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RawEnv("REDIS_HOST") + ":" + config.RawEnv("REDIS_PORT"),
		Password: config.RawEnv("REDIS_PASSWORD"),
		DB:       0,
	})
	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
	return redisClient
}

func GetRedis(key string) (string, error) {
	val, err := RedisConn().Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf(key, " does not exists.")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func SetRedis(key, val string) error {
	var expireTime int
	expireTime = StringToInt(config.RawEnv("REDIS_DEFAULT_EXPIRED_SECOND"))
	err := RedisConn().Set(key, val, time.Duration(expireTime)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}
func StringToInt(str string) int {
	intK, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return int(intK)
}
