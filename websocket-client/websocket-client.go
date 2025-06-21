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
	connMutex sync.RWMutex
	closing   bool
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func connect(cfg *config.WebsocketConfig) (*websocket.Conn, error) {
	headers := http.Header{
		"X-API-TOKEN": []string{cfg.WSApiToken},
	}

	// Connect to the WebSocket server with custom headers
	c, _, err := websocket.DefaultDialer.Dial(cfg.WSURL, headers)
	if err != nil {
		slog.Error("error connecting to websocket server", "reason", err)
		return nil, eris.Wrap(err, "error connecting to WebSocket server")
	}
	return c, nil
}

// listenForMessages listens for incoming message from websocket hub. Use DEBUG=true to print the message.
func listenForMessages() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	// Handle read message in a separate goroutine
	go func() {
		defer close(done)

		// This outer loop is the reconnection loop.
		for {
			connMutex.RLock()
			if closing {
				connMutex.RUnlock()
				return
			}
			currentConn := conn
			connMutex.RUnlock()

			// --- Live Session ---
			// This block manages an active connection. It will exit when the connection dies.
			pingTicker := time.NewTicker(pingPeriod)
			currentConn.SetReadDeadline(time.Now().Add(pongWait))
			currentConn.SetPongHandler(func(string) error {
				currentConn.SetReadDeadline(time.Now().Add(pongWait))
				return nil
			})

			// Start a reader goroutine for this connection. It sends errors back on a channel.
			readErrChan := make(chan error, 1)
			go func() {
				for {
					_, message, err := currentConn.ReadMessage()
					if err != nil {
						readErrChan <- err
						return
					}
					websocketBusinessLogic(message)
				}
			}()

			// The event loop for the active session.
			sessionActive := true
			for sessionActive {
				select {
				case err := <-readErrChan:
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						slog.Info("websocket connection closed by server")
					} else {
						slog.Error("read error on websocket, will reconnect", "reason", err)
					}
					sessionActive = false // Exit session loop to trigger reconnect
				case <-pingTicker.C:
					currentConn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := currentConn.WriteMessage(websocket.PingMessage, nil); err != nil {
						slog.Warn("failed to write ping, will reconnect", "reason", err)
						sessionActive = false // Exit session loop to trigger reconnect
					}
				}
			}
			pingTicker.Stop()

			// --- Reconnection Logic ---
			connMutex.RLock()
			if closing {
				connMutex.RUnlock()
				return
			}
			connMutex.RUnlock()

			index := 1
			cfg := config.GetConfig().WebsocketConfig
			reconnectInterval := time.Second * time.Duration(cfg.WSReconnectInterval)
			for {
				slog.Info("attempting to reconnect to websocket server...")
				newConn, err := connect(&cfg)
				if err != nil {
					slog.Error("failed to reconnect to the websocket server", "reason", err)
					time.Sleep(reconnectInterval)
				} else {
					slog.Info("reconnected to the websocket server", "attempts", index)
					connMutex.Lock()
					if closing { // Double-check after acquiring lock
						connMutex.Unlock()
						newConn.Close()
						return
					}
					conn = newConn
					connMutex.Unlock()
					break // Exit reconnect loop, outer loop will start a new session.
				}
				index++
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			slog.Info("interrupt signal received, closing websocket connection")

			if err := CloseWebsocketClient(); err != nil {
				slog.Error("error during websocket client shutdown", "reason", err)
			}

			select {
			case <-done:
				slog.Debug("read loop terminated gracefully")
			case <-time.After(time.Second):
				slog.Warn("timed out waiting for read loop to close")
			}
			return
		}
	}
}

func InitWebsocketClient() error {
	slog.Debug("initializing websocket client")
	cfg := config.GetConfig().WebsocketConfig
	if cfg.WSURL == "" {
		return eris.New("websocket URL not set")
	}
	if cfg.WSApiToken == "" {
		return eris.New("websocket API Token not set")
	}

	// Establish WebSocket connection
	newConn, err := connect(&cfg)
	if err != nil {
		slog.Error("unable to establish connection to the websocket server", "reason", err)
		return eris.Wrap(err, "establishing connection to the WebSocket server")
	}
	connMutex.Lock()
	conn = newConn
	closing = false
	connMutex.Unlock()

	// Start a goroutine to listen for messages from the WebSocket server
	go listenForMessages()

	slog.Info("successfully established connection to the websocket server")
	return nil
}

func CloseWebsocketClient() error {
	connMutex.Lock()
	if closing {
		connMutex.Unlock()
		return nil // Already closing or closed
	}
	closing = true
	c := conn
	connMutex.Unlock()

	if c == nil {
		return nil // Nothing to close
	}

	slog.Debug("closing websocket connection")
	// Send a close message to the peer.
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		slog.Warn("failed to write close message, connection may be broken", "reason", err)
	}

	// Close the underlying network connection.
	// This will cause ReadMessage in the listener to return an error, which will then
	// check the 'closing' flag and exit the goroutine.
	return c.Close()
}

func GetWSConn() *websocket.Conn {
	connMutex.RLock()
	defer connMutex.RUnlock()
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

	connMutex.RLock()
	defer connMutex.RUnlock()

	if conn == nil {
		return eris.New("cannot send message: websocket connection not available")
	}

	// Set a deadline for the write operation to prevent it from blocking indefinitely.
	if err := conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return eris.Wrap(err, "setting write deadline")
	}

	// WriteJSON is safe for concurrent use. The mutex here protects the `conn` variable itself.
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
