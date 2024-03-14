package main

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestRoundRobinLoadBalancer_ChooseServer(t *testing.T) {
	// 模拟一组实例
	instances := []Instance{
		{App: "1", Status: UP.String()},
		{App: "2", Status: STARTING.String()},
		{App: "3", Status: DOWN.String()},
	}

	application := Application{
		Instances: instances,
	}

	// 创建一个 RoundRobinLoadBalancer
	lb := NewRoundRobinLoadBalancer()

	// 测试 ChooseServer 方法
	t.Run("TestChooseServer", func(t *testing.T) {
		ctx := context.Background()

		// 测试正常情况，应该选择 UP 或 STARTING 状态的实例
		chosenInstance, err := lb.ChooseServer(ctx, application)
		require.NoError(t, err)
		assert.NotNil(t, chosenInstance)
		assert.Contains(t, []string{UP.String(), STARTING.String()}, chosenInstance.Status)

		// 测试没有可用实例的情况
		emptyApplication := Application{}
		_, err = lb.ChooseServer(ctx, emptyApplication)
		assert.Error(t, err)
		assert.Equal(t, errors.New("<nil instance>"), err)

		// 测试在规定时间内未能选择可用实例的情况
		instancesWithoutAvailable := []Instance{
			{App: "4", Status: DOWN.String()},
			{App: "5", Status: DOWN.String()},
		}
		applicationWithoutAvailable := Application{
			Instances: instancesWithoutAvailable,
		}
		_, err = lb.ChooseServer(ctx, applicationWithoutAvailable)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no up servers available from load balancer")
	})
}

// MockRule 是 Rule 接口的模拟实现
type MockRule struct {
	mock.Mock
}

func (m *MockRule) Choose(ctx context.Context, instances []Instance) (*Instance, error) {
	args := m.Called(ctx, instances)
	return args.Get(0).(*Instance), args.Error(1)
}

func TestRoundRobinLoadBalancer_MockChooseServer(t *testing.T) {
	// 模拟一组实例
	instances := []Instance{
		{App: "1", Status: UP.String()},
		{App: "2", Status: STARTING.String()},
		{App: "3", Status: DOWN.String()},
	}

	application := Application{
		Instances: instances,
	}

	// 创建一个 RoundRobinLoadBalancer
	lb := NewRoundRobinLoadBalancer()

	// 使用 MockRule 模拟 Rule 接口
	mockRule := new(MockRule)
	lb.rl = mockRule

	// 测试 ChooseServer 方法
	t.Run("TestChooseServer", func(t *testing.T) {
		ctx := context.Background()

		// 设置 MockRule 期望的行为
		mockRule.On("Choose", ctx, mock.Anything).Return(&Instance{App: "1", Status: UP.String()}, nil)

		// 测试模拟请求是否生效
		chosenInstance, err := lb.ChooseServer(ctx, application)
		require.NoError(t, err)
		assert.NotNil(t, chosenInstance)
		assert.Equal(t, "1", chosenInstance.App)

	})
}

/*
*

	2024/02/07 12:38:59 Heartbeat application instance successfully

=== RUN   TestRoundRobinLoadBalancer_Round

=== RUN   TestRoundRobinLoadBalancer_Round/TestChooseServer
&{    1    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    2    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    3    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    1    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    2    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    3    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    1    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    2    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    3    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    1    UP <nil> <nil> <nil> <nil> map[]      0 }
--- PASS: TestRoundRobinLoadBalancer_Round (0.00s)

	--- PASS: TestRoundRobinLoadBalancer_Round/TestChooseServer (0.00s)

PASS
*/
func TestRoundRobinLoadBalancer_Round(t *testing.T) {
	// 模拟一组实例
	instances := []Instance{
		{App: "1", Status: UP.String()},
		{App: "2", Status: UP.String()},
		{App: "3", Status: UP.String()},
	}

	application := Application{
		Instances: instances,
	}

	// 创建一个 RoundRobinLoadBalancer
	lb := NewRoundRobinLoadBalancer()

	// 测试 ChooseServer 方法
	index := 1
	t.Run("TestChooseServer", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			server, _ := lb.ChooseServer(context.Background(), application)
			appIndex, _ := strconv.Atoi(server.App)
			assert.Equal(t, appIndex, index)
			index++
		}

	})
}
