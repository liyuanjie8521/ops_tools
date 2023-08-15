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
	ipAddressList := []string{"192.168.1.224", "192.168.1.106", "192.168.1.123", "192.168.1.125", "192.168.1.139", "192.168.1.140", "192.168.1.141", "192.168.1.142", "192.168.1.143", "192.168.1.144", "192.168.1.146", "192.168.1.147", "192.168.1.148", "192.168.1.149", "192.168.1.67", "192.168.1.139", "192.168.1.44", "192.168.1.162", "192.168.1.163", "192.168.1.164", "192.168.1.165", "192.168.1.166", "192.168.1.167", "192.168.1.168", "192.168.1.169", "192.168.1.181", "192.168.1.190", "192.168.1.42", "192.168.1.43", "192.168.1.44", "192.168.1.26", "192.168.1.10", "192.168.1.47", "192.168.1.46", "192.168.1.228", "192.168.1.51", "192.168.1.48", "192.168.1.18", "192.168.1.42", "192.168.1.129", "192.168.1.127", "192.168.1.131", "192.168.1.68", "192.168.1.82", "192.168.1.103", "192.168.1.94", "192.168.1.95", "192.168.1.126", "192.168.1.182"}

	// 循环遍历IP地址列表进行注册service
	for _, IP := range ipAddressList {
		port := 9100
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
			ID:      fmt.Sprintf("node_exporter-%s:%d", IP, port),
			Name:    "node_exporter",
			Tags:    []string{"node_monitor"},
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
