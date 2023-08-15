package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	// "github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	// "github.com/shirou/gopsutil/process"
)

func collect() {
	c, _ := cpu.Info()
	// cpu使用率
	cc, _ := cpu.Percent(time.Second, false)
	// cpu load值
	avg, _ := load.Avg()
	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	// totalPercent, _ := cpu.Percent(3*time.Second, false)
	// perPercents, _ := cpu.Percent(3*time.Second, true)
	// 物理内存
	v, _ := mem.VirtualMemory()
	// 交换内存
	swapMemory, _ := mem.SwapMemory()
	d, _ := disk.Usage("/")
	// 磁盘分区
	// partitions, _ := disk.Partitions(true)
	// 磁盘分区IO信息
	// counters, _ := disk.IOCounters()

	// 机器信息
	n, _ := host.Info()
	// 终端用户
	users, _ := host.Users()
	nv, _ := net.IOCounters(true)
	// 机器启动时间戳
	boottime, _ := host.BootTime()
	btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")

	// 所有进程名称和PID
	// processes, _ := process.Processes()
	// dockerID列表
	// list, _ := docker.GetDockerIDList()

	fmt.Println()
	fmt.Printf("  OS         : %v(%v) %v \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	for _, user := range users {
		fmt.Println("  All users  :", user.User)
	}

	fmt.Printf("  Hostname   : %v \n", n.Hostname)
	fmt.Println()

	/*
		// 打印CPU使用率,每5秒一次,总共1次
		for i := 1; i < 2; i++ {
			time.Sleep(time.Millisecond * 5000)
			percent, _ := cpu.Percent(time.Second, false)
			fmt.Printf("num: %v  cpu percent: %v", i, percent)
			fmt.Println()
		}
	*/
	fmt.Println()
	// 显示cpu load值:
	fmt.Println(avg)
	fmt.Println()

	fmt.Println("  CPU Information:")
	if len(c) > 1 {
		for _, sub_cpu := range c {
			modelname := sub_cpu.ModelName
			cores := sub_cpu.Cores
			fmt.Printf("    CPU info   : %v %v cores \n", modelname, cores)
		}
	} else {
		sub_cpu := c[0]
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		fmt.Printf("    CPU info   : %v %v cores \n", modelname, cores)
	}
	fmt.Println()
	fmt.Printf("    Physical : %d \n", physicalCnt)
	fmt.Printf("    Logical  : %d \n", logicalCnt)
	fmt.Println()
	fmt.Printf("    CPU percent : %f%% \n", cc[0])
	fmt.Println()

	fmt.Println("  Mem Information:")
	// 显示物理内存信息
	fmt.Printf("    Total: %v GB\n", v.Total/1024/1024/1024)
	fmt.Printf("    Free: %v GB\n", v.Available/1024/1024/1024)
	fmt.Printf("    Used:%v\n", v.Used/1024/1024/1024)
	fmt.Printf("    Usage:%f%%\n", v.UsedPercent)
	fmt.Println()
	// 显示交换内存
	fmt.Printf("    Swap Memory : %v\n", swapMemory)
	fmt.Println()

	fmt.Println("  Disk Information:")
	fmt.Printf("    HD         : Total: %v GB Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)

	/*
			// 打印磁盘分区信息
			for _, part := range partitions {
				fmt.Printf("part:%v\n", part.String())
				usage, _ := disk.Usage(part.Mountpoint)
				fmt.Printf("disk info:   used : %v  free: %v\n", usage.UsedPercent, usage.Free)
			}
		// 打印磁盘分区IO信息
		for k, v := range counters {
			fmt.Printf("%v,%v\n", k, v)
		}
	*/

	fmt.Println("  Network I/O Information:")
	fmt.Printf("    Network    : %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)

	fmt.Println()
	/*
		// 打印所有进程名称和PID
		for _, process := range processes {
			fmt.Println(process.Pid)
			// fmt.Println(process.Name())
		}

		fmt.Println()
			// 打印docker ID列表
			for _, v := range list {
				fmt.Println(v)
			}
	*/

	fmt.Printf("    SystemBoot : %v\n", btime)
	fmt.Println()
	fmt.Println(" Os Resources Collect is finished.please view informations.")
	fmt.Println()
}

func main() {
	collect()
}
