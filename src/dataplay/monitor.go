package main

import (
	"fmt"
	"net"
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
	Data map[string]interface{}
}

var MonitoringCollection []MonitoringData
var MonitoringLastFlush int
var MonitoringRedisDB int = 1

func StoreMonitoringData(httpEndPoint, httpRequest, httpUrl, httpMethod string, httpCode int, executionTime int64) (e error) {
	timeNow, nanoTime := GetUnixNanoTimeStamp()
	host, _ := GetLocalIp()

	monitoringData := MonitoringData{}
	monitoringData.Hash = nanoTime
	monitoringData.Url = httpEndPoint
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

func FlushMonitoringData(lastFlush int64) error {
	if len(MonitoringCollection) == 0 {
		return nil
	}

	c, err := GetRedisConnection()
	if err != nil {
		return fmt.Errorf("Could not connect to Redis.")
	}

	defer c.Close()

	r := c.Cmd("SELECT", MonitoringRedisDB) // DB 1
	if r.Err != nil {
		return fmt.Errorf("Could not SELECT DB %d from Redis.", MonitoringRedisDB)
	}

	for i, monitor := range MonitoringCollection {
		if monitor.Hash < lastFlush {
			if i > len(MonitoringCollection) {
				break
			}

			hash := strconv.FormatInt(monitor.Hash, 10)

			// Store hash into DB
			r = c.Cmd("HMSET", monitor.Url+":"+hash, monitor.Data)
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

	go func() {
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
	}()
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

	if err != nil {
		return GetHostName()
	}

	for _, address := range addrs {
		fmt.Println("GetLocalIp address", address)
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return GetHostName()
}
