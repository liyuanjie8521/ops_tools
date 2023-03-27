package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	// "github.com/shirou/gopsutil/process"
)

func collect() {
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	totalPercent, _ := cpu.Percent(3*time.Second, false)
	perPercents, _ := cpu.Percent(3*time.Second, true)
	v, _ := mem.VirtualMemory()
	d, _ := disk.Usage("/")
	n, _ := host.Info()
	nv, _ := net.IOCounters(true)
	boottime, _ := host.BootTime()
	btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")

	if len(c) > 1 {
		for _, sub_cpu := range c {
			modelname := sub_cpu.ModelName
			cores := sub_cpu.Cores
			fmt.Printf(" CPU info   : %v %v cores \n", modelname, cores)
		}
	} else {
		sub_cpu := c[0]
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		fmt.Printf("CPU: %v %v cores \n", modelname, cores)
	}
	fmt.Printf("physical count: %d logical count: %d\n", var)
	fmt.Printf("  CPU Used   : used %f%% \n", cc[0])
	fmt.Printf(" Mem info   : Total: %v GB Free: %v GB Used:%v Usage:%f%%\n", v.Total/1024/1024/1024, v.Available/1024/1024/1024, v.Used/1024/1024/1024, v.UsedPercent)

	fmt.Printf("  Network    : %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)
	fmt.Printf("  SystemBoot : %v\n", btime)
	fmt.Printf("  HD         : Total: %v GB Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	fmt.Printf("  OS         : %v(%v) %v \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	fmt.Printf("  Hostname   : %v \n", n.Hostname)
}

func main() {
	collect()
}
