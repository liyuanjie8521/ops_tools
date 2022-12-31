package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/er10yi/nmap-go/nmap"
)

// nmap基础扫描
func main() {
	// 定义变量,用于接收命令行的参数值
	var host string
	var ports string

	// 要扫描的主机IP地址
	flag.StringVar(&host, "h", "127.0.0.1", "scan target host ip address,default 127.0.0.1")
	// 要扫描的主机端口
	flag.StringVar(&ports, "p", "3306,6379,9200,9300", "<port ranges> Only scan specified ports,default scan mysql、redis、es services always in use ports.")

	// 转换
	flag.Parse()

	// 通过NewNmap()创建nmap
	// AddTargets增加目标,Addp增加端口范围
	scanner := nmap.NewNmap().AddTargets(host).Addp(ports)

	// Run运行
	runResult := scanner.Run()

	// 获取警告信息
	warn := runResult.WarnOut
	if warn != "" {
		fmt.Printf("warn:\n%s", warn)
	}

	// 获取错误信息
	err := runResult.ErrOut
	if err != nil {
		log.Fatal("error:", err)
	}

	// 获取运行的xml结果
	result := runResult.Result

	// 解析xml结果
	parseResult := scanner.ParseXmlResult(result)
	if err != nil {
		log.Fatal(err)
	}

	xmlResult := parseResult.(*nmap.NmapXMLResult)

	// 格式化输出xml结果
	scanner.PrettyResult(xmlResult)

	// 导出xml结果到Excel
	scanner.ExportResult(xmlResult)

	// 导出xml结果到txt, 用于导入魔方
	scanner.ExportTxtResult(xmlResult)
}
