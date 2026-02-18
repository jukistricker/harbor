package model

type CalculatedMetric struct {
	TargetID          string  `json:"target_id"`
	Timestamp         uint64  `json:"timestamp"`
	CPUUsagePercent   float64 `json:"cpu_usage_percent"`
	MemPrivateBytes   uint64  `json:"mem_private_bytes"`
	MemGrowthBytes    int64   `json:"mem_growth_bytes"`
	HandleCount       uint32  `json:"handle_count"`
	ThreadCount       uint32  `json:"thread_count"`
	ActiveConnections uint32  `json:"active_connections"`
	NetOutErrors      uint64  `json:"net_out_errors"`
	Mode              string  `json:"mode"`
}