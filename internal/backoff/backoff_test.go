package backoff

//
//import (
//	"github.com/stretchr/testify/assert"
//	"math/rand"
//	"testing"
//	"time"
//)
//
//func TestExponentialBackoff(t *testing.T) {
//	// Seed the random number generator
//	rand.Seed(time.Now().UnixNano())
//
//	// Seed the random number generator
//	rand.Seed(time.Now().UnixNano())
//
//	t.Run("BackoffInitialRetry", func(t *testing.T) {
//		strategy := Exponential{Config: DefaultConfig}
//		backoff := strategy.Backoff(0)
//		assert.Equal(t, DefaultConfig.BaseDelay, backoff, "Unexpected backoff duration for initial retry")
//	})
//
//	t.Run("BackoffWithMaxDelay", func(t *testing.T) {
//		maxDelay := time.Second * 10
//		config := Config{BaseDelay: time.Millisecond, MaxDelay: maxDelay, Multiplier: 2, Jitter: 0.1}
//		strategy := Exponential{Config: config}
//		retries := 5
//		backoff := strategy.Backoff(retries)
//		assert.LessOrEqual(t, backoff, maxDelay, "Backoff duration exceeds maximum allowed delay")
//	})
//
//	t.Run("RandomizedBackoff", func(t *testing.T) {
//		strategy := Exponential{Config: DefaultConfig}
//		retries := 2
//		backoff1 := strategy.Backoff(retries)
//		backoff2 := strategy.Backoff(retries)
//		assert.NotEqual(t, backoff1, backoff2, "Randomized backoff durations should not be equal")
//	})
//}
