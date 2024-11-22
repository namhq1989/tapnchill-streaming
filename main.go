package main

import (
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

func broadcastingChannel(id string) Channel {
	channel := NewChannel(id)
	go channel.Broadcast()

	return channel
}

func broadcastingRoute(mux *http.ServeMux, id string, channel *Channel) {
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

var channels = make([]*Channel, 0)
var channelIDs = []string{"mix"}

func main() {
	mux := http.NewServeMux()

	for _, id := range channelIDs {
		// if id != "mix" {
		// 	if err := os.MkdirAll("audios/"+id, os.ModePerm); err != nil {
		// 		panic(fmt.Errorf("failed to create directory: %v", err))
		// 	}
		// }

		channel := broadcastingChannel(id)
		broadcastingRoute(mux, id, &channel)
		channels = append(channels, &channel)
	}

	// go startSendingChannelStats()

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}

// func startSendingChannelStats() {
// 	ticker := time.NewTicker(1 * time.Minute)
// 	defer ticker.Stop()
//
// 	for range ticker.C {
// 		sendChannelStats()
// 	}
// }

// var tapnchillServer = os.Getenv("TAPNCHILL_SERVER")
//
// func sendChannelStats() {
// 	url := fmt.Sprintf("%s/api/webhook/channel-stats", tapnchillServer)
//
// 	stats := make([]map[string]interface{}, len(channels))
//
// 	for i, channel := range channels {
// 		stats[i] = map[string]interface{}{
// 			"id":        channel.GetID(),
// 			"audiences": channel.GetNumOfConn(),
// 		}
// 	}
//
// 	payload, err := json.Marshal(stats)
// 	if err != nil {
// 		log.Printf("Failed to marshal stats: %v\n", err)
// 		return
// 	}
//
// 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
// 	if err != nil {
// 		log.Printf("Failed to send channel stats: %v\n", err)
// 		return
// 	}
// 	defer func() { _ = resp.Body.Close() }()
//
// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("Received non-OK response: %v\n", resp.Status)
// 	}
// }
