package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	connections map[int]*websocket.Conn
	mu          sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[int]*websocket.Conn),
	}
}

func (m *Manager) Add(userID int, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connections[userID] = conn
}

func (m *Manager) Remove(userID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.connections, userID)
}

func (m *Manager) Get(userID int) (*websocket.Conn, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connections[userID]

	return conn, exists
}
