package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mintrage/sysmon/internal/models"
)

// type Metrics struct {
// 	OS       string `json:"os"`
// 	CPUs     int    `json:"cpus"`
// 	AllocRAM uint64 `json:"alloc_ram"`
// }

type App struct {
	DB *sql.DB
}

func (a *App) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}
	var m models.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	insertSQL := "INSERT INTO metrics (os, cpus, alloc_ram) VALUES ($1, $2, $3)"
	_, err = a.DB.Exec(insertSQL, m.OS, m.CPUs, m.AllocRAM)
	if err != nil {
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *App) latestMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET", http.StatusMethodNotAllowed)
		return
	}
	selectSQL := "SELECT os, cpus, alloc_ram FROM metrics ORDER BY id DESC LIMIT 1"
	var m models.Metrics
	row := a.DB.QueryRow(selectSQL)
	err := row.Scan(&m.OS, &m.CPUs, &m.AllocRAM)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)

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
	sql := "CREATE TABLE IF NOT EXISTS metrics (id SERIAL PRIMARY KEY, os VARCHAR(50), cpus INT, alloc_ram BIGINT);"
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешное подключение к БД!")
	app := App{
		DB: db,
	}

	http.HandleFunc("/api/metrics", app.metricsHandler)
	http.HandleFunc("/api/metrics/latest", app.latestMetricsHandler)

	srv := &http.Server{
		Addr: ":8080",
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Ошибка сервера: %v", err)
				return
			}
		}
	}()

	fmt.Println("Server started succesfully on port 8080...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Получен сигнал на завершение, выключаем сервер...")

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err = srv.Shutdown(ctxWithTimeout)
	if err != nil {
		log.Fatal("Принудительное завершение сервера:", err)
	}

	fmt.Println("Сервер успешно и безопасно остановлен.")

}
