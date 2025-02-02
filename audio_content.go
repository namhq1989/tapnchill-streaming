package main

import "time"

type AudioContent struct {
	name     string
	path     string
	duration time.Duration
}

func NewAudioContent(name, path string, duration time.Duration) AudioContent {
	return AudioContent{name: name, path: path, duration: duration}
}

func (c *AudioContent) GetDuration() time.Duration {
	return c.duration
}

func (c *AudioContent) GetName() string {
	return c.name
}

func (c *AudioContent) GetPath() string {
	return c.path
}
