package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	var ips []string
	for i := 1; i <= 255; i++ {
		ip := fmt.Sprintf("192.168.1.%d", i)
		ips = append(ips, ip)
	}

	var wg sync.WaitGroup
	results := make(chan string)

	fmt.Println("Scanning 192.168.1.0/24 network for Redis instance...")
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", ip+":6379", time.Second)
			if err == nil {
				password := ""
				_, err = conn.Write([]byte("AUTH Trusfort@20151010\r\n"))
				if err == nil {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err == nil && string(buf[:n]) == "+OK\r\n" {
						password = "Trusfort@20151010"
					}
				}
				conn.Close()
				results <- fmt.Sprintf("{\"ip\":\"%s\",\"port\":%d,\"password\":\"%s\"}", ip, 6379, password)
			}
		}(ip)
	}

	go func() {
		var count int
		for result := range results {
			count++
			fmt.Printf("\rScanned %d/%d IP addresses...", count, len(ips))
			fmt.Println(result + ",")
		}
	}()

	wg.Wait()
	close(results)

	fmt.Println("\n\tDone.")
	fmt.Println("Writing results to file...")
	var data []map[string]interface{}
	for result := range results {
		var m map[string]interface{}
		json.Unmarshal([]byte(result), &m)
		data = append(data, m)
	}
	file, _ := os.Create("multi_redis.json")
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
	fmt.Println("Done.")
}
