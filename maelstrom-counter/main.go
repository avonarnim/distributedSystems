package main

import (
	"context"
    "encoding/json"
    "log"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var ctx = context.TODO()


func main () {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)

	n.Handle("add", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// log.Println("in broadcasting")

		// Update the message type to return back.
		curVal, _ := kv.ReadInt(ctx, "counter")
		kv.Write(ctx, "counter", curVal+1)
		body["type"] = "add_ok"
		delete(body, "delta")
	
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// log.Println("in read")

		// Update the message type to return back.
		body["type"] = "read_ok"
		body["value"], _ = kv.ReadInt(ctx, "counter")
	
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})
	
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

