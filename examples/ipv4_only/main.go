package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
)

func main() {
	log.Println("Starting IPv4-only zeroconf example...")

	// 服务端：注册一个只使用IPv4的服务
	log.Println("Registering IPv4-only service...")
	server, err := zeroconf.RegisterWithOptions(
		"MyService",                          // 服务实例名
		"_myapp._tcp",                        // 服务类型
		"local.",                             // 域
		8080,                                 // 端口
		[]string{"version=1.0", "path=/api"}, // TXT记录
		nil,                                  // 网络接口（nil表示使用所有可用接口）
		zeroconf.SelectServerIPTraffic(zeroconf.IPv4), // 只使用IPv4
	)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	defer server.Shutdown()

	log.Println("Service registered successfully (IPv4 only)")

	// 客户端：创建一个只监听IPv4的解析器
	log.Println("Creating IPv4-only resolver...")
	resolver, err := zeroconf.NewResolver(
		zeroconf.SelectIPTraffic(zeroconf.IPv4), // 只使用IPv4
	)
	if err != nil {
		log.Fatalf("Failed to create resolver: %v", err)
	}

	// 浏览服务
	entries := make(chan *zeroconf.ServiceEntry)
	go func() {
		for entry := range entries {
			log.Printf("Found service: %s at %s:%d (IPv4: %v, IPv6: %v)",
				entry.Instance, entry.HostName, entry.Port,
				entry.AddrIPv4, entry.AddrIPv6)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Println("Browsing for services (IPv4 only)...")
	err = resolver.Browse(ctx, "_myapp._tcp", "local.", entries)
	if err != nil {
		log.Fatalf("Failed to browse: %v", err)
	}

	// 等待信号或超时
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("Received interrupt signal")
	case <-ctx.Done():
		log.Println("Context timeout")
	}

	log.Println("Shutting down...")
}
