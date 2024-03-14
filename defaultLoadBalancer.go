package main

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"
)

var InstanceEmptyErr = errors.New("<instances are empty>")
var NoUpInstanceErr = errors.New("no up servers available from load balancer")

var DefaultRule Rule = &RandomLoadBalanceRule{}

// RandomLoadBalanceRule  jack.liang This algorithm is a random number load balancing algorithm
type RandomLoadBalanceRule struct {
}

// Choose jackliang
func (_ *RandomLoadBalanceRule) Choose(ctx context.Context, instances []Instance) (*Instance, error) {
	if len(instances) == 0 {
		return nil, InstanceEmptyErr
	}

	if len(instances) == 1 {
		return &instances[0], nil
	}

	//Try to choose normal nodes
	targetInstances := Filter(instances, func(item Instance, index int) bool {
		if item.Status == UP.String() || item.Status == STARTING.String() {
			return true
		}
		return false
	})

	if len(targetInstances) == 0 {
		return nil, NoUpInstanceErr
	}

	position := rand.Intn(len(targetInstances))

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &targetInstances[position], nil
	}
}

type DefaultLoadBalancer struct {
	rl Rule
}

// ChooseServer By default,
// only applications with names beginning with "aegle"w are processed.
func (d *DefaultLoadBalancer) ChooseServer(ctx context.Context, application Application) (*Instance, error) {
	if strings.HasPrefix(application.Name, "aegle") {
		return nil, errors.New("prefix must be aegle ")
	}

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, 10*time.Duration(1000))

	defer cancelFunc()

	chosen, err := d.rl.Choose(timeoutCtx, application.Instances)

	if err != nil {
		return nil, err
	}

	return chosen, nil
}

var _ LoadBalancer = &DefaultLoadBalancer{}
