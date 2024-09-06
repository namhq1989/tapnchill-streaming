package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type Channel struct {
	id                  string
	bufferSize          int
	delayTime           int
	files               []AudioContent
	connPool            *ConnectionPool
	numOfConn           int
	startedAt           time.Time
	currentAudioContent *AudioContent
}

func NewChannel(id string) Channel {
	channel := Channel{
		id:                  id,
		bufferSize:          8192,
		delayTime:           150,
		connPool:            NewConnectionPool(),
		numOfConn:           0,
		startedAt:           time.Now(),
		currentAudioContent: nil,
	}

	err := channel.updateFileList()
	if err != nil {
		panic(err)
	}

	return channel

}

func (c *Channel) AddConnection(connection *Connection) {
	c.connPool.AddConnection(connection)
	c.numOfConn++
}

func (c *Channel) DeleteConnection(connection *Connection) {
	c.connPool.DeleteConnection(connection)
	c.numOfConn--
}

func (c *Channel) updateFileList() error {
	audios, err := ListMP3Files("audios/" + c.id)
	if err != nil {
		return err
	}
	c.files = audios
	return nil
}

func (c *Channel) calculateDelayTime(fileInfo os.FileInfo, duration time.Duration, bufferSize int) {
	fileSize := fileInfo.Size()
	totalBuffers := fileSize / int64(bufferSize)
	c.delayTime = int(duration / time.Duration(totalBuffers) / time.Millisecond)
}

func (c *Channel) shuffleFiles(files []AudioContent) []AudioContent {
	rand.New(rand.NewSource(time.Now().UnixMilli()))
	rand.Shuffle(len(files), func(i, j int) {
		files[i], files[j] = files[j], files[i]
	})
	return files
}

func (c *Channel) Broadcast() {
	var lastPlayed *AudioContent

	for {
		updateErr := c.updateFileList()
		if updateErr != nil {
			log.Printf("Error updating file list: %v", updateErr)
			panic(updateErr)
		}

		if len(c.files) == 0 {
			fmt.Printf("channel %s is empty, restarting ...\n", c.id)
			time.Sleep(1 * time.Minute)
			continue
		}

		shuffledFiles := c.shuffleFiles(c.files)

		if lastPlayed != nil && shuffledFiles[0].GetPath() == lastPlayed.GetPath() {
			if len(shuffledFiles) > 1 {
				shuffledFiles[0], shuffledFiles[1] = shuffledFiles[1], shuffledFiles[0]
			}
		}

		for _, f := range shuffledFiles {
			file, err := os.Open(f.GetPath())
			if err != nil {
				log.Printf("Error opening file %s: %v", f.GetPath(), err)
				panic(err)
			}

			fileInfo, err := file.Stat()
			if err != nil {
				log.Printf("Error getting file info for %s: %v", f.GetName(), err)
				panic(err)
			}

			c.calculateDelayTime(fileInfo, f.GetDuration(), c.bufferSize)
			buffer := make([]byte, c.bufferSize)

			content, err := io.ReadAll(file)
			if err != nil {
				log.Printf("Error reading file %s: %v", f.GetName(), err)
				panic(err)
			}

			tempFile := bytes.NewReader(content)
			ticker := time.NewTicker(time.Millisecond * time.Duration(c.delayTime))

			log.Printf("Broadcasting %s, duration: %.0f seconds, delay time: %d ...\n", f.GetName(), f.GetDuration().Seconds(), c.delayTime)

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

			lastPlayed = &f
		}
	}
}

func (c *Channel) GetID() string { return c.id }

func (c *Channel) GetNumOfConn() int {
	return c.numOfConn
}
