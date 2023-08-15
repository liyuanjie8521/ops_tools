package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "os/exec"
	"strconv"
	// "strings"
	"bytes"
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

	// 生成254个IP地址
	numIPs := 254
	IPRange := "192.168.1."
	IPs := []string{}

	for i := 1; i <= numIPs; i++ {
		IP := IPRange + strconv.Itoa(i)
		IPs = append(IPs, IP)
	}
	// 将每个IP地址注册到Consul接口
	for _, IP := range IPs {
		port := 6379
		url := fmt.Sprintf("http://192.168.1.180:9121/scrape?target=%s:%d", IP, port)
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
			ID:      fmt.Sprintf("multi_redis_instance-%s:%d", IP, port),
			Name:    "multi_redis_instance",
			Tags:    []string{"redis", "test", "middleware_monitor"},
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
