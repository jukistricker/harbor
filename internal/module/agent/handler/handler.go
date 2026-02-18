package handler

import (
	"encoding/json"
	"fmt"
	"harbor/internal/module/agent/model"
	"harbor/internal/module/metric"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler xử lý kết nối từ Agent (Rust)
func WSHandler(metricMod *metric.Module) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("❌ Upgrade error: %v", err)
			return
		}

		agentID := r.URL.Query().Get("id")
		if agentID == "" {
			agentID = r.RemoteAddr
		}

		// Đăng ký vào Hub
		model.Hub.Lock()
		model.Hub.Agents[agentID] = &model.ActiveAgent{ID: agentID, Conn: conn}
		model.Hub.Unlock()

		log.Printf("📥 Agent [%s] connected", agentID)

		defer func() {
			model.Hub.Lock()
			delete(model.Hub.Agents, agentID)
			model.Hub.Unlock()
			conn.Close()
			log.Printf("🔌 Agent [%s] disconnected", agentID)
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			metricMod.Service.ProcessRawMetric(agentID, msg)
			// Chuyển dữ liệu sang module metric để xử lý lưu DB
			log.Printf("📊 Data from [%s]: %s", agentID, string(msg))
			// metricMod.Service.ProcessRawMetric(agentID, msg) 
		}
	}
}

// ControlHandler nhận lệnh từ UI/Postman để điều khiển Agent cụ thể
func ControlHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AgentID string `json:"agent_id"`
		Action  string `json:"action"` // START | STOP
		Port    int    `json:"port"`
		Target  string `json:"target"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	model.Hub.RLock()
	agent, exists := model.Hub.Agents[req.AgentID]
	model.Hub.RUnlock()

	if !exists {
		http.Error(w, "Agent not found", http.StatusNotFound)
		return
	}

	// Forward lệnh xuống Agent qua WebSocket
	cmd, _ := json.Marshal(req)
	err := agent.Conn.WriteMessage(websocket.TextMessage, cmd)
	if err != nil {
		http.Error(w, "Failed to send command", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Command %s sent to %s", req.Action, req.AgentID)
}