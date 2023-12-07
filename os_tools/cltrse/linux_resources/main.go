package main

//用于提供系统信息、操作系统的I/O、网络、内存、CPU等信息的获取和处理功能
import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"

	// "github.com/shirou/gopsutil/load"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	// "runtime"
)

const (
	user = iota
	nice
	system
	idle
	iowait
	irq
	softirq
	steal
	guest
	guest_nice
)

func printOSTime() {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("Current Time: %s", currentTime)
}

// 获取主机信息
func printLinuxHostInfo() {
	info, _ := host.Info()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"os items", "os information"})
	table.Append([]string{"Hostname", info.Hostname})
	table.Append([]string{"OS", info.OS})
	table.Append([]string{"Platform", info.Platform})
	table.Append([]string{"KernelVersion", info.KernelVersion})
	uptime := secondsToDurationString(info.Uptime)
	table.Append([]string{"Uptime", uptime})
	table.Render()
}

// 采集CPU相关信息
func printCPUInfo() {
	info, _ := cpu.Info()

	// 创建map用于合并和去重型号和核心数
	cpuInfo := make(map[string]int)

	for i := 0; i < len(info); i++ {
		modelName := info[i].ModelName
		cores := int(info[i].Cores)
		cpuInfo[modelName] += cores
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CPU_ModelName", "Cores"})

	for modelName, cores := range cpuInfo {
		table.Append([]string{modelName, fmt.Sprintf("%d", cores)})
	}

	table.Render()
}

// 统计CPU使用率相关信息
// func printCPUUsageInfo() {
// 	percent,_ := cpu.Percent(0, false)
// 	table := tablewriter.NewWriter(os.Stdout)
// 	table.SetHeader([]string{"CPU","Usage"})
// 	for i, p := range percent{
// 		table.Append([]string{fmt.Sprintf("CPU%d",i),fmt.Sprintf("%.2f%%", p)})
// 	}
// 	table.Render()
// }

// 获取CPU总核心数的使用率情况
func getCPUUseRate(sampleTime time.Duration) (float64, error) {
	workPre, totalPre, err := jiffies()
	if err != nil {
		return 0, err
	}

	time.Sleep(sampleTime)

	workAfter, totalAfter, err := jiffies()
	if err != nil {
		return 0, err
	}

	work := workAfter - workPre
	total := totalAfter - totalPre

	return float64(work) / float64(total), nil
}

var ErrInvalidStatFile = errors.New("invalid statistic file")

func jiffies() (workTime, totalTime int, err error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}

	lines := strings.SplitN(string(contents), "\n", 2)
	if len(lines) == 0 {
		return 0, 0, ErrInvalidStatFile
	}

	fields := strings.Fields(lines[0])
	if len(fields) != 11 || fields[0] != "cpu" {
		return 0, 0, ErrInvalidStatFile
	}

	v := make([]int, len(fields))
	for i := 1; i < len(fields); i++ {
		val, err := strconv.Atoi(fields[i])
		if err != nil {
			return 0, 0, fmt.Errorf("%w: %v", ErrInvalidStatFile, err)
		}
		v[i] = val
	}

	workTime = v[user] + v[nice] + v[system] + v[irq] + +v[softirq] + v[steal]
	idleTime := v[idle] + v[iowait]

	return workTime, workTime + idleTime, nil
}

// 分别获取CPU每个核心数的使用率信息(分核版)
func jiffiesPerCore() (workTimes, totalTimes []int, err error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(contents), "\n")
	if len(lines) == 0 {
		return nil, nil, ErrInvalidStatFile
	}

	coreValues := make([][]int, 0)
	for _, line := range lines[1:] {
		if !strings.HasPrefix(line, "cpu") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 11 {
			return nil, nil, ErrInvalidStatFile
		}

		v := make([]int, len(fields))
		for i := 1; i < len(fields); i++ {
			val, err := strconv.Atoi(fields[i])
			if err != nil {
				return nil, nil, fmt.Errorf("%w: %v", ErrInvalidStatFile, err)
			}
			v[i] = val
		}
		coreValues = append(coreValues, v)
	}

	for _, v := range coreValues {
		workTime := v[user] + v[nice] + v[system] + v[irq] + +v[softirq] + v[steal]
		idleTime := v[idle] + v[iowait]
		workTimes = append(workTimes, workTime)
		totalTimes = append(totalTimes, workTime+idleTime)
	}

	return workTimes, totalTimes, nil
}

func getCPUUseRatePerCore(sampleTime time.Duration) (rates []float64, err error) {
	workPre, totalPre, err := jiffiesPerCore()
	if err != nil {
		return nil, err
	}

	time.Sleep(sampleTime)

	workAfter, totalAfter, err := jiffiesPerCore()
	if err != nil {
		return nil, err
	}

	if len(workPre) != len(workAfter) {
		return nil, errors.New("unexpected jiffies")
	}

	for i, _ := range workPre {
		work := workAfter[i] - workPre[i]
		total := totalAfter[i] - totalPre[i]
		rate := float64(work) / float64(total)
		rates = append(rates, rate)
	}

	return rates, nil
}

// 内存信息
func printMemoryInfo() {
	info, _ := mem.VirtualMemory()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Total", "Available", "Used", "Used Percentage"})
	table.Append([]string{fmt.Sprintf("%d bytes", info.Total), fmt.Sprintf("%d bytes", info.Available), fmt.Sprintf("%d bytes", info.Used), fmt.Sprintf("%.2f%%", info.UsedPercent)})
	table.Render()
}

// 统计内存使用率信息
func printMemoryUsageInfo() {
	vmem, _ := mem.VirtualMemory()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Total", "Available", "Used", "Used Percentage"})
	table.Append([]string{bytesToString(vmem.Total), bytesToString(vmem.Available), bytesToString(vmem.Used), fmt.Sprintf("%.2f%%", vmem.UsedPercent)})
	table.Render()
}

// 获取磁盘信息
func printDiskInfo() {
	partitions, _ := disk.Partitions(false)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Device", "Mountpoint", "Fstype", "Total", "Free", "Used", "Used Percentage"})

	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		table.Append([]string{partition.Device, partition.Mountpoint, usage.Fstype, bytesToString(usage.Total), bytesToString(usage.Free), bytesToString(usage.Used), fmt.Sprintf("%.2f%%", usage.UsedPercent)})
	}
	table.Render()
}

// 统计硬盘使用率信息
func printDiskUsageInfo() {
	partitions, _ := disk.Partitions(false)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Device", "Mountpoint", "Fstype", "Total", "Free", "Used", "Used Percentage"})
	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		table.Append([]string{partition.Device, partition.Mountpoint, usage.Fstype, bytesToString(usage.Total), bytesToString(usage.Free), bytesToString(usage.Used), fmt.Sprintf("%.2f%%", usage.UsedPercent)})
	}
	table.Render()
}

// 统计网络IO信息
func printNetworkIOStats() {
	stats, _ := net.IOCounters(true)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Interface", "Download", "Upload"})
	for _, stat := range stats {
		table.Append([]string{stat.Name, fmt.Sprintf("%s/s", bytesToString(uint64(stat.BytesRecv))), fmt.Sprintf("%s/s", bytesToString(uint64(stat.BytesSent)))})
	}
	table.Render()
}

func secondsToDurationString(seconds uint64) string {
	minutes := (seconds / 60) % 60
	hours := (seconds / (60 * 60)) % 24
	days := seconds / (60 * 60 * 24)

	return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
}

func bytesToString(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	fmt.Println("\n\n\n====================  os system resources information  ===============================================")
	printOSTime()
	fmt.Println("\n\nLinux Host Information:")
	fmt.Println("Linux主机信息:")
	printLinuxHostInfo()
	fmt.Println("\nCPU Information:")
	fmt.Println("CPU信息:")
	printCPUInfo()
	// fmt.Println("\nMemory Information:")
	// fmt.Println("内存信息:")
	// printMemoryInfo()
	// fmt.Println("\nDisk Information:")
	// fmt.Println("硬盘信息:")
	// printDiskInfo()

	// fmt.Println("\nCPU Usage Information:")
	// printCPUUsageInfo()
	ratess, err := getCPUUseRate(1 * time.Second)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println()
	fmt.Println("\nGet Current Linux Os CPU Total UseRate Information:")
	fmt.Println("获取当前linux系统中CPU总使用率情况:")
	fmt.Printf("CPU Total UseRate：%f%%\n", ratess*100)
	fmt.Println()

	rates, err := getCPUUseRatePerCore(1 * time.Second)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println()
	fmt.Println("\nGet Current Linux Os CPU UseRate Per Core Information:")
	fmt.Println("获取当前linux系统中CPU每个核心数的使用率情况:")
	for i, rate := range rates {
		fmt.Printf("CPU %d UseRate：%f%%\n", i+1, rate*100)
	}
	fmt.Println()

	fmt.Println("\nMemory Usage Information:")
	printMemoryUsageInfo()
	fmt.Println("\nDisk Usage Information:")
	printDiskUsageInfo()
	fmt.Println("\nNetwork Information:")
	printNetworkIOStats()

	fmt.Println("\n\n\n====================  os system resources information collect finished  ============================")
	fmt.Println("\n\n\n")
}
