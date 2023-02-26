package main

import (
    "encoding/json"
    "log"
	// "sync"
	// "sync/atomic"
	"github.com/google/uuid"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var idCounter uint64
// var mu sync.Mutex

func main () {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// Update the message type to return back.
		body["type"] = "generate_ok"
		id := uuid.New().String()
		body["id"] = id
		// mu.Lock()
		// atomic.AddUint64(&idCounter, 1)
		// body["id"] = atomic.LoadUint64(&idCounter)
		// mu.Unlock()
	
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})
	
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

