package main

import "sync"

type Connection struct {
	bufferChannel chan []byte
	buffer        []byte
}

type ConnectionPool struct {
	ConnectionMap map[*Connection]struct{}
	mu            sync.Mutex
}

const BUFFERSIZE = 8192

func NewConnection() *Connection {
	connection := &Connection{bufferChannel: make(chan []byte), buffer: make([]byte, BUFFERSIZE)}
	return connection
}

func (cp *ConnectionPool) AddConnection(connection *Connection) {
	defer cp.mu.Unlock()
	cp.mu.Lock()

	cp.ConnectionMap[connection] = struct{}{}

}

func (cp *ConnectionPool) DeleteConnection(connection *Connection) {
	defer cp.mu.Unlock()
	cp.mu.Lock()

	delete(cp.ConnectionMap, connection)

}

func (cp *ConnectionPool) Broadcast(buffer []byte) {
	defer cp.mu.Unlock()
	cp.mu.Lock()

	for connection := range cp.ConnectionMap {
		copy(connection.buffer, buffer)
		select {
		case connection.bufferChannel <- connection.buffer:
		default:
		}
	}
}

func NewConnectionPool() *ConnectionPool {
	connectionMap := make(map[*Connection]struct{})
	return &ConnectionPool{ConnectionMap: connectionMap}
}
