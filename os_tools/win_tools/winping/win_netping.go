package main

import (
	"flag"
	"fmt"
	// "gitee.com/liumou_site/gns"
	"gitee.com/liumou_site/logger"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

// 定义一个计数器
var wg sync.WaitGroup

func ping(host string) {
	//logger.Debug("开始测试", host)
	r := exec.Command("ping", host)
	err := r.Run()
	if err != nil {
		logger.Error("失败:", host)
	} else {
		logger.Info("成功:", host)
	}
	// 完成一个计数器
	wg.Done()
}

func main() {
	StartTime := time.Now()
	// 定义变量,用于接收命令行的参数值
	var host string
	// 要扫描的主机IP地址
	flag.StringVar(&host, "h", "192.168.1.", "scan target host ip network default net 192.168.1.")
	// 转换
	flag.Parse()

	// 添加254个任务
	wg.Add(254)
	// 生成254个IP(切片)
	/*
		IpSub, err := gns.IpGenerateList("192.168.1", 1, 254)
		if err == nil {
			for _, ip := range IpSub {
				go ping(ip)
			}
		}
	*/
	fmt.Println()
	fmt.Println()
	fmt.Println("开始扫描指定网段的所有IP地址,请查看扫描结果")
	fmt.Println()
	for i := 1; i <= 254; i++ {
		//fmt.Println(ip + strconv.Itoa(i))
		true_ip := host + strconv.Itoa(i)
		go ping(true_ip)
	}
	// 等待所有任务完成
	wg.Wait()
	fmt.Println()
	fmt.Println("已完成扫描指定网段的所有IP地址")
	fmt.Println()
	UseTime := time.Since(StartTime)
	fmt.Println("    耗时时长为:", UseTime, " 秒")

}
