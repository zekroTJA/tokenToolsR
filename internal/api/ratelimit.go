package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zekroTJA/ratelimit"
)

const (
	limiterLimit = 5 * time.Second
	limiterBurst = 3
)

var limiters = make(map[string]*ratelimit.Limiter)

func CheckRatelimit(w http.ResponseWriter, r *http.Request) bool {
	addr := getIPAddr(r)
	if strings.Contains(addr, ":") {
		split := strings.Split(addr, ":")
		addr = strings.Join(split[0:len(split)-1], ":")
	}

	limiter, ok := limiters[addr]
	if !ok {
		limiter = ratelimit.NewLimiter(limiterLimit, limiterBurst)
		limiters[addr] = limiter
	}

	a, res := limiter.Reserve()

	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", res.Burst))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", res.Reset.UnixNano()))

	if !a {
		SendResponse(w, 427, nil)
	}

	return a
}

func getIPAddr(r *http.Request) string {
	// forwardedfor := ctx.Request.Header.PeekBytes(headerXForwardedFor)
	forwardedfor := r.Header.Get("X-Forwarded-For")
	if forwardedfor != "" {
		return forwardedfor
	}

	return r.RemoteAddr
}
