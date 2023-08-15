package main

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"

	// "github.com/vmware/govmomi/performance"
	"net/url"
	// "github.com/vmware/govmomi/property"
	// "github.com/vmware/govmomi/view"
	"log"
)

func main() {
	username := "root"
	password := "123456"
	// esxiServers := []string{"192.168.1.156"}
	esxiServers := "192.168.1.156"

	// for _, server := range esxiServers {
	/*
		client, err := govmomi.NewClient(context.Background(), &url.URL{
			Scheme: "https",
			Host:   server,
			Path:   "/sdk",
		}, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(client)
		defer client.Logout(context.Background())
	*/

	u := &url.URL{
		Scheme: "https",
		Host:   esxiServers,
		Path:   "/sdk",
	}
	ctx := context.Background()
	u.User = url.UserPassword(username, password)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		// panic(err)
		log.Fatalf("Failed to connect to %q: %v", esxiServers, err)
	}
	// fmt.Println(client)
	defer client.Logout(ctx)
	// Set up finder to explore the inventory hierarchy
	finder := find.NewFinder(client.Client, true)

	// Set up view manager to retrieve performance counters
	//m := performance.NewManager(client.Client)

	// Retrieve host system object from inventory
	host, err := finder.HostSystem(ctx, esxiServers)
	if err != nil {
		log.Fatalf("Failed to find host %q: %v", esxiServers, err)
	}
	fmt.Println(host)
	// }
}
