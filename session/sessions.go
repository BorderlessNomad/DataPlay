package session

import (
	"crypto/rand"
	"fmt"
	"github.com/fzzy/radix/redis"
	"net/http"
	"time"
)

// session.res.Header().Set("Set-Cookie", fmt.Sprintf("SessionId=; path=%s;", session.manager.path))
func IsUserLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	cookie, _ := req.Cookie("DPSession")
	c, err := GetRedisConnection()
	defer c.Close()
	if cookie != nil {
		if err != nil {
			return false
		}
		defer c.Close()
		r := c.Cmd("GET", cookie.Value)
		i, err := r.Int()
		if err != nil {
			return false
		}
		if i != 0 {
			return true
		} else {
			return false // there is no zero user.
		}
	} else {
		return false
	}
}

func GetUserID(res http.ResponseWriter, req *http.Request) int {
	cookie, _ := req.Cookie("DPSession")
	c, err := GetRedisConnection()
	defer c.Close()
	if cookie != nil {
		if err != nil {
			return 0
		}
		defer c.Close()
		r := c.Cmd("GET", cookie.Value)
		i, err := r.Int()
		if err != nil {
			return 0
		}
		if i != 0 {
			return i
		} else {
			return 0 // there is no zero user.
		}
	} else {
		return 0
	}
}

func SetSession(res http.ResponseWriter, req *http.Request, userid int) (e error) {
	NewSessionID := randString(64)
	c, err := GetRedisConnection()
	defer c.Close()
	if err != nil {
		return fmt.Errorf("Could not connect to redis server to make session")
	}
	r := c.Cmd("SET", NewSessionID, userid)
	if r.Err != nil {
		return fmt.Errorf("Could not store session in Redis D:")
	}
	res.Header().Set("Set-Cookie", fmt.Sprintf("DPSession=%s; path=/; expires=Thu, 01-Jan-2030 00:00:00 GMT;", NewSessionID))
	return e
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func GetRedisConnection() (c *redis.Client, err error) {
	c, err = redis.DialTimeout("tcp", "10.0.0.2:6379", time.Duration(10)*time.Second)
	return c, err
}
