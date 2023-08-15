package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"go.etcd.io/etcd/clientv3"
)

func main() {
	now := time.Now().Format("2006-01-02 15:04:05")

	// color print
	info := color.New(color.FgGreen).SprintFunc()
	warn := color.New(color.FgYellow).SprintFunc()
	error := color.New(color.FgRed).SprintFunc()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.1.180:2379"},
		DialTimeout: time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed,err: %v\n", err)
		return
	}
	// info := color.New(color.BgHiBlack, color.FgGreen).SprintFunc()
	fmt.Printf("%s  %s\n", info(now), info("Info: connect etcd success."))
	fmt.Printf("%s  %s\n", warn(now), warn("warn: connect etcd success."))
	fmt.Printf("%s  %s\n", error(now), error("error: connect etcd success."))

	defer cli.Close()

	// Watch操作
	wch := cli.Watch(context.Background(), "/name")
	for resp := range wch {
		for _, ev := range resp.Events {
			fmt.Printf("Type: %v,Key:%v,Value:%v\n", ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}
}
