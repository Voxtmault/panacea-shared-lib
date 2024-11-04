package websocketclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/voxtmault/panacea-shared-lib/config"
	"github.com/voxtmault/panacea-shared-lib/websocket-client/types"

	"github.com/gorilla/websocket"
	"github.com/rotisserie/eris"
)

var conn *websocket.Conn

func connectWebSocket(serverURL string, headers http.Header) (*websocket.Conn, error) {
	// Connect to the WebSocket server with custom headers
	connection, _, err := websocket.DefaultDialer.Dial(serverURL, headers)
	if err != nil {
		return nil, fmt.Errorf("error connecting to WebSocket server: %w", err)
	}
	return connection, nil
}

// listenForMessages listens for incoming message from websocket hub. Use DEBUG=true to print the message.
func listenForMessages() {

	for {
		// Read message from the server
		_, message, err := conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "connection reset by peer") ||
				strings.Contains(err.Error(), strconv.Itoa(websocket.CloseGoingAway)) ||
				strings.Contains(err.Error(), strconv.Itoa(websocket.CloseAbnormalClosure)) {
				log.Println("Websocket connection closed abnormally, attempting to reconnect", err)

				// Close the current connection
				if err := conn.Close(); err != nil {
					log.Println("Error closing websocket connection:", err)
				}

				// Attempt to reconnect
				for {
					conn, err = connectWebSocket(config.GetConfig().WebsocketConfig.WSURL, http.Header{
						"X-API-TOKEN": []string{config.GetConfig().WebsocketConfig.WSApiToken},
					})
					if err != nil {
						log.Println("Reconnect attempt failed:", err)
						time.Sleep(time.Second * time.Duration(config.GetConfig().WebsocketConfig.WSReconnectInterval))
					} else {
						log.Println("Reconnected to websocket server")
						break
					}
				}

			} else {
				log.Println("Error reading message:", err)
				return
			}
		}
		// Print the received message
		if config.GetConfig().DebugMode {
			log.Println("Received message:", string(message))
		}

		// Business Logic Here
		websocketBusinessLogic(message)
	}
}

func InitWebsocketClient() error {
	if config.GetConfig().WebsocketConfig.WSURL == "" {
		return eris.New("Websocket URL not set")
	}
	if config.GetConfig().WebsocketConfig.WSApiToken == "" {
		return eris.New("Websocket API Token not set")
	}

	// Authenticate as an API
	headers := http.Header{
		"X-API-TOKEN": []string{config.GetConfig().WebsocketConfig.WSApiToken},
	}

	// Establish WebSocket connection
	var err error
	conn, err = connectWebSocket(config.GetConfig().WebsocketConfig.WSURL, headers)
	if err != nil {
		return err
	}

	// Start a goroutine to listen for messages from the WebSocket server
	go listenForMessages()

	slog.Info("Successfully Connected To Websocket Server")
	return nil
}

func CloseWebsocketClient() error {
	slog.Info("Closing Websocket Connection")
	// Ensure the WebSocket connection is closed
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func GetWSConn() *websocket.Conn {
	return conn
}

// DO NOT MARSHALL WHATEVER MESSAGE (3rd Parameter) YOU MIGHT HAVE, THE FUNCTION IS ALREADY GOING TO MARSHALL IT
//
// DO NOT COMPLAIN TO ME IF SHIT'S NOT GETTING HANDLED CORRECTLY
func SendMessage(ctx context.Context, messageType types.EventList, message interface{}) error {

	var msg Event
	var err error

	msg.Type = messageType
	msg.Payload, err = json.Marshal(message)
	if err != nil {
		return eris.Wrap(err, "Marshalling Payload")
	}

	// log.Println("Sending message:", string(msg.Payload))

	conn.WriteJSON(msg)

	return nil
}

func websocketBusinessLogic(event []byte) {

	var message Event

	if event == nil {
		// Usually from reconnect message
		return
	}

	// Unmarshall to get the event
	if err := json.Unmarshal(event, &message); err != nil {
		slog.Error("Websocket Business Logic", "Unmarshalling Message", err)
		return
	}

	// Check if message type is supported, if it is then call the appropriate function in model to handle this
	if handler, exists := eventHandlers[message.Type]; exists {
		handler(message)
	} else {
		slog.Info("Websocket Business Logic", "Unsupported Message Type", message.Type)
	}
}
