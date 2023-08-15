package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	// "os/exec"
	// "strings"
	"github.com/cheggaaa/pb/v3"
	"time"
)

type RedisInfo struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

func main() {
	// 扫描192.168.1网段的所有6379端口
	var redisList []RedisInfo
	bar := pb.StartNew(100)
	for i := 1; i <= 100; i++ {
		ip := fmt.Sprintf("192.168.1.%d", i)
		addr := fmt.Sprintf("%s:%d", ip, 6379)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			bar.Increment()
			continue
		}
		defer conn.Close()

		// 获取Redis密码,需要具备相应权限和授权
		password := "Trusfort@20151010"

		redis := RedisInfo{
			IP:       ip,
			Port:     "6379",
			Password: password,
		}

		redisList = append(redisList, redis)
		bar.Increment()
		time.Sleep(50 * time.Millisecond)
	}
	bar.Finish()

	// 如果存在Redis实例,则将IP地址、端口和密码写入JSON文件
	if len(redisList) > 0 {
		file, err := os.Create("multi_redis1.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		err = encoder.Encode(redisList)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		fmt.Println("Redis information saved to multi_redis1.json")
	} else {
		fmt.Println("No Redis instance found")
	}
}
