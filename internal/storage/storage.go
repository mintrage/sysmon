package storage

import (
	"database/sql"

	"github.com/mintrage/sysmon/internal/models"
)

type Storage struct {
	DB *sql.DB
}

func (s *Storage) SaveMetric(m models.Metrics) error {
	insertSQL := "INSERT INTO metrics (server_name, os, cpu_usage, alloc_ram) VALUES ($1, $2, $3, $4)"
	_, err := s.DB.Exec(insertSQL, m.ServerName, m.OS, m.CPUUsage, m.AllocRAM)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetLatestMetric() (models.Metrics, error) {
	selectSQL := "SELECT server_name, os, cpu_usage, alloc_ram FROM metrics ORDER BY id DESC LIMIT 1"
	var m models.Metrics
	row := s.DB.QueryRow(selectSQL)
	err := row.Scan(&m.ServerName, &m.OS, &m.CPUUsage, &m.AllocRAM)
	if err != nil {
		return m, err
	}
	return m, nil
}
