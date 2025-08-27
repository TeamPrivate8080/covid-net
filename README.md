# Covid V2 Driven Network
> ⚠️ Please don't attack websites without the owner's consent.  

This repository contains two main engines:

1. **Core Attack Tool** (`go run main.go`) – The core engine of the Covid network controller.  
2. **Bot Scanner** (`go run BotScan.go`) – Scans and analyzes your proxy bots worldwide. You can stop the bot scanner at any time, each bot will be stored immediately once it is scanned and works successfully
---

## Features & Methods

* 💣 Layer4

| Method  | Description |
|---------|------------|
| TCP     | Strong TCP Flood with multiplexed packets & connections. Work great for SSH, Web, RDP ports - (Power depends on your server & bots)
| UDP-REF | [SOON] Reflected UDP flooder with high bandwidth throughput

* 🧨 Layer7:

| Method    | Description |
|-----------|------------|
| HTTP-SPAM | [SOON] Powerful multiplexed rotating bot flood with a high number of legitimate requests (Up to 40 million requests a second on a dedicated)


## Requirements

- Go 1.20+ installed
- UNIX-like environment (Linux, macOS, or Windows with WSL)
- Input file: `proxies.txt` (non-scanned)
- Output file: `working.txt` (after scanning)

---

## Fresh Installation (Linux / WSL)

Run these commands to install dependencies, clone the repo, and prepare the environment:

```bash
# Update package lists
sudo apt update && sudo apt upgrade -y

# Install Git if not installed
sudo apt install git -y

# Install wget & tar for Go installation
sudo apt install wget tar -y

# Download Go 1.20+ (update version if needed)
wget https://go.dev/dl/go1.20.10.linux-amd64.tar.gz

# Remove old Go installation if exists
sudo rm -rf /usr/local/go

# Extract Go to /usr/local
sudo tar -C /usr/local -xzf go1.20.10.linux-amd64.tar.gz

# Set Go environment variables
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version

# Clone the repository
git clone https://github.com/TeamPrivate8080/covid-net.git
cd covid-net

# Download dependencies
go mod tidy
```

<p float="left">
  <img src="https://github.com/user-attachments/assets/78fa0f68-7cd3-4e05-b549-d25ba7b0ff2f" width="49%" />
  <img src="https://github.com/user-attachments/assets/d216a2cf-b947-43c5-ab8d-07aaa1617a93" width="49%" />
</p>
