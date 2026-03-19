package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Metrics struct {
	OS       string `json:"os"`
	CPUs     int    `json:"cpus"`
	AllocRAM uint64 `json:"alloc_ram"`
}

func collectMetrics() Metrics {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return Metrics{
		OS:       runtime.GOOS,
		CPUs:     runtime.NumCPU(),
		AllocRAM: m.Alloc,
	}
}

func main() {
	for {
		jsonData, err := json.MarshalIndent(collectMetrics(), "", "  ")
		if err != nil {
			fmt.Println("Error marshalling to json")
			return
		}
		fmt.Println(string(jsonData))
		time.Sleep(2 * time.Second)
	}
}
