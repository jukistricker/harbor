package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"harbor/internal/module/metric/model"
	"log"
	"path/filepath"
	"sync"
)

type Processor struct {
	db         *sql.DB
	lastStates map[uint32]model.AgentPayload
	buffer     []model.CalculatedMetric
	mu         sync.Mutex
}

func NewProcessor(db *sql.DB) *Processor {
	return &Processor{
		db:         db,
		lastStates: make(map[uint32]model.AgentPayload),
		buffer:     make([]model.CalculatedMetric, 0, 100),
	}
}

func (p *Processor) Process(current model.AgentPayload) {
	p.mu.Lock()
	defer p.mu.Unlock()

	prev, exists := p.lastStates[current.Metadata.TargetPID]
	p.lastStates[current.Metadata.TargetPID] = current

	if !exists {
		return 
	}

	deltaTime := float64(current.Metadata.Timestamp - prev.Metadata.Timestamp)
	deltaCycles := float64(current.Layer2.TotalCycles - prev.Layer2.TotalCycles)
	
	cpuPercent := 0.0
	if deltaTime > 0 {
		// Windows FileTime unit (100ns) -> 1s = 10^7 units
		cpuPercent = (deltaCycles / (deltaTime * 10000000.0)) * 100.0
	}

	memGrowth := int64(current.Layer2.MemPrivate) - int64(prev.Layer2.MemPrivate)

	calc := model.CalculatedMetric{
		TargetID:          fmt.Sprintf("%s_%d", filepath.Base(current.Metadata.ProcessPath), current.Metadata.TargetPID),
		Timestamp:         current.Metadata.Timestamp,
		CPUUsagePercent:   cpuPercent,
		MemPrivateBytes:   current.Layer2.MemPrivate,
		MemGrowthBytes:    memGrowth,
		HandleCount:       current.Layer2.HandleCount,
		ThreadCount:       current.Layer2.ThreadCount,
		ActiveConnections: current.Layer3.ActiveConns,
		NetOutErrors:      current.Layer4.NetOutErrors,
		Mode:              current.Mode,
	}
	p.buffer = append(p.buffer, calc)

	if len(p.buffer) >= 50 {
		p.Flush()
	}
}

func (p *Processor) ProcessRawMetric(agentID string, rawData []byte) {
    var payload model.AgentPayload
    if err := json.Unmarshal(rawData, &payload); err != nil {
        log.Printf("⚠️ [Metric] Failed to unmarshal from %s: %v", agentID, err)
        return
    }

    // Gán AgentID vào Metadata nếu bạn muốn phân biệt các máy
    payload.Metadata.AgentID = agentID 

    // Gọi hàm Process để tính toán Delta và Buffer
    p.Process(payload)
}

func (p *Processor) Flush() {
	p.mu.Lock()
    count := len(p.buffer) // Lấy số lượng trước khi clear
    if count == 0 {
        p.mu.Unlock()
        return
    }
    p.mu.Unlock()

	// Bắt đầu Transaction để tối ưu Batch Insert cho SQLite
	tx, err := p.db.Begin()
	if err != nil {
		log.Printf("[DB] Failed to begin transaction: %v", err)
		return
	}

	stmt, err := tx.Prepare(`
		INSERT INTO harbor_metrics (
			target_id, timestamp, handle_count, thread_count, 
			mem_private, total_cycles, active_connections, 
			disk_read, disk_write, net_out_errors, mode
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Printf("[DB] Failed to prepare statement: %v", err)
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, m := range p.buffer {
		// Lưu ý: Lưu total_cycles thực tế (Raw) 
		// Nếu muốn lưu CPU %, sửa schema init.sql thêm cột
		_, err := stmt.Exec(
			m.TargetID, m.Timestamp, m.HandleCount, m.ThreadCount,
			m.MemPrivateBytes, uint64(m.CPUUsagePercent*100), // Trick: Lưu CPU*100 vào total_cycles nếu chưa kịp sửa schema
			m.ActiveConnections, 0, 0, m.NetOutErrors, m.Mode,
		)
		if err != nil {
			log.Printf("[DB] Insert error: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[DB] Transaction commit error: %v", err)
	}

	// Clear buffer nhưng giữ capacity
	p.mu.Lock()
    p.buffer = p.buffer[:0]
    p.mu.Unlock()
    log.Printf("💾 [Processor] Flushed %d metrics to DB", count)
}