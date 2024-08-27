package websocketclient

import (
	"encoding/json"

	"github.com/voxtmault/panacea-shared-lib/websocket-client/types"
)

// Event is a struct that represents a message from the websocket server
type Event struct {
	Response types.EventResponse `json:"response"`
	Type     types.EventList     `json:"type"`
	Payload  json.RawMessage     `json:"payload"`
}

type EventHandler func(event Event)

var eventHandlers = make(map[types.EventList]EventHandler)

// RegisterEventHandler registers an event handler for a specific event type, customizing the behavior of the clients
func RegisterEventHandler(eventType types.EventList, handler EventHandler) {
	eventHandlers[eventType] = handler
}
