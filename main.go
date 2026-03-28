package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var version = "dev"

func main() {
	port    := flag.Int("port", 8765, "port for local server")
	watch   := flag.String("watch", "", "JSON file to watch for changes")
	noOpen  := flag.Bool("no-open", false, "do not auto-open browser")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	hub := newHub()
	go hub.run()

	if *watch != "" {
		go pollFile(*watch, hub)
	} else {
		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to stat stdin: %v\n", err)
			os.Exit(1)
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "error: pipe JSON data via stdin or use --watch <file>")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "usage:")
			fmt.Fprintln(os.Stderr, "  cat data.json | jviz")
			fmt.Fprintln(os.Stderr, "  jstats \"count by status\" < data.json | jviz")
			fmt.Fprintln(os.Stderr, "  jviz --watch data.json")
			os.Exit(1)
		}
		go readStdin(hub)
	}

	addr := fmt.Sprintf(":%d", *port)
	url  := fmt.Sprintf("http://localhost:%d", *port)

	if !*noOpen {
		go func() {
			time.Sleep(300 * time.Millisecond)
			openBrowser(url)
		}()
	}

	fmt.Fprintf(os.Stderr, "jviz: serving at %s  (Ctrl+C to stop)\n", url)
	if err := serve(addr, hub); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd, args = "open", []string{url}
	case "windows":
		cmd, args = "cmd", []string{"/c", "start", url}
	default:
		cmd, args = "xdg-open", []string{url}
	}
	_ = exec.Command(cmd, args...).Start()
}

// readStdin decodes a stream of JSON arrays from stdin.
func readStdin(hub *Hub) {
	decodeJSONStream(os.Stdin, hub.send)
}

// decodeJSONStream reads JSON arrays from r and calls send for each.
// This handles both single-output (cat data.json | jviz) and
// streaming (while true; do ...; done | jviz) cases.
func decodeJSONStream(r io.Reader, send func([]map[string]any)) {
	dec := json.NewDecoder(r)
	for {
		var rows []map[string]any
		err := dec.Decode(&rows)
		if err == io.EOF {
			return
		}
		if err != nil {
			// Skip to the next newline to avoid looping on the same bad input,
			// then rebuild the decoder from the remaining buffered bytes + r.
			buf := dec.Buffered()
			b := make([]byte, 1)
			for {
				if _, rerr := buf.Read(b); rerr != nil || b[0] == '\n' {
					break
				}
			}
			dec = json.NewDecoder(io.MultiReader(buf, r))
			continue
		}
		send(rows)
	}
}

// pollFile watches path for content changes every second.
func pollFile(path string, hub *Hub) {
	var lastModTime int64
	for {
		info, err := os.Stat(path)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if info.ModTime().UnixNano() == lastModTime {
			time.Sleep(time.Second)
			continue
		}
		lastModTime = info.ModTime().UnixNano()

		data, err := os.ReadFile(path)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		var rows []map[string]any
		if err := json.Unmarshal(data, &rows); err == nil {
			hub.send(rows)
		}
		time.Sleep(time.Second)
	}
}
