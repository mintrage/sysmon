package storage

import (
	"database/sql"

	"github.com/mintrage/sysmon/internal/models"
)

type Storage struct {
	DB *sql.DB
}

func (s *Storage) SaveMetric(m models.Metrics) error {
	insertSQL := "INSERT INTO metrics (os, cpus, alloc_ram) VALUES ($1, $2, $3)"
	_, err := s.DB.Exec(insertSQL, m.OS, m.CPUs, m.AllocRAM)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetLatestMetric() (models.Metrics, error) {
	selectSQL := "SELECT os, cpus, alloc_ram FROM metrics ORDER BY id DESC LIMIT 1"
	var m models.Metrics
	row := s.DB.QueryRow(selectSQL)
	err := row.Scan(&m.OS, &m.CPUs, &m.AllocRAM)
	if err != nil {
		return m, err
	}
	return m, nil
}
