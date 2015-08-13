package main

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"os"
	"strconv"
)

func GetRedisConnection() (client *redis.Client, err error) {
	redisHost := "109.231.121.62"
	redisPort := "6379"
	redisPoolSize := 100

	if os.Getenv("DP_REDIS_HOST") != "" {
		redisHost = os.Getenv("DP_REDIS_HOST")
	}

	if os.Getenv("DP_REDIS_PORT") != "" {
		redisPort = os.Getenv("DP_REDIS_PORT")
	}

	if os.Getenv("DP_REDIS_POOLSIZE") != "" {
		redisPoolSize, _ = strconv.Atoi(os.Getenv("DP_REDIS_POOLSIZE"))
	}

	Logger.Println("Connecting to Redis " + redisHost + ":" + redisPort)

	p, err := pool.New("tcp", redisHost+":"+redisPort, redisPoolSize)
	if err != nil {
		Logger.Println("Could not initiate Redis Pool.")
		return nil, err
	}

	client, err = p.Get()
	if err != nil {
		Logger.Println("Could not connect to the Redis server.")
		return nil, err
	}

	defer p.Put(client)

	Logger.Println("Connected!")

	return client, nil
}
