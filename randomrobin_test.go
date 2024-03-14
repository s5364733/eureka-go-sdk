package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// MockRule 是 Rule 接口的模拟实现
type RandomMockRule struct {
	mock.Mock
}

func (m *RandomMockRule) Choose(ctx context.Context, instances []Instance) (*Instance, error) {
	args := m.Called(ctx, instances)
	return args.Get(0).(*Instance), args.Error(1)
}

func TestRandomLoadBalanceRule_Choose(t *testing.T) {
	// 模拟一组实例
	instances := []Instance{
		{App: "1", Status: UP.String()},
		{App: "2", Status: STARTING.String()},
		{App: "3", Status: DOWN.String()},
	}

	// 使用 MockRule 模拟 Rule 接口
	mockRule := new(RandomMockRule)

	// 测试 Choose 方法
	t.Run("TestChoose", func(t *testing.T) {
		ctx := context.Background()

		// 设置 MockRule 期望的行为
		mockRule.On("Choose", ctx, mock.Anything).Return(&Instance{App: "1", Status: UP.String()}, nil)

		// 测试模拟请求是否生效
		chosenInstance, err := mockRule.Choose(ctx, instances)
		require.NoError(t, err)
		assert.NotNil(t, chosenInstance)
		assert.Equal(t, "1", chosenInstance.App)

	})
}

func TestRandomLoadBalanceRule_ChooseRefreshSeed(t *testing.T) {
	// 模拟一组实例
	instances := []Instance{
		{App: "1", Status: UP.String()},
		{App: "2", Status: STARTING.String()},
		{App: "3", Status: DOWN.String()},
	}

	rule := &RandomLoadBalanceRule{}

	// 测试 Choose 方法
	t.Run("TestChoose", func(t *testing.T) {
		ctx := context.Background()

		// 测试模拟请求是否生效
		chosenInstance, err := rule.Choose(ctx, instances)
		require.NoError(t, err)
		assert.NotNil(t, chosenInstance)
		assert.Equal(t, "2", chosenInstance.App)

	})
}

func TestRandomLoadBalanceRule_ChooseDownFilter(t *testing.T) {

	// 模拟一组实例
	instances := []Instance{
		{App: "1", Status: DOWN.String()},
		{App: "2", Status: DOWN.String()},
		{App: "3", Status: DOWN.String()},
	}

	rule := &RandomLoadBalanceRule{}

	// 测试 Choose 方法
	t.Run("TestChoose", func(t *testing.T) {
		ctx := context.Background()

		// 理论上如果三个节点都挂了应该返回
		chosenInstance, err := rule.Choose(ctx, instances)
		assert.Error(t, err)
		assert.Nil(t, chosenInstance)
		assert.ErrorIs(t, err, NoUpInstanceErr)
	})
}

/*
*
GOROOT=D:\fz\env\go_env\go #gosetup
GOPATH=D:\fz\work #gosetup
D:\fz\env\go_env\go\bin\go.exe test -c -o C:\Users\13543\AppData\Local\JetBrains\GoLand2023.2\tmp\GoLand\___TestRandomLoadBalanceRule_TestChoose_in_github_com_xuanbo_eureka_client.test.exe github.com/xuanbo/eureka-client #gosetup
D:\fz\env\go_env\go\bin\go.exe tool test2json -t C:\Users\13543\AppData\Local\JetBrains\GoLand2023.2\tmp\GoLand\___TestRandomLoadBalanceRule_TestChoose_in_github_com_xuanbo_eureka_client.test.exe -test.v -test.paniconexit0 -test.run ^\QTestRandomLoadBalanceRule\E$/^\QTestChoose\E$ #gosetup
2024/02/07 12:17:36 Heartbeat application instance successfully
=== RUN   TestRandomLoadBalanceRule
=== RUN   TestRandomLoadBalanceRule/TestChoose
&{    3    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    1    UP <nil> <nil> <nil> <nil> map[]      0 }
&{    3    UP <nil> <nil> <nil> <nil> map[]      0 }
--- PASS: TestRandomLoadBalanceRule (0.00s)

	--- PASS: TestRandomLoadBalanceRule/TestChoose (0.00s)

PASS
*/
func TestRandomLoadBalanceRule(t *testing.T) {

	// 模拟一组实例
	instances := []Instance{
		{App: "1", Status: UP.String()},
		{App: "2", Status: UP.String()},
		{App: "3", Status: UP.String()},
	}

	rule := &RandomLoadBalanceRule{}

	// 测试 Choose 方法
	t.Run("TestChoose", func(t *testing.T) {
		ctx := context.Background()
		for i := 0; i < 3; i++ {
			// 基本轮询
			chosenInstance, _ := rule.Choose(ctx, instances)
			fmt.Println(chosenInstance)
		}

	})
}
