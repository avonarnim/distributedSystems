package main

import (
    "encoding/json"
    "log"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var broadcastMessages = []int{}
// var mu sync.Mutex

func main () {
	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// Update the message type to return back.
		broadcastMessages = append(broadcastMessages, int(body["message"].(float64)))

		for _, node := range n.NodeIDs() {
			if node != n.ID() {
				n.Send(node, body)
			}
		}

		body = map[string]any{}
		body["type"] = "broadcast_ok"

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
		body["messages"] = broadcastMessages
	
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// log.Println("in topology")

		// Update the message type to return back.
		body["type"] = "topology_ok"
		delete(body, "topology")
		
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})
	
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

