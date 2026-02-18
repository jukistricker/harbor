package model

type AgentPayload struct {
    Metadata struct {
		AgentID     string `json:"agent_id"`
        Timestamp    uint64 `json:"timestamp"`
        TargetPID    uint32 `json:"target_pid"`
        ProcessPath  string `json:"process_path"`
    } `json:"metadata"`
    Layer2 struct {
        HandleCount  uint32 `json:"handle_count"`
        ThreadCount  uint32 `json:"thread_count"`
        MemPrivate   uint64 `json:"memory_private_bytes"`
        TotalCycles  uint64 `json:"total_cycles"`
    } `json:"layer2_raw"`
    Layer3 struct {
        ActiveConns  uint32 `json:"active_connections"`
    } `json:"layer3_raw"`
    Layer4 struct {
        DiskRead     uint64 `json:"disk_read_bytes"`
        NetOutErrors uint64 `json:"net_out_errors"`
    } `json:"layer4_raw"`
    Mode string `json:"mode"`
}

