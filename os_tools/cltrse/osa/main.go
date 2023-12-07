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


func main() {
	rates, err := getCPUUseRate(1 * time.Second)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("当前CPU利用率：%f%%\n", rates*100)
	fmt.Println()
}