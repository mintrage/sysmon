package models

// Структура должна быть с большой буквы, чтобы другие пакеты могли ее импортировать
type Metrics struct {
	OS       string `json:"os"`
	CPUs     int    `json:"cpus"`
	AllocRAM uint64 `json:"alloc_ram"`
}
