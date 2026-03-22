package models

// Структура должна быть с большой буквы, чтобы другие пакеты могли ее импортировать
type Metrics struct {
	ServerName string  `json:"server_name"`
	OS         string  `json:"os"`
	CPUUsage   float64 `json:"cpu_usage"`
	AllocRAM   uint64  `json:"alloc_ram"`
}
