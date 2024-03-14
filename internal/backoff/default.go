package backoff

//
//import (
//	"context"
//	"time"
//)
//
//// Config defines the configuration options for backoff.
//type Config struct {
//	BaseDelay time.Duration
//	//乘数是重试失败后乘以退避的因子。 理想情况下应大于 1。
//	Multiplier float64
//	//抖动是使退避随机化的因素。
//	Jitter float64
//	//MaxDelay 是退避延迟的上限。
//	MaxDelay time.Duration
//}
//
//var DefaultConfig = Config{
//	BaseDelay:  1.0 * time.Second,
//	Multiplier: 1.2, //默认按1.2倍放大延迟，最大10S之后结束重试，
//	// 您可以改写 BackoffStrategy.Backoff func 来满足您的需求
//	Jitter:   0.2,
//	MaxDelay: 20 * time.Second,
//}
//
//var (
//	BackoffStrategy = DefaultExponential
//	BackoffFunc     = func(ctx context.Context, retries int) bool {
//		d := BackoffStrategy.Backoff(retries)
//		timer := time.NewTimer(d)
//		select {
//		case <-timer.C:
//			return true
//		case <-ctx.Done():
//			timer.Stop()
//			return false
//		}
//	}
//)
