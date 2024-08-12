package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	relaxingChannel := NewChannel("relaxing", []AudioContent{
		NewAudioContent("Relaxing 1", "audios/relaxing-1.mp3", 22*time.Second, []string{"relaxing", "chilling", "focusing"}),
		NewAudioContent("Relaxing 2", "audios/relaxing-2.mp3", 15*time.Second, []string{"relaxing", "sleeping", "focusing"}),
	})
	go relaxingChannel.Broadcast()

	drivingChannel := NewChannel("driving", []AudioContent{
		NewAudioContent("Driving 1", "audios/driving-1.mp3", 38*time.Second, []string{"relaxing", "driving"}),
	})
	go drivingChannel.Broadcast()

	http.HandleFunc("/relaxing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "audio/mpeg")
		w.Header().Add("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("Could not create flusher")
		}

		connection := NewConnection()
		relaxingChannel.AddConnection(connection)
		log.Printf("%s has connected to the audio stream\n", r.Host)

		log.Printf("total connection %d \n", relaxingChannel.GetNumOfConn())

		for {
			buf := <-connection.bufferChannel
			if _, err := w.Write(buf); err != nil {
				relaxingChannel.DeleteConnection(connection)
				log.Printf("%s's connection to the audio stream has been closed\n", r.Host)
				return
			}
			flusher.Flush()
			clear(connection.buffer)
		}
	})

	http.HandleFunc("/driving", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "audio/mpeg")
		w.Header().Add("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("Could not create flusher")
		}

		connection := NewConnection()
		drivingChannel.AddConnection(connection)
		log.Printf("%s has connected to the audio stream\n", r.Host)

		log.Printf("total connection %d \n", relaxingChannel.GetNumOfConn())

		for {
			buf := <-connection.bufferChannel
			if _, err := w.Write(buf); err != nil {
				drivingChannel.DeleteConnection(connection)
				log.Printf("%s's connection to the audio stream has been closed\n", r.Host)
				return
			}
			flusher.Flush()
			clear(connection.buffer)
		}
	})

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
