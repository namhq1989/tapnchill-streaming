package main

import "time"

type AudioContent struct {
	name     string
	path     string
	duration time.Duration
	topics   []string
}

func NewAudioContent(name, path string, duration time.Duration, topics []string) AudioContent {
	return AudioContent{name: name, path: path, duration: duration, topics: topics}
}

func (c *AudioContent) GetDuration() time.Duration {
	return c.duration
}

func (c *AudioContent) GetTopics() []string {
	return c.topics
}

func (c *AudioContent) GetName() string {
	return c.name
}

func (c *AudioContent) GetPath() string {
	return c.path
}
