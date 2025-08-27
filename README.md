Please Don't Attack websites without the owner's consent.

# Covid V2 Driven Network

This repository contains two tools:

1. **Attack Tool** (`go run main.go`) – Main Core engine of the Covid network controller
2. **Bot Scanner** (`go run BotScan.go`) – Scans and analyzes your proxy bots. (world wide)

> ⚠️ **Warning:** Only use these tools for local testing and with explicit permission. Unauthorized use may be illegal.

---

## Requirements

- Go 1.20+ installed
- A UNIX-like environment (Linux, macOS, or Windows with WSL)
- Output after scanning bots will always be working.txt by default
- Input by proxies.txt (non-scanned)
---

## Installation

Clone the repository and download dependencies:

```bash
https://github.com/bliphotelnl-hash/covid-net.git
cd covid-net
go mod tidy

Commands:
go run main.go
go run BotScan.go
