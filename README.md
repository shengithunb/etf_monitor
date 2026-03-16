ETF Monitor
一个用于监控ETF（交易型开放式指数基金）行情的轻量级工具。

项目简介
ETF Monitor 是一个用 Go 语言编写的命令行工具，用于实时监控和分析ETF基金的市场数据。它可以帮助投资者快速获取ETF的实时行情、历史数据和技术指标。
安装使用
方法一：直接下载可执行文件
从 Releases 页面下载对应平台的可执行文件：
Windows: etf-monitor.exe

方法二：从源码构建

构建步骤
bash
# 1. 克隆仓库
git clone https://github.com/shengithunb/etf_monitor.git
cd etf_monitor


# 3. 构建可执行文件
go build -o etf-monitor ./cmd/etf_monitor

# 或者构建特定平台的版本
## Windows
GOOS=windows GOARCH=amd64 go build -o etf-monitor.exe ./cmd/etf_monitor

## Linux
GOOS=linux GOARCH=amd64 go build -o etf-monitor-linux ./cmd/etf_monitor

## macOS
GOOS=darwin GOARCH=amd64 go build -o etf-monitor-macos ./cmd/etf_monitor
