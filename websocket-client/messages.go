package websocketclient

import "encoding/json"

// Event is a struct that represents a message from the websocket server
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event)

var eventHandlers = make(map[string]EventHandler)

// RegisterEventHandler registers an event handler for a specific event type, customizing the behavior of the clients
func RegisterEventHandler(eventType string, handler EventHandler) {
	eventHandlers[eventType] = handler
}
