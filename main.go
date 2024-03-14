package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := SpawnClient(func(cfg *Config) {
		cfg.DefaultZone = "http://localhost:8761/eureka/"
		cfg.InstanceID = "aegle-eureka-example"
		cfg.App = "aegle-eureka-example"
		cfg.Port = 10000
		cfg.RenewalIntervalInSecs = 10
		cfg.RegistryFetchIntervalSeconds = 15
		cfg.DurationInSecs = 30
		cfg.Metadata = map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		}
	})
	go func() {
		// http server
		http.HandleFunc("/v1/services", func(writer http.ResponseWriter, request *http.Request) {
			// full applications from eureka server
			apps := "{\"instance\":{\"hostName\":\"172.30.224.1\",\"homePageUrl\":\"http://172.30.224.1:10000\",\"statusPageUrl\":\"http://172.30.224.1:10000/info\",\"app\":\"aegle-example\",\"ipAddr\":\"172.30.224.1\",\"vipAddress\":\"aegle-example\",\"secureVipAddress\":\"aegle-example\",\"status\":\"UP\",\"port\":{\"$\":10000,\"@enabled\":\"true\"},\"dataCenterInfo\":{\"name\":\"MyOwn\",\"@class\":\"com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo\"},\"leaseInfo\":{\"renewalIntervalInSecs\":30,\"durationInSecs\":90},\"metadata\":{\"NODE_GROUP_ID\":0,\"PRODUCT_CODE\":\"DEFAULT\",\"PRODUCT_ENV_CODE\":\"DEFAULT\",\"PRODUCT_VERSION_CODE\":\"DEFAULT\",\"SERVICE_VERSION_CODE\":\"DEFAULT\",\"VERSION\":\"0.1.0\"},\"overriddenstatus\":\"UNKNOWN\",\"instanceId\":\"aegle-example\"}}"

			b, _ := json.Marshal(apps)
			_, _ = writer.Write(b)
		})

		// start http server
		if err := http.ListenAndServe(":10000", nil); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(3 * time.Second)

	// http server
	http.HandleFunc("/rpc/service", func(writer http.ResponseWriter, request *http.Request) {
		// full applications from eureka server
		withJson, err := client.restRpcWithJson(context.TODO(), "http://aegle-eureka-example/v1/services", "POST")

		if err != nil {
			fmt.Println(err)
		}

		_, _ = writer.Write([]byte(withJson))
	})

	// start http server
	if err := http.ListenAndServe(":10001", nil); err != nil {
		fmt.Println(err)
	}

}
