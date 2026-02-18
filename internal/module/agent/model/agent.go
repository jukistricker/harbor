package model

import (
	"sync"
	"github.com/gorilla/websocket"
)

// ActiveAgent lưu trữ kết nối thực tế
type ActiveAgent struct {
	ID   string
	Conn *websocket.Conn
}

// AgentHub quản lý danh sách các Agent đang online
type AgentHub struct {
	sync.RWMutex
	Agents map[string]*ActiveAgent
}

// Khởi tạo Hub toàn cục cho module
var Hub = &AgentHub{
	Agents: make(map[string]*ActiveAgent),
}