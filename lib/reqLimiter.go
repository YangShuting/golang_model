package lib

import (
	"sync"
	"time"

	"github.com/labstack/echo"
)

type ReqLimiter struct {
	RequestConnect map[string]int
	Lock           sync.Mutex
}

type ReqLimiterService struct {
	Time     time.Duration
	MaxCount int
	ReqLimit ReqLimiter
}

func (reqlimit *ReqLimiterService) GetIPAndUri(c echo.Context) string {
	remoteIP := c.RealIP()
	requestUri := c.Request().RequestURI
	key := remoteIP + ":" + requestUri
	return key
}

func (reqlimit *ReqLimiterService) Increase(key string) {
	if v, exist := reqlimit.ReqLimit.RequestConnect[key]; exist {
		reqlimit.ReqLimit.RequestConnect[key] = v + 1
	} else {
		reqlimit.ReqLimit.RequestConnect[key] = 1
	}
}

func (reqlimit *ReqLimiterService) IsAvaliable(key string) bool {
	if v, exist := reqlimit.ReqLimit.RequestConnect[key]; exist {
		reqlimit.ReqLimit.RequestConnect[key] = v + 1
	}
	return reqlimit.ReqLimit.RequestConnect[key] < reqlimit.MaxCount
}

func NewReqLimiterService(timeD time.Duration, maxCount int) *ReqLimiterService {
	reqLimit := &ReqLimiterService{
		Time:     timeD,
		MaxCount: maxCount,
	}
	reqLimit.ReqLimit.RequestConnect = make(map[string]int)
	go func() {
		ticker := time.NewTicker(timeD)
		for {
			<-ticker.C
			reqLimit.ReqLimit.Lock.Lock()
			for key, _ := range reqLimit.ReqLimit.RequestConnect {
				reqLimit.ReqLimit.RequestConnect[key] = 0
			}
			reqLimit.ReqLimit.Lock.Unlock()
		}
	}()
	return reqLimit
}
