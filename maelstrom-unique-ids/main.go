package main

import (
	"encoding/json"
	_ "errors"
	"log"
	"strings"

	"github.com/pkg/errors"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	length   = 12
)

// New Generates a unique ID.
func New() (string, error) {
	return nanoid.Generate(alphabet, length)
}

// Must is the same as new, but panic on error
func Must() string {
	return nanoid.MustGenerate(alphabet, length)
}

// Validate checks if a given field name's public ID value is valid
// according to the contraints defined by the application.
func Validate(id string) error {
	if id == "" {
		return errors.Errorf("id cannot be blank")
	}
	if len(id) != length {
		return errors.Errorf("id should be %d characters long", length)
	}
	if strings.Trim(id, alphabet) != "" {
		return errors.Errorf("id has invalid characters")
	}
	return nil
}

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		generated_id := Must()
		body["type"] = "generate_ok"
		body["id"] = generated_id

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
