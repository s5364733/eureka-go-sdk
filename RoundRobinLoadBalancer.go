package main

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var _ LoadBalancer = &RoundRobinLoadBalancer{}

func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{
		rl: &RoundRobinRule{},
	}
}

type RoundRobinLoadBalancer struct {
	rl Rule
}

func (r *RoundRobinLoadBalancer) ChooseServer(ctx context.Context, application Application) (*Instance, error) {

	if len(application.Instances) == 0 {
		return nil, errors.New("<nil instance>")
	}

	if len(application.Instances) == 1 {
		return &application.Instances[0], nil
	}

	targetInstances := Filter(application.Instances, func(item Instance, index int) bool {
		if item.Status == UP.String() || item.Status == STARTING.String() {
			return true
		}
		return false
	})

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, 10*time.Duration(1000))

	defer cancelFunc()

	//Here you need to implement the load balancing algorithm
	chooseOne, err := r.rl.Choose(timeoutCtx, targetInstances)
	if err != nil {
		return nil, err
	}

	return chooseOne, nil
}

var _ Rule = &RoundRobinRule{}

type RoundRobinRule struct {
}

// Choose selects an instance from the list using round-robin.
func (_ *RoundRobinRule) Choose(ctx context.Context, instances []Instance) (*Instance, error) {
	if instances == nil {
		return nil, errors.New("<nil>")
	}

	count := len(instances)

	var index int
	var inst *Instance
	for inst == nil && (index) < 10 {
		position := incrementAndGetModulo(int32(count))
		instance := instances[position]

		if instance.Status == UP.String() || instance.Status == STARTING.String() {
			inst = &instance
		}

		if inst != nil {
			break
		}
		index++
	}

	if index >= 10 {
		println("no available alive servers after 10 tries from load balancer")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return inst, nil
	}
}

var nextServerCyclicCounter int32 = 0

/**
  private int incrementAndGetModulo(int modulo) {
      for (;;) {
          int current = nextServerCyclicCounter.get();
          int next = (current + 1) % modulo;
          if (nextServerCyclicCounter.compareAndSet(current, next))
              return next;
      }
  }
*/
// incrementAndGetModulo increments and gets the next counter value, modulo n
func incrementAndGetModulo(modulo int32) int32 {
	next := (nextServerCyclicCounter) % modulo
	//atomic.AddInt32(&nextServerCyclicCounter, next)
	atomic.AddInt32(&nextServerCyclicCounter, 1)
	return next
}
