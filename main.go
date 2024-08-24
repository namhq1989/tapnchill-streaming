package main

import (
	"fmt"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func broadcastingChannel(name string) Channel {
	audios, err := ListMP3Files("audios/" + name)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(audios); i++ {
		fmt.Printf("%d. %s duration: %.0f seconds \n", i+1, audios[i].GetName(), audios[i].GetDuration().Seconds())
	}

	channel := NewChannel("name", audios)
	go channel.Broadcast()

	return channel
}

func broadcastingRoute(mux *http.ServeMux, id string, channel Channel) {
	mux.HandleFunc("/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "audio/mpeg")
		w.Header().Add("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("Could not create flusher")
		}

		connection := NewConnection()
		channel.AddConnection(connection)
		log.Printf("%s has connected to the audio stream %s\n, total connection %d", r.Host, id, channel.GetNumOfConn())

		for {
			buf := <-connection.bufferChannel
			if _, err := w.Write(buf); err != nil {
				channel.DeleteConnection(connection)
				log.Printf("%s's connection to the audio stream %s has been closed, total connection %d\n", r.Host, id, channel.GetNumOfConn())
				return
			}
			flusher.Flush()
			clear(connection.buffer)
		}
	})
}

func main() {
	mux := http.NewServeMux()

	chillingChannel := broadcastingChannel("chilling")
	broadcastingRoute(mux, "chilling", chillingChannel)

	gamingChannel := broadcastingChannel("gaming")
	broadcastingRoute(mux, "gaming", gamingChannel)

	motivatingChannel := broadcastingChannel("motivating")
	broadcastingRoute(mux, "motivating", motivatingChannel)

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}
