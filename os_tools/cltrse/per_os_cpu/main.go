package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"strings"
	"errors"
	"strconv"
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

var ErrInvalidStatFile = errors.New("invalid statistic file")


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

func main() {
	go func() {
		for i := 0; ; {i++}
	}()
	for {
		rates, err := getCPUUseRatePerCore(1 * time.Second)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println()
		fmt.Println("分别获取当前linux系统中 CPU 每个核心的使用情况")
		for i, rate := range rates {
			fmt.Printf("当前CPU %d 利用率：%f%%\n", i+1, rate*100)
		}
		fmt.Println()
	}
}