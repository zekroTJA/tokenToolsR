package ratelimit

import (
	"strings"
	"time"
)

var RATE_LIMITS = map[string]int64{
	"general": 1,
	"guilds":  2,
}

var requestTimes = map[[2]string]int64{}

func CheckRatelimit(addr, endpoint string) bool {
	limit := RATE_LIMITS[endpoint]
	addrSplit := strings.Split(addr, ":")[0]
	if addr == "" || addrSplit == "[" {
		return false
	}
	if reqTime, ok := requestTimes[[2]string{addr, endpoint}]; ok {
		if reqTime+limit > time.Now().Unix() {
			return true
		}
	}
	requestTimes[[2]string{addr, endpoint}] = time.Now().Unix()
	return false
}
