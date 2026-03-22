package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	ServerName string  `json:"server_name"`
	OS         string  `json:"os"`
	CPUUsage   float64 `json:"cpu_usage"`
	AllocRAM   uint64  `json:"alloc_ram"`
}

func collectMetrics(name string) Metrics {

	v, errMem := mem.VirtualMemory()
	c, errCPU := cpu.Percent(0, false)
	var cpuUsage float64
	if errCPU == nil && len(c) > 0 {
		cpuUsage = c[0]
	}

	alloc := uint64(0)
	if errMem == nil {
		alloc = v.Used
	}

	return Metrics{
		ServerName: name,
		OS:         runtime.GOOS,
		CPUUsage:   cpuUsage,
		AllocRAM:   alloc,
	}
}

func main() {
	agentName := os.Getenv("SYSMON_AGENT_NAME")
	if agentName == "" {
		host, err := os.Hostname()
		if err != nil {
			agentName = "unknown-agent"
		} else {
			agentName = host
		}
	}

	serverURL := os.Getenv("SYSMON_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080/api/metrics"
	}

	for {
		jsonData, err := json.Marshal(collectMetrics(agentName))
		if err != nil {
			fmt.Println("Error marshalling to json")
			return
		}
		//fmt.Println(string(jsonData))
		resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Не удалось отправить данные")
		} else {
			resp.Body.Close()
			fmt.Println("Метрики отправлены, статус:", resp.Status)
		}
		time.Sleep(2 * time.Second)
	}
}
