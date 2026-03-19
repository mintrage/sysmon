package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Metrics struct {
	OS       string `json:"os"`
	CPUs     int    `json:"cpus"`
	AllocRAM uint64 `json:"alloc_ram"`
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}
	var m Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	fmt.Printf("Получены метрики: %+v\n", m)
}

func main() {
	http.HandleFunc("/api/metrics", metricsHandler)
	fmt.Println("Server started succesfully")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server")
	}

}
