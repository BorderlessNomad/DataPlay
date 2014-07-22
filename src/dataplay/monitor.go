package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

/**
 * @brief Store monitoring data into Redis
 * @details Following method stores requested data and associated parameters
 * into redis. Utlizing 'Listed Hash'.
 *
 * @param httpEndPoint [string] e.g. "api"
 * @param httpRequest [string] e.g. "search"
 * @param httpUrl [string], full request e.g. "/api/search/nhs"
 * @param httpMethod [string] e.g. GET, POST etc
 * @param httpCode [int] e.g. 200, 500
 * @param executionTime [int], total execution time to complete request
 */
func StoreMonitoringData(httpEndPoint, httpRequest, httpUrl, httpMethod string, httpCode int, executionTime int64) (e error) {
	timeNow, nanoTime := GetUnixNanoTimeStamp()

	c, err := GetRedisConnection()
	if err != nil {
		return fmt.Errorf("Could not connect to Redis.")
	}

	defer c.Close()

	r := c.Cmd("SELECT", 1) // DB 1
	if r.Err != nil {
		return fmt.Errorf("Could not select database from Redis.")
	}

	host, err := GetHostName()

	data := map[string]interface{}{
		"endpoint":  httpEndPoint,
		"request":   httpRequest,
		"method":    httpMethod,
		"url":       httpUrl,
		"code":      httpCode,
		"duration":  executionTime,
		"timestamp": int32(timeNow.Unix()),
		"host":      host,
	}

	key := strconv.FormatInt(nanoTime, 10)

	// Store hash into DB
	r = c.Cmd("HMSET", httpEndPoint+":"+key, data)
	if r.Err != nil {
		return fmt.Errorf("Could not store HMSET in Redis.")
	}

	// Store pointer to hash's key in DB
	r = c.Cmd("SADD", httpEndPoint, key)
	if r.Err != nil {
		return fmt.Errorf("Could not store SADD in Redis.")
	}

	return nil
}

func GetUnixNanoTimeStamp() (time.Time, int64) {
	now := time.Now()
	nanos := now.UnixNano()

	return now, nanos
}

/**
 * @brief Get Hostname
 * @details Retrieve hostname of a system along with NS-Lookup
 */
func GetHostName() (name string, err error) {
	name, err = os.Hostname()

	return
}
