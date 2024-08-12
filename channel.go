package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"
)

type Channel struct {
	name                string
	bufferSize          int
	delayTime           int
	files               []AudioContent
	connPool            *ConnectionPool
	numOfConn           int
	startedAt           time.Time
	currentAudioContent *AudioContent
}

func NewChannel(name string, files []AudioContent) Channel {
	return Channel{
		name:                name,
		bufferSize:          4096,
		delayTime:           100,
		files:               files,
		connPool:            NewConnectionPool(),
		numOfConn:           0,
		startedAt:           time.Now(),
		currentAudioContent: nil,
	}
}

func (c *Channel) AddConnection(connection *Connection) {
	c.connPool.AddConnection(connection)
	c.numOfConn++
}

func (c *Channel) DeleteConnection(connection *Connection) {
	c.connPool.DeleteConnection(connection)
	c.numOfConn--
}

func (c *Channel) Broadcast() {
	buffer := make([]byte, c.bufferSize)

	for {
		for _, f := range c.files {
			file, err := os.Open(f.GetPath())
			if err != nil {
				log.Printf("Error opening file %s: %v", f.GetPath(), err)
				panic(err)
			}

			content, err := io.ReadAll(file)
			if err != nil {
				log.Printf("Error reading file %s: %v", f.GetName(), err)
				panic(err)
			}

			tempFile := bytes.NewReader(content)
			ticker := time.NewTicker(time.Millisecond * time.Duration(c.delayTime))

			log.Printf("Broadcasting %s, duration: %.0f seconds ...\n", f.GetName(), f.GetDuration().Seconds())

			for range ticker.C {
				clear(buffer)

				_, err = tempFile.Read(buffer)
				if err == io.EOF {
					log.Printf("%s ended, move to next file\n", f.GetName())

					ticker.Stop()

					time.Sleep(1 * time.Second)
					break
				}

				c.connPool.Broadcast(buffer)
			}
			func() { _ = file.Close() }()
		}
	}
}

func (c *Channel) GetNumOfConn() int {
	return c.numOfConn
}
