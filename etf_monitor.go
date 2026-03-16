package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	// 配置参数
	etfCode := "sh515210" // 钢铁ETF代码
	maxPrice := 1.589      // 最大价格阈值，请根据实际情况调整
	minPrice := 1.41      // 最小价格阈值，请根据实际情况调整
	checkInterval := 3   // 检查间隔（秒）

	fmt.Printf("开始监控钢铁ETF (%s)...\n", etfCode)
	fmt.Printf("价格阈值: 最低 %.2f, 最高 %.2f\n", minPrice, maxPrice)
	fmt.Printf("检查间隔: %d 秒\n", checkInterval)
	fmt.Println("按 Ctrl+C 停止监控")

	// 主监控循环
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		price, err := getETFTPrice(etfCode)
		if err != nil {
			fmt.Printf("获取价格失败: %v\n", err)
			continue
		}

		fmt.Printf("%s - 当前价格: %.3f\n", time.Now().Format("15:04:05"), price)

		// 检查阈值
		if price > maxPrice {
			fmt.Printf("⚠️ 警告: 价格 %.3f 超过最大值 %.2f\n", price, maxPrice)
			playAlertSound()
		} else if price < minPrice {
			fmt.Printf("⚠️ 警告: 价格 %.3f 低于最小值 %.2f\n", price, minPrice)
			playAlertSound()
		}
	}
}

// getETFTPrice 获取ETF实时价格
func getETFTPrice(code string) (float64, error) {
	// 构建新浪API URL
	url := fmt.Sprintf("http://hq.sinajs.cn/rn=%d&list=%s", time.Now().Unix(), code)

	// 创建HTTP请求，添加必要的请求头
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加请求头模仿浏览器访问
	req.Header.Set("Host", "hq.sinajs.cn")
	req.Header.Set("Referer", "https://finance.sina.com.cn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %w", err)
	}

	// 将GB18030编码转换为UTF-8（新浪API使用GB18030编码）
	content, err := GB18030ToUTF8(body)
	if err != nil {
		// 如果转换失败，尝试直接使用原始内容
		content = string(body)
	}

	// 解析价格数据
	price, err := parseSinaPriceResponse(content)
	if err != nil {
		return 0, fmt.Errorf("解析价格失败: %w", err)
	}

	return price, nil
}

// parseSinaPriceResponse 解析新浪API的响应
func parseSinaPriceResponse(response string) (float64, error) {
	// 响应格式: var hq_str_sh515210="钢铁ETF,1.234,1.200,...";
	// 找到等号分割
	parts := strings.SplitN(response, "=", 2)
	if len(parts) < 2 {
		return 0, fmt.Errorf("无效的响应格式，找不到等号")
	}

	// 提取等号右边的部分，去除前后的空格和引号
	data := strings.TrimSpace(parts[1])
	// 去除末尾的分号
	data = strings.TrimSuffix(data, ";")
	// 去除前后的双引号
	data = strings.Trim(data, "\"")

	// 按逗号分割字段
	fields := strings.Split(data, ",")
	if len(fields) < 4 {
		return 0, fmt.Errorf("数据字段不足，期望至少4个字段，实际%d个", len(fields))
	}

	// 第三个字段是当前价格（索引3）
	priceStr := fields[3]
	if priceStr == "" {
		return 0, fmt.Errorf("价格字段为空")
	}

	// 转换为浮点数
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("价格转换失败 '%s': %w", priceStr, err)
	}

	return price, nil
}

// GB18030ToUTF8 将GB18030编码转换为UTF-8
func GB18030ToUTF8(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

// playAlertSound 播放警报声音
func playAlertSound() {
	// 方法1: 使用控制台响铃（可能在某些终端中不起作用）
	fmt.Print("\a")

	// 方法2: 使用PowerShell播放系统声音（Windows）
	cmd := exec.Command("powershell", "-c", "[System.Media.SystemSounds]::Beep.Play()")
	err := cmd.Run()
	if err != nil {
		// 如果PowerShell失败，尝试使用简单的蜂鸣
		fmt.Println("声音播放失败，请检查系统声音设置")
	}
}

// 初始化函数，检查操作系统
func init() {
	// 检查是否是Windows系统
	if os.Getenv("OS") != "Windows_NT" {
		fmt.Println("警告: 此程序主要在Windows系统上测试")
		fmt.Println("非Windows系统可能需要调整声音播放方式")
	}
}