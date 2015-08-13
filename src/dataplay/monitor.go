package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
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

type MonitoringData struct {
	Hash int64
	Url  string
	Host string
	Data map[string]interface{}
}

var MonitoringCollection []MonitoringData
var MonitoringLastFlush int
var MonitoringEndPoint string = "api"
var MonitoringRedisDB int = 1

func StoreMonitoringData(httpEndPoint, httpRequest, httpUrl, httpMethod string, httpCode int, executionTime int64) (e error) {
	timeNow, nanoTime := GetUnixNanoTimeStamp()
	host, _ := GetLocalIp()

	monitoringData := MonitoringData{}
	monitoringData.Hash = nanoTime
	monitoringData.Url = httpEndPoint
	monitoringData.Host = host
	monitoringData.Data = map[string]interface{}{
		"endpoint":  httpEndPoint,
		"request":   httpRequest,
		"method":    httpMethod,
		"url":       httpUrl,
		"code":      httpCode,
		"duration":  executionTime,
		"timestamp": int32(timeNow.Unix()),
		"host":      host,
	}

	MonitoringCollection = append(MonitoringCollection, monitoringData)

	return nil
}

/**
 * @todo Create connection pool with Redis and use same connection for all monitoring related pushes.
 * Since it uses different DB than session will be faster and efficient.
 */
func FlushMonitoringData(lastFlush int64) error {
	if len(MonitoringCollection) == 0 {
		return nil
	}

	c, err := GetRedisConnection()
	if err != nil {
		return fmt.Errorf("Could not connect to Redis.")
	}

	defer c.Close()

	r := c.Cmd("SELECT", MonitoringRedisDB)
	if r.Err != nil {
		return fmt.Errorf("Could not SELECT DB %d from Redis.", MonitoringRedisDB)
	}

	for i, monitor := range MonitoringCollection {
		if monitor.Hash < lastFlush {
			if i > len(MonitoringCollection) {
				break
			}

			/**
			 * api
			 * 	1428676155432861500
			 * 	1428676155722498000
			 * 	1428676155981220800
			 * 	1428677145605146200
			 *
			 * 1.2.3.4
			 * 	1.2.3.4: 1428676155432861500
			 * 	1.2.3.4: 1428676155981220800
			 *
			 * 5.6.7.8
			 * 	5.6.7.8: 1428676155722498000
			 * 	5.6.7.8: 1428677145605146200
			 */
			hash := strconv.FormatInt(monitor.Hash, 10)

			// Store hash into DB
			r = c.Cmd("HMSET", monitor.Host+":"+hash, monitor.Data)
			if r.Err != nil {
				return fmt.Errorf("Could not store HMSET in Redis.")
			}

			// Store pointer to hash's key in DB
			r = c.Cmd("SADD", monitor.Url, hash)
			if r.Err != nil {
				return fmt.Errorf("Could not store SADD in Redis.")
			}

			if i == len(MonitoringCollection) {
				MonitoringCollection = make([]MonitoringData, 0)
			} else {
				MonitoringCollection = MonitoringCollection[:i+copy(MonitoringCollection[i:], MonitoringCollection[i+1:])]
			}
		}
	}

	return nil
}

func AsyncMonitoringPush() {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})

	go MonitoringPush(ticker, quit)
}

func MonitoringPush(ticker *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-ticker.C:
			_, nanoTime := GetUnixNanoTimeStamp()
			FlushMonitoringData(nanoTime)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func GetPerformanceInfo(res http.ResponseWriter, req *http.Request) string {
	endPoint := MonitoringEndPoint
	host, _ := GetLocalIp()

	c, err := GetRedisConnection()
	if err != nil {
		Logger.Println("Could not connect to Redis.", err)
		http.Error(res, "Could not connect to Redis.", http.StatusInternalServerError)
		return ""
	}

	defer c.Close()

	Logger.Println("SELECT DB", MonitoringRedisDB)
	r := c.Cmd("SELECT", MonitoringRedisDB)
	if r.Err != nil {
		Logger.Println("Could not select database from Redis.", r.Err)
		http.Error(res, "Could not select database from Redis.", http.StatusInternalServerError)
		return ""
	}

	// SORT api LIMIT 0 100 GET 5.6.7.8:*->duration BY 5.6.7.8:*->timestamp DESC
	Logger.Println("SORT", endPoint, "LIMIT", 0, 100, "GET", host+":*->duration", "BY", host+":*->timestamp", "DESC")
	sortedData, err := c.Cmd("SORT", endPoint, "LIMIT", 0, 100, "GET", host+":*->duration", "BY", host+":*->timestamp", "DESC").List()
	if err != nil {
		Logger.Println("Could not select keys from Redis.", err)
		http.Error(res, "Could not select keys from Redis.", http.StatusInternalServerError)
		return ""
	}

	data := make([]float64, 0)
	for _, val := range sortedData {
		v, _ := strconv.ParseFloat(val, 10)
		if v > 0 {
			data = append(data, v)
		}
	}

	mean := Mean(data)
	variation := Variation(data)
	standev := StandDev(data)

	info := map[string]interface{}{
		"host":      host,
		"endpoint":  endPoint,
		"mean":      math.Ceil(mean / 1000),
		"standev":   math.Ceil(standev / 1000),
		"variation": variation,
		"timestamp": time.Now(),
	}

	b, _ := json.Marshal(info)

	Logger.Println("Data:", string(b))

	return string(b)
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

/**
 * @brief Get local IP address
 * @details Loop through all network interfaces and get local
 * IP address of the host
 */
func GetLocalIp() (address string, err error) {
	addrs, err := net.InterfaceAddrs()

	if err == nil {
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	return GetHostName()
}
