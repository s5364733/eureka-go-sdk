package main

import (
	"context"
	"errors"
	"github.com/requests"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var _ Exchange = &DefaultExchange{}

type DefaultExchange struct {
	httpClient *http.Client
}

func NewDefaultExchange() *DefaultExchange {
	return &DefaultExchange{
		httpClient: &http.Client{
			Timeout: 6 * time.Second,
		},
	}

}
func (httpExchange *DefaultExchange) Exchange(ctx context.Context, url string, method string) (interface{}, error) {
	text, err := requests.Request(url, method, httpExchange.httpClient).Send().Text()
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return text, nil
	}
}

type Exchange interface {
	Exchange(ctx context.Context, url string, method string) (interface{}, error)
}

type EurekaClient interface {
	Start()
	refresh()
	heartbeat()
	doRegister() error
	doUnRegister() error
	handleSignal()
	doHeartbeat() error
	doRefresh() error
	GetApplicationInstance(name string) *Application
}

// Client eureka客户端
type Client struct {
	logger Logger

	// for monitor system signal
	signalChan chan os.Signal
	mutex      sync.RWMutex
	running    bool

	Config   *Config
	Instance *Instance

	// eureka服务中注册的应用
	Applications *Applications

	Lb LoadBalancer

	Exc Exchange
}

// SetLogger 设置日志实现
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// Start 启动时注册客户端，并后台刷新服务列表，以及心跳
func (c *Client) Start() {
	c.mutex.Lock()
	c.running = true
	c.mutex.Unlock()
	// 心跳
	go c.heartbeat()
	time.Sleep(1 * time.Second)
	// 刷新本地服务列表
	go c.refresh()
	//// 监听退出信号，自动删除注册信息
	go c.handleSignal()
}

// refresh 刷新服务列表
func (c *Client) refresh() {
	timer := time.NewTimer(0)
	interval := time.Duration(c.Config.RenewalIntervalInSecs) * time.Second
	for c.running {
		select {
		case <-timer.C:
			if err := c.doRefresh(); err != nil {
				c.logger.Error("Refresh application instance failed", err)
			} else {
				c.logger.Debug("Refresh application instance successfully")
			}
		}
		// reset interval
		timer.Reset(interval)
	}
	// stop
	timer.Stop()
}

// heartbeat 心跳
func (c *Client) heartbeat() {
	timer := time.NewTimer(0)
	interval := time.Duration(c.Config.RegistryFetchIntervalSeconds) * time.Second
	for c.running {
		select {
		case <-timer.C:
			err := c.doHeartbeat()
			if err == nil {
				c.logger.Debug("Heartbeat application instance successfully")
			} else if errors.Is(err, ErrNotFound) {
				// heartbeat not found, need register
				err = c.doRegister()
				if err == nil {
					c.logger.Info("Register application instance successfully")
				} else {
					c.logger.Error("Register application instance failed", err)
				}
			} else {
				c.logger.Error("Heartbeat application instance failed", err)
			}
		}
		// reset interval
		timer.Reset(interval)
	}
	// stop
	timer.Stop()
}

func (c *Client) doRegister() error {
	return Register(c.Config.DefaultZone, c.Config.App, c.Instance)
}

func (c *Client) doUnRegister() error {
	return UnRegister(c.Config.DefaultZone, c.Instance.App, c.Instance.InstanceID)
}

func (c *Client) doHeartbeat() error {
	return Heartbeat(c.Config.DefaultZone, c.Instance.App, c.Instance.InstanceID)
}

func (c *Client) doRefresh() error {
	// get all applications
	applications, err := Refresh(c.Config.DefaultZone)
	if err != nil {
		return err
	}

	if applications != nil && len(applications.Applications) > 0 {
		// set applications
		c.mutex.Lock()
		c.Applications = applications
		c.mutex.Unlock()
	}
	return nil
}

// handleSignal 监听退出信号，删除注册的实例
func (c *Client) handleSignal() {
	if c.signalChan == nil {
		c.signalChan = make(chan os.Signal)
	}
	signal.Notify(c.signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		switch <-c.signalChan {
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGKILL:
			fallthrough
		case syscall.SIGTERM:
			c.logger.Info("Receive exit signal, client instance going to de-register")
			err := c.doUnRegister()
			if err != nil {
				c.logger.Error("De-register application instance failed", err)
			} else {
				c.logger.Info("De-register application instance successfully")
			}
			os.Exit(0)
		}
	}
}

// SpawnClient 创建客户端
func SpawnClient(opts ...Option) *Client {
	var cfg Config
	var err error
	if opts == nil {
		cfg, err = ParseConfig("config.toml")
		if err != nil {
			panic(err)
		}
	}

	cfg.extraConfig()

	for _, opt := range opts {
		opt.ExtraCfg(&cfg)
	}

	loadBalancer := NewRoundRobinLoadBalancer()

	instance := NewInstance(&cfg)

	client := &Client{
		logger:   NewLogger(),
		Config:   &cfg,
		Instance: instance,
		Lb:       loadBalancer,
		Exc:      NewDefaultExchange(),
	}

	client.Start()

	return client
}

func (cfg *Config) extraConfig() {
	if cfg.IP == "" {
		cfg.IP = GetLocalIP()
	}

	if cfg.HostName == "" {
		cfg.HostName = GetLocalIP()
	}
}

// GetApplicationInstance 根据服务名获取注册的服务实例列表
func (c *Client) GetApplicationInstance(name string) *Application {
	instances := make([]Instance, 0)
	c.mutex.Lock()
	if c.Applications != nil {
		for _, app := range c.Applications.Applications {
			if strings.EqualFold(app.Name, name) {
				instances = append(instances, app.Instances...)
			}
		}
	}
	c.mutex.Unlock()

	if len(instances) == 0 {
		return nil
	}

	return &Application{
		Instances: instances,
	}
}
