package main

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type Metrics struct {
	OS   string `json:"os"`
	CPUs int    `json:"cpus"`
}

func collectMetrics() Metrics {

	return Metrics{
		OS:   runtime.GOOS,
		CPUs: runtime.NumCPU(),
	}
}

func main() {
	jsonData, err := json.MarshalIndent(collectMetrics(), "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to json")
		return
	}
	fmt.Println(string(jsonData))
}
