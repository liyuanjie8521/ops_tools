package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ConsulService struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Port    int      `json:"port"`
	Address string   `json:"address"`
	Checks  []Check  `json:"checks"`
}

type Check struct {
	HTTP     string `json:"http"`
	Interval string `json:"interval"`
}

func main() {
	// 定义consul服务器IP地址和端口
	consulUrl := "http://192.168.1.224:8500/v1/agent/service/register"
	// consulUrl := "http://192.168.1.224:8500/v1/agent/service/deregister/"

	// 定义多台node_exporter的IP地址
	ipAddressList := []string{"192.168.1.180"}

	// 循环遍历IP地址列表进行注册service
	for _, IP := range ipAddressList {
		port := 9121
		url := fmt.Sprintf("http://%s:%d/metrics", IP, port)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error response from %s:%s", IP, resp.Status)
			continue
		}
		// 构造service信息
		service := ConsulService{
			ID:      fmt.Sprintf("redis_exporter-%s:%d", IP, port),
			Name:    "redis_exporter",
			Tags:    []string{"middleware_monitor"},
			Port:    port,
			Address: IP,
			Checks: []Check{
				Check{
					HTTP:     url,
					Interval: "5s",
				},
			},
		}

		// 发送注册请求
		requestBody, err := json.Marshal(service)
		if err != nil {
			fmt.Println("Failed to marshal service JSON: ", err)
			continue
		}
		// 发送注册请求
		request, err := http.NewRequest("PUT", consulUrl, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Failed to create request: ", err)
			continue
		}
		resp, err = http.DefaultClient.Do(request)
		if err != nil {
			fmt.Println("Failed to register service: ", err)
			continue
		}
		defer resp.Body.Close()
		// resp, err = http.NewRequest("PUT", consulUrl, bytes.NewBuffer(requestBody))
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error registering service:%s", resp.Status)
			continue
		}
		fmt.Println("Service registered successfully.")
	}
}
