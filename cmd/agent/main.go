package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

type Metrics struct {
	ServerName string `json:"server_name"`
	OS         string `json:"os"`
	CPUs       int    `json:"cpus"`
	AllocRAM   uint64 `json:"alloc_ram"`
}

func collectMetrics(name string) Metrics {

	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("")
		return Metrics{
			ServerName: name,
			OS:         runtime.GOOS,
			CPUs:       runtime.NumCPU(),
			AllocRAM:   0,
		}
	}

	return Metrics{
		ServerName: name,
		OS:         runtime.GOOS,
		CPUs:       runtime.NumCPU(),
		AllocRAM:   v.Used,
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
