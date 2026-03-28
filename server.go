package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sync"

	_ "embed"

	"embed"
)

//go:embed all:static
var staticFS embed.FS

// Hub manages SSE clients and the most recent data snapshot.
type Hub struct {
	mu      sync.RWMutex
	current []map[string]any
	clients map[chan []byte]struct{}
	in      chan []map[string]any
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[chan []byte]struct{}),
		in:      make(chan []map[string]any, 8),
	}
}

// run processes incoming data and fans it out to all SSE clients.
func (h *Hub) run() {
	for rows := range h.in {
		h.mu.Lock()
		h.current = rows
		h.mu.Unlock()

		msg := sseMessage(rows)

		h.mu.RLock()
		for ch := range h.clients {
			select {
			case ch <- msg:
			default: // drop if client is too slow
			}
		}
		h.mu.RUnlock()
	}
}

func (h *Hub) send(rows []map[string]any) { h.in <- rows }

func (h *Hub) subscribe() chan []byte {
	ch := make(chan []byte, 16)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) unsubscribe(ch chan []byte) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
}

func (h *Hub) snapshot() []map[string]any {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.current
}

// sseMessage encodes rows as an SSE data event.
func sseMessage(rows []map[string]any) []byte {
	b, _ := json.Marshal(rows)
	return fmt.Appendf(nil, "data: %s\n\n", b)
}

// serve starts the HTTP server.
func serve(addr string, hub *Hub) error {
	mux := http.NewServeMux()

	// Current data snapshot (used by UI on reconnect)
	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		rows := hub.snapshot()
		if rows == nil {
			rows = []map[string]any{}
		}
		_ = json.NewEncoder(w).Encode(rows)
	})

	// SSE stream
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := hub.subscribe()
		defer hub.unsubscribe(ch)

		// Send current snapshot immediately so the browser gets data on connect/reconnect
		if rows := hub.snapshot(); rows != nil {
			_, _ = w.Write(sseMessage(rows))
			flusher.Flush()
		}

		for {
			select {
			case msg := <-ch:
				_, _ = w.Write(msg)
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	// Serve static files (index.html, chart.js, …) — must be last
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	return http.ListenAndServe(addr, mux)
}
