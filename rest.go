package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/internal/backoff"
	"net/url"
	"time"
)

func (c *Client) restRpcWithJson(ctx context.Context, urlx string, method string) (string, error) {
	fmt.Printf("rest rpc with json luanching url %s method %s \t\n", urlx, method)

	raw, err := url.Parse(urlx)

	if err != nil {
		return "", fmt.Errorf("rest rpc parsing error %w", err)
	}

	var instance *Instance
	if raw.Hostname() != "" {
		app := c.GetApplicationInstance(raw.Host)
		if app == nil {
			return "", errors.New("<nil>")
		}

		server, err1 := c.Lb.ChooseServer(ctx, *app)
		if err1 != nil {
			return "", fmt.Errorf("rest rpc error %w", err1)
		}
		instance = server
	}

	if instance == nil {
		return "", NoUpInstanceErr
	}

	var res any

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)

	defer cancelFunc()

	for {
		resJson, err := c.Exc.Exchange(timeoutCtx, fmt.Sprintf("%s://%s:%d%s",
			raw.Scheme,
			instance.IPAddr,
			instance.Port.Port,
			raw.Path),
			method)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return "", err
			}
			d := backoff.BackCfg.Duration()
			_ = fmt.Errorf("%s, reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}

		res = resJson

		return res.(string), nil
	}

}
