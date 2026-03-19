package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
		jsonData, err := json.Marshal(collectMetrics())
		if err != nil {
			fmt.Println("Error marshalling to json")
			return
		}
		//fmt.Println(string(jsonData))
		resp, err := http.Post("http://localhost:8080/api/metrics", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Не удалось отправить данные")
		} else {
			resp.Body.Close()
			fmt.Println("Метрики отправлены, статус:", resp.Status)
		}
		time.Sleep(2 * time.Second)
	}
}
