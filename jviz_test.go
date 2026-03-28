package main

import (
	"strings"
	"testing"
	"time"
)

// --- sseMessage ---

func TestSSEMessage_Format(t *testing.T) {
	rows := []map[string]any{{"name": "Alice", "age": 30}}
	msg := sseMessage(rows)
	s := string(msg)
	if !strings.HasPrefix(s, "data: ") {
		t.Errorf("sseMessage should start with 'data: ', got: %q", s)
	}
	if !strings.HasSuffix(s, "\n\n") {
		t.Errorf("sseMessage should end with '\\n\\n', got: %q", s)
	}
}

func TestSSEMessage_EmptySlice(t *testing.T) {
	msg := sseMessage([]map[string]any{})
	if string(msg) != "data: []\n\n" {
		t.Errorf("unexpected output for empty slice: %q", string(msg))
	}
}

// --- Hub ---

func TestHub_SnapshotUpdatesAfterSend(t *testing.T) {
	h := newHub()
	go h.run()

	rows := []map[string]any{{"k": "v"}}
	h.send(rows)

	// Wait for run() to process
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if s := h.snapshot(); len(s) > 0 {
			if s[0]["k"] == "v" {
				return
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Error("snapshot was not updated after send")
}

func TestHub_SnapshotNilBeforeSend(t *testing.T) {
	h := newHub()
	if h.snapshot() != nil {
		t.Error("snapshot should be nil before any send")
	}
}

func TestHub_SubscribeReceivesMessage(t *testing.T) {
	h := newHub()
	go h.run()

	ch := h.subscribe()
	defer h.unsubscribe(ch)

	rows := []map[string]any{{"x": 1}}
	h.send(rows)

	select {
	case msg := <-ch:
		s := string(msg)
		if !strings.Contains(s, `"x"`) {
			t.Errorf("received message missing expected field: %q", s)
		}
	case <-time.After(time.Second):
		t.Error("timed out waiting for message on subscribed channel")
	}
}

func TestHub_UnsubscribeStopsDelivery(t *testing.T) {
	h := newHub()
	go h.run()

	ch := h.subscribe()
	h.unsubscribe(ch)

	h.send([]map[string]any{{"x": 1}})

	select {
	case <-ch:
		t.Error("should not receive message after unsubscribe")
	case <-time.After(100 * time.Millisecond):
		// expected: no message delivered
	}
}

func TestHub_MultipleSubscribers(t *testing.T) {
	h := newHub()
	go h.run()

	ch1 := h.subscribe()
	ch2 := h.subscribe()
	defer h.unsubscribe(ch1)
	defer h.unsubscribe(ch2)

	h.send([]map[string]any{{"n": 42}})

	for _, ch := range []chan []byte{ch1, ch2} {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Error("timed out waiting for message on subscriber channel")
		}
	}
}

// --- decodeJSONStream ---

func TestDecodeJSONStream_Single(t *testing.T) {
	input := `[{"a":"1"},{"a":"2"}]`
	var got [][]map[string]any
	decodeJSONStream(strings.NewReader(input), func(rows []map[string]any) {
		got = append(got, rows)
	})
	if len(got) != 1 || len(got[0]) != 2 {
		t.Errorf("expected 1 batch of 2 rows, got %d batches", len(got))
	}
}

func TestDecodeJSONStream_Multiple(t *testing.T) {
	input := `[{"n":1}]` + "\n" + `[{"n":2}]` + "\n" + `[{"n":3}]`
	var got [][]map[string]any
	decodeJSONStream(strings.NewReader(input), func(rows []map[string]any) {
		got = append(got, rows)
	})
	if len(got) != 3 {
		t.Errorf("expected 3 batches, got %d", len(got))
	}
}

func TestDecodeJSONStream_SkipsInvalidTokens(t *testing.T) {
	input := `garbage` + "\n" + `[{"ok":true}]`
	var got [][]map[string]any
	decodeJSONStream(strings.NewReader(input), func(rows []map[string]any) {
		got = append(got, rows)
	})
	if len(got) != 1 {
		t.Errorf("expected 1 valid batch after garbage, got %d", len(got))
	}
}

func TestDecodeJSONStream_Empty(t *testing.T) {
	var called bool
	decodeJSONStream(strings.NewReader(""), func(_ []map[string]any) {
		called = true
	})
	if called {
		t.Error("send should not be called for empty input")
	}
}
