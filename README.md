/*
Package main provides a command-line interface (CLI) tool for searching through AWS CloudWatch logs. This tool is designed to help users efficiently query and filter log data from their AWS CloudWatch logs, making it easier to diagnose issues, monitor applications, and analyze log data.

## Installation

To install the CLI tool, you need to download the appropriate release binary for your CPU and operating system architecture. Follow the instructions below for your respective environment:

### Bash (Linux/macOS)

1. Determine your operating system and architecture:
    ```bash
    uname -s
    uname -m
    ```

2. Download the release binary:
    ```bash
    OS=$(uname -s)
    ARCH=$(uname -m)
    curl -Lo cloudwatch-search https://github.com/your-repo/cloudwatch-search/releases/latest/download/cloudwatch-search-$OS-$ARCH
    ```

3. Make the binary executable:
    ```bash
    chmod +x cloudwatch-search
    ```

4. Move the binary to a directory in your PATH:
    ```bash
    sudo mv cloudwatch-search /usr/local/bin/
    ```

### Windows PowerShell

1. Determine your operating system and architecture:
    ```powershell
    $OS = (Get-CimInstance Win32_OperatingSystem).Caption
    $ARCH = (Get-CimInstance Win32_Processor).Architecture
    ```

2. Download the release binary:
    ```powershell
    $url = "https://github.com/your-repo/cloudwatch-search/releases/latest/download/cloudwatch-search-$OS-$ARCH.exe"
    Invoke-WebRequest -Uri $url -OutFile cloudwatch-search.exe
    ```

3. Move the binary to a directory in your PATH:
    ```powershell
    Move-Item -Path .\cloudwatch-search.exe -Destination "C:\Program Files\"
    [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\Program Files\", [System.EnvironmentVariableTarget]::Machine)
    ```

## Usage

This CLI tool is essential for developers and system administrators who need to search and analyze log data from AWS CloudWatch. It simplifies the process of querying logs, allowing users to filter logs based on various criteria such as log group, log stream, time range, and specific patterns.

### Example Commands


By using this CLI tool, users can quickly and efficiently access and analyze their CloudWatch logs, leading to faster issue resolution and better insights into their applications and systems.
*/