package main

import "context"

// LoadBalancer Eureka client loadbalancer
type LoadBalancer interface {
	// ChooseServer Choose Loadbalance to choose  instance
	// when instances size  > 0
	// You can choose appropriate rules based on
	//the characteristics of your application
	ChooseServer(ctx context.Context, application Application) (*Instance, error)
}

type Rule interface {
	// Choose Loadbalance to choose  instance
	// when instances size  > 0
	Choose(ctx context.Context, instances []Instance) (*Instance, error)
}

type noopLoadBalancer struct {
}

func (d noopLoadBalancer) ChooseServer(_ context.Context, _ Application) (*Instance, error) {
	return nil, nil
}

var _ LoadBalancer = noopLoadBalancer{}
