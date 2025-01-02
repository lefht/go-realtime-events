package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	flag.Parse()

	http.HandleFunc("/events", events)
	http.ListenAndServe(":8080", withCORS(http.HandlerFunc(events)))
}

func events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")

	// run ping command and get the output
	cmd := exec.Command("ping", "-c 4", "google.com")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to get command output", http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start command", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		content := fmt.Sprintf("data: %s\n\n", scanner.Text())
		w.Write([]byte(content))
		w.(http.Flusher).Flush()
		time.Sleep(time.Millisecond * 400)
	}

	if err := cmd.Wait(); err != nil {
		http.Error(w, "Command execution failed", http.StatusInternalServerError)
		return
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}
