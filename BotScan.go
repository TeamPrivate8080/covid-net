package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	inputFile := "proxies.txt"
	outputFile := "working.txt"
	testTarget := "8.8.8.8:53"
	timeout := 5 * time.Second
	concurrency := 50

	rand.Seed(time.Now().UnixNano())

	BaseFileDeduplicator(inputFile)
	BaseFileDeduplicator(outputFile)

	existing := make(map[string]struct{})
	if f, err := os.Open(outputFile); err == nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				existing[line] = struct{}{}
			}
		}
		f.Close()
	}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Failed to open %s: %v\n", inputFile, err)
		return
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			if _, ok := existing[line]; !ok {
				proxies = append(proxies, line)
			}
		}
	}

	if len(proxies) == 0 {
		fmt.Println("No new proxies to check.")
		return
	}

	rand.Shuffle(len(proxies), func(i, j int) { proxies[i], proxies[j] = proxies[j], proxies[i] })

	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, concurrency)

	for _, proxyAddr := range proxies {
		wg.Add(1)
		sem <- struct{}{}
		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()

			if CheckNewSocks5(p, testTarget, timeout) {
				mu.Lock()
				if _, ok := existing[p]; !ok {
					existing[p] = struct{}{}
					f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err == nil {
						f.WriteString(p + "\n")
						f.Close()
					}
					fmt.Println("[WORKING]", p)
				}
				mu.Unlock()
			} else {
				fmt.Println("[FAILED] ", p)
			}
		}(proxyAddr)
	}

	wg.Wait()
	fmt.Println("Proxy scan finished.")
}

func CheckNewSocks5(proxyAddr, target string, timeout time.Duration) bool {
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return false
	}

	connChan := make(chan net.Conn, 1)
	errChan := make(chan error, 1)

	go func() {
		conn, err := dialer.Dial("tcp", target)
		if err != nil {
			errChan <- err
		} else {
			connChan <- conn
		}
	}()

	select {
	case conn := <-connChan:
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(timeout))
		return true
	case <-time.After(timeout):
		return false
	case <-errChan:
		return false
	}
}

func BaseFileDeduplicator(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	lines := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines[line] = struct{}{}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to rewrite %s: %v\n", filename, err)
		return
	}
	defer f.Close()

	for line := range lines {
		f.WriteString(line + "\n")
	}
}
