package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Metrics struct {
	OS       string `json:"os"`
	CPUs     int    `json:"cpus"`
	AllocRAM uint64 `json:"alloc_ram"`
}

type App struct {
	DB *sql.DB
}

func (a *App) metricsHandler(w http.ResponseWriter, r *http.Request) {
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
	dsn := "postgres://postgres:12345678@localhost:5432/sysmon"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешное подключение к БД!")
	app := App{
		DB: db,
	}

	http.HandleFunc("/api/metrics", app.metricsHandler)
	fmt.Println("Server started succesfully on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
