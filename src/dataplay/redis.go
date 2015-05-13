package main

import (
	"github.com/fzzy/radix/redis"
	"os"
	"strconv"
	"time"
)

func GetRedisConnection() (client *redis.Client, err error) {
	redisHost := "109.231.124.16"
	redisPort := "6379"
	redisTimeout := time.Duration(10) * time.Second

	if os.Getenv("DP_REDIS_HOST") != "" {
		redisHost = os.Getenv("DP_REDIS_HOST")
	}

	if os.Getenv("DP_REDIS_PORT") != "" {
		redisPort = os.Getenv("DP_REDIS_PORT")
	}

	if os.Getenv("DP_REDIS_TIMEOUT") != "" {
		timeout, _ := strconv.Atoi(os.Getenv("DP_REDIS_TIMEOUT"))
		redisTimeout = time.Duration(timeout) * time.Second
	}

	Logger.Println("Connecting to Redis " + redisHost + ":" + redisPort + " with " + strconv.FormatFloat(redisTimeout.Seconds(), 'f', -1, 64) + " secs timeout...")
	client, err = redis.DialTimeout("tcp", redisHost+":"+redisPort, redisTimeout)

	if err != nil {
		Logger.Println("Could not connect to the redis server.")
		return nil, err
	}

	Logger.Println("Connected!")

	return client, nil
}
