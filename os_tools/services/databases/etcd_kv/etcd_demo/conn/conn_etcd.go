package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.1.180:2379"},
		DialTimeout: 5 * time.Second})
	if err != nil {
		fmt.Printf("init etcd failed,err:%v\n", err)
		return
	}
	fmt.Println("init etcd success.")

	defer cli.Close()
	// pub操作
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "/name", "jerry")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed,err: %v\n", err)
		return
	}
	fmt.Println("put to etcd success.")

	// get操作

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/name")
	cancel()

	if err != nil {
		fmt.Println("get from etcd failed, err:", err)
		return
	}

	for _, v := range resp.Kvs {
		fmt.Printf("Key:%s,Value:%s\n", v.Key, v.Value)
	}
}
