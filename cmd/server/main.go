package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mintrage/sysmon/internal/handler"
	"github.com/mintrage/sysmon/internal/storage"
)

func main() {
	dsn := "postgres://postgres:12345678@db:5432/sysmon"
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

	store := &storage.Storage{DB: db}

	h := &handler.Handler{Storage: store}

	http.HandleFunc("/api/metrics", h.MetricsHandler)
	http.HandleFunc("/api/metrics/latest", h.LatestMetricsHandler)

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
