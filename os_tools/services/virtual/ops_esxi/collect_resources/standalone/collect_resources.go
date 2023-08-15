package main

// 用于通过vSphere API连接和查询多个ESXi服务器的硬件资源和虚拟机信息

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	esxiUsername = "root"
	esxiPassword = "123456"
	esxiHostname = "192.168.1.156"
)

func main() {
	// Set up connection to ESXi server
	u := url.URL{
		Scheme: "https",
		Host:   esxiHostname,
		Path:   "/sdk",
	}
	ctx := context.Background()
	conn, err := govmomi.NewClient(ctx, &u, true)
	if err != nil {
		log.Fatalf("Failed to connect to %q: %v", esxiHostname, err)
	}
	defer conn.Logout(ctx)

	// Set up finder to explore the inventory hierarchy
	finder := find.NewFinder(conn.Client, true)

	// Set up view manager to retrieve performance counters
	m := performance.NewManager(conn.Client)

	// Retrieve host system object from inventory
	host, err := finder.HostSystem(ctx, esxiHostname)
	if err != nil {
		log.Fatalf("Failed to find host %q: %v", esxiHostname, err)
	}

	// Retrieve summary hardware info from host system
	hardware, err := host(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve hardware resource summary: %v", err)
	}
	fmt.Println("CPU count:", hardware.NumCpuCores)
	fmt.Println("Memory capacity (bytes):", hardware.MemoryCapacity)

	// Set up performance query spec to retrieve CPU usage over the past minute
	querySpec := performance.NewQuerySpec(host.Reference(), []performance.MetricId{
		performance.MetricId{
			CounterId:  6,
			InstanceId: "*",
		},
	}, &types.PerfQuerySpec{
		MaxSample:  1,
		IntervalId: 20,
	})
	querySpec.ParseRealtime()

	// Retrieve performance counters for the query spec
	results, err := m.Query(ctx, querySpec)
	if err != nil {
		log.Fatalf("Failed to retrieve performance counter: %v", err)
	}
	if len(results) == 0 {
		fmt.Println("No performance counter results.")
		os.Exit(0)
	}

	// Print CPU usage over the past minute
	for _, result := range results[0].Series {
		for _, value := range result.Value {
			fmt.Println("CPU usage over the past minute (%):", value.Value[0])
		}
	}

	// Retrieve VM objects from inventory
	vms, err := finder.VirtualMachineList(ctx, "*")
	if err != nil {
		log.Fatalf("Failed to retrieve virtual machine list: %v", err)
	}

	// Print summary info for each VM
	for _, vm := range vms {
		fmt.Println("VM name:", vm.Name())
		fmt.Println("VM power state:", vm.Runtime.PowerState)
		fmt.Println("VM CPU count:", vm.Config.Hardware.NumCPU)
		fmt.Println("VM memory (MB):", vm.Config.Hardware.MemoryMB)
	}

	// Log out of connection
	if err := conn.Logout(ctx); err != nil {
		log.Fatalf("Failed to log out of connection: %v", err)
	}
}
