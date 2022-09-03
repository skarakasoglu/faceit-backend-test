package health

// Response health check response model
type Response struct {
	Status             bool               `json:"status"`
	DatabaseConnection DatabaseConnection `json:"database_connection"`
	MemoryStats        MemoryStats        `json:"memory_stats"`
}

// DatabaseConnection health check db connection status model
type DatabaseConnection struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error"`
}

// MemoryStats is a model which contains the information about the memory usage of the app
type MemoryStats struct {
	CurrentAllocated                uint64 `json:"current_allocated"`
	TotalAllocated                  uint64 `json:"total_allocated"`
	TotalMemory                     uint64 `json:"total_memory"`
	CompletedGarbageCollectorCycles uint32 `json:"completed_garbage_collector_cycles"`
}
