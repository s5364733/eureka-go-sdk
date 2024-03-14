package backoff

import (
	"github.com/jpillora/backoff"
	"time"
)

var BackCfg = &backoff.Backoff{
	//These are the defaults
	Min:    100 * time.Millisecond,
	Max:    30 * time.Second,
	Factor: 2,
	Jitter: false,
}

//
//
//
//var ExceededRetryErr = errors.New("retrying to reach limit, failing fast")
//
//var r = rand.New(rand.NewSource(time.Now().UnixNano()))
//
//// Strategy Retry strategy
//type Strategy interface {
//	Backoff(retries int) time.Duration
//}
//
//// default DefaultExponential
//
//var DefaultExponential = Exponential{Config: DefaultConfig}
//
//type Exponential struct {
//	Config Config
//}
//
//func (bc Exponential) Backoff(retries int) time.Duration {
//	if retries == 0 {
//		return bc.Config.BaseDelay
//	}
//	backoff, max := float64(bc.Config.BaseDelay), float64(bc.Config.MaxDelay)
//	for backoff < max && retries > 0 {
//		backoff *= bc.Config.Multiplier
//		retries--
//	}
//	if backoff > max {
//		backoff = max
//	}
//	//随机化退避延迟，这样如果一组请求同时启动，它们就不会同步运行。
//	// Randomize backoff delays so that if a cluster of requests start at
//	// the same time, they won't operate in lockstep.
//	backoff *= 1 + bc.Config.Jitter*(r.Float64()*2-1)
//	if backoff < 0 {
//		return 0
//	}
//	return time.Duration(backoff)
//}
