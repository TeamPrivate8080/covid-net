# Covid V2 Driven Network

**Please Don't Attack websites without the owner's consent.**  

This repository contains two main tools:

1. **Attack Tool** (`go run main.go`) – The core engine of the Covid network controller.  
2. **Bot Scanner** (`go run BotScan.go`) – Scans and analyzes your proxy bots worldwide.

> ⚠️ **Warning:** Only use these tools for local testing or on systems you own with explicit permission. Unauthorized use is illegal.

## Features And Methods

 * 💣 Layer4

   * <img src="https://raw.githubusercontent.com/kgretzky/pwndrop/master/media/pwndrop-logo-512.png" width="16" height="16" alt="tcp"> TCP | Strong TCP Flood multiplexed packets & connections, works with scanned bots by default.


## Requirements

- Go 1.20+ installed
- UNIX-like environment (Linux, macOS, or Windows with WSL)
- Input file: `proxies.txt` (non-scanned)
- Output file: `working.txt` (after scanning)

## Installation

Clone the repository and download dependencies:

```bash
git clone https://github.com/bliphotelnl-hash/covid-net.git
cd covid-net
go mod tidy
