package database

import (
	"github.com/fzzy/radix/redis"
	"os"
	"strconv"
	"time"
)

type Redis struct {
	redis.Client
	Host    string
	Port    string
	Timeout time.Duration
}

func (self *Redis) Connect() (err error) {
	self.Host = "109.231.121.62"
	self.Port = "6379"
	self.Timeout = time.Duration(10) * time.Second

	if os.Getenv("DP_REDIS_HOST") != "" {
		self.Host = os.Getenv("DP_REDIS_HOST")
	}

	if os.Getenv("DP_REDIS_PORT") != "" {
		self.Port = os.Getenv("DP_REDIS_PORT")
	}

	if os.Getenv("DP_REDIS_TIMEOUT") != "" {
		timeout, _ := strconv.Atoi(os.Getenv("DP_REDIS_TIMEOUT"))
		self.Timeout = time.Duration(timeout) * time.Second
	}

	Logger.Println("Connecting to Redis " + self.Host + ":" + self.Port)

	var client *redis.Client
	client, err = redis.DialTimeout("tcp", self.Host+":"+self.Port, self.Timeout)

	if err != nil {
		Logger.Println("Could not connect to the redis server.")
		return err
	}

	self.Client = *client

	Logger.Println("Connected!")

	return nil
}
