package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/go-mp3"
)

func ListMP3Files(id string) ([]AudioContent, error) {
	var (
		contents []AudioContent
		err      error
	)

	contents, err = walkDir(id)
	if err != nil {
		return nil, err
	}
	return contents, nil

	// if id != "mix" {
	// 	contents, err = walkDir(id)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return contents, nil
	// }
	//
	// for _, channel := range channelIDs {
	// 	if channel == "mix" {
	// 		continue
	// 	}
	//
	// 	channelContents, walkDirErr := walkDir(channel)
	// 	if walkDirErr != nil {
	// 		return nil, walkDirErr
	// 	}
	// 	contents = append(contents, channelContents...)
	// }

	// return contents, nil
}

func walkDir(id string) ([]AudioContent, error) {
	var contents []AudioContent

	err := filepath.Walk("audios/"+id, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(info.Name()) == ".mp3" {
			file, openErr := os.Open(path)
			if openErr != nil {
				return openErr
			}
			defer func() { _ = file.Close() }()

			decoder, decodeErr := mp3.NewDecoder(file)
			if decodeErr != nil {
				return decodeErr
			}

			// Calculate the duration
			bytesPerSecond := float64(decoder.SampleRate() * 4) // 2 bytes per sample * 2 channels
			durationSeconds := float64(decoder.Length()) / bytesPerSecond
			duration := time.Duration(durationSeconds) * time.Second

			// Create an AudioContent struct with the file details
			content := NewAudioContent(info.Name(), path, duration)

			// Append the struct to the slice
			contents = append(contents, content)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return contents, nil
}
