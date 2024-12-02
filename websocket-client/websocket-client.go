package websocketclient

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/voxtmault/panacea-shared-lib/config"
	"github.com/voxtmault/panacea-shared-lib/websocket-client/types"

	"github.com/gorilla/websocket"
	"github.com/rotisserie/eris"
)

var (
	conn      *websocket.Conn
	connMutex sync.Mutex
	closing   bool
)

func connectWebSocket(serverURL string) error {
	var err error
	headers := http.Header{
		"X-API-TOKEN": []string{config.GetConfig().WebsocketConfig.WSApiToken},
	}

	// Connect to the WebSocket server with custom headers
	conn, _, err = websocket.DefaultDialer.Dial(serverURL, headers)
	if err != nil {
		slog.Error("error connecting to websocket server", "reason", err)
		return eris.Wrap(err, "error connecting to WebSocket server")
	}
	return nil
}

// listenForMessages listens for incoming message from websocket hub. Use DEBUG=true to print the message.
func listenForMessages() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Handle read message in a separate goroutine
	go func() {
		defer close(done)
		for {
			// Read message from the server
			_, message, err := conn.ReadMessage()
			if err != nil {
				// For every read error, we will attempt to reconnect to the server while also checking
				// if the connection is intented to be closed or not

				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					slog.Info("websocket connection closed by the server")
					return
				}

				connMutex.Lock()
				if closing {
					connMutex.Unlock()
					return
				}
				connMutex.Unlock()

				slog.Error("unable to read message from the websocket server", "reason", err)
				if closeErr := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); closeErr != nil {
					slog.Error("unable to write close message to the websocket server (reconnect)", "reason", err)
				}

				if err := conn.Close(); err != nil {
					slog.Error("unable to close the websocket connection (reconnect)", "reason", err)
					conn = nil
				}

				// Attempt to reconnect
				index := 1
				for {
					connMutex.Lock()
					if closing {
						connMutex.Unlock()
						return
					}
					connMutex.Unlock()

					err := connectWebSocket(config.GetConfig().WebsocketConfig.WSURL)
					if err != nil {
						slog.Error("failed to reconnect to the websocket server", "reason", err)
						time.Sleep(time.Second * time.Duration(config.GetConfig().WebsocketConfig.WSReconnectInterval))
					} else {
						slog.Debug("reconnected to the websocket server", "attempts", index)
						break
					}
					index++
				}
			}

			if config.GetConfig().DebugMode {
				slog.Debug("Received message:", "message", string(message))
			}

			// Handle business logic
			websocketBusinessLogic(message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			slog.Info("interrupt signal received, closing websocket connection")

			connMutex.Lock()
			closing = true
			connMutex.Unlock()

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				slog.Error("unable to write close message to the websocket server", "reason", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}

			slog.Debug("successfully closed websocket connection")
			return
		}
	}
}

func InitWebsocketClient() error {
	slog.Debug("initializing websocket client")
	if config.GetConfig().WebsocketConfig.WSURL == "" {
		return eris.New("websocket URL not set")
	}
	if config.GetConfig().WebsocketConfig.WSApiToken == "" {
		return eris.New("websocket API Token not set")
	}

	// Establish WebSocket connection
	if err := connectWebSocket(config.GetConfig().WebsocketConfig.WSURL); err != nil {
		slog.Error("unable to establish connection to the websocket server", "reason", err)
		return eris.Wrap(err, "establishing connection to the WebSocket server")
	}

	// Start a goroutine to listen for messages from the WebSocket server
	go listenForMessages()

	slog.Info("successfully established connection to the websocket server")
	return nil
}

// Deprecated: InitWebsocketClient has already implemented the close mechanism. This function will just return nil if called to avoid panic errors
func CloseWebsocketClient() error {
	// connMutex.Lock()
	// closing = true
	// connMutex.Unlock()

	// if conn == nil {
	// 	slog.Debug("successfully closed websocket connection")
	// 	return nil
	// }

	// slog.Debug("closing websocket connection")
	// // Ensure the WebSocket connection is closed
	// if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
	// 	slog.Error("unable to write close message to the websocket server", "reason", err)
	// 	return eris.Wrap(err, "writing close message to the WebSocket server")
	// }

	// connMutex.Lock()
	// if err := conn.Close(); err != nil {
	// 	slog.Error("unable to close the websocket connection", "reason", err)
	// 	return eris.Wrap(err, "closing the WebSocket connection")
	// }
	// connMutex.Unlock()

	// slog.Debug("successfully closed websocket connection")
	return nil
}

func GetWSConn() *websocket.Conn {
	return conn
}

// SendMessage will marshall the provided message before sending it to the websocket server
func SendMessage(ctx context.Context, messageType types.EventList, message interface{}) error {

	var msg Event
	var err error

	msg.Type = messageType
	msg.Payload, err = json.Marshal(message)
	if err != nil {
		slog.Error("unable to marshall websocket message", "reason", err)
		return eris.Wrap(err, "marshalling websocket payload")
	}

	if err = conn.WriteJSON(msg); err != nil {
		slog.Error("unable to send message to the websocket server", "reason", err)
		return eris.Wrap(err, "sending message to the WebSocket server")
	}

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
		slog.Error("unable to unmarshall websocket message", "reason", err)
		return
	}

	// Check if message type is supported, if it is then call the appropriate function in model to handle this
	if handler, exists := eventHandlers[message.Type]; exists {
		handler(message)
	} else {
		slog.Info("unable to handle websocket message, unsupported message type", "received type", message.Type)
	}
}
