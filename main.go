package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Statistics struct {
	VisitCount   int    `json:"visit_count"`
	TotalTraffic string `json:"total_traffic"`
	UpdateTime   string `json:"update_time"`
	Start        string `json:"start"`
	End          string `json:"end"`
}

func formatTraffic(bytes int64) string {
	if bytes >= 1024*1024*1024 {
		return fmt.Sprintf("%.2fG", float64(bytes)/(1024*1024*1024))
	}
	if bytes >= 1024*1024 {
		return fmt.Sprintf("%.2fM", float64(bytes)/(1024*1024))
	}
	if bytes >= 1024 {
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	}
	return fmt.Sprintf("%dB", bytes)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("用法: %s <日志文件>\n", os.Args[0])
		os.Exit(1)
	}

	file, _ := os.Open(os.Args[1])
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("关闭文件时出错: %v\n", err)
			os.Exit(1)
		}
	}(file)

	totalRequests := 0
	totalTraffic := int64(0)
	today := time.Now().Format("02/Jan/2006")
	dateRegex := regexp.MustCompile(today)
	trafficRegex := regexp.MustCompile(`/bmclapi/[^"]+" 200 (\d+)`)
	timeRegex := regexp.MustCompile(`\[(\d{2}/[A-Za-z]{3}/\d{4}:\d{2}:\d{2}:\d{2})`)

	foundToday := false
	var startTime, endTime time.Time

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !foundToday {
			if !dateRegex.MatchString(line) {
				continue
			}
			foundToday = true
		}

		// 记录时间范围
		if timeMatch := timeRegex.FindStringSubmatch(line); len(timeMatch) == 2 {
			t, err := time.Parse("02/Jan/2006:15:04:05", timeMatch[1])
			if err == nil {
				if startTime.IsZero() {
					startTime = t
				}
				endTime = t
			}
		}

		matches := trafficRegex.FindStringSubmatch(line)
		if len(matches) == 2 {
			totalRequests++
			traffic, _ := strconv.ParseInt(matches[1], 10, 64)
			totalTraffic += traffic
		}
	}

	stats := Statistics{
		VisitCount:   totalRequests,
		TotalTraffic: formatTraffic(totalTraffic),
		UpdateTime:   time.Now().Format("2006/01/02 15:04:05"),
		Start:        startTime.Format("2006/01/02 15:04:05"),
		End:          endTime.Format("2006/01/02 15:04:05"),
	}

	jsonData, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Println(string(jsonData))
}
