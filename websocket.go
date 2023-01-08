package fishfish

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WSEventType string

const (
	WSEventTypeDomainCreate = "domain_create"
	WSEventTypeDomainUpdate = "domain_update"
	WSEventTypeDomainDelete = "domain_delete"
	WSEventTypeURLCreate    = "url_create"
	WSEventTypeURLUpdate    = "url_update"
	WSEventTypeURLDelete    = "url_delete"
)

type WSEvent struct {
	Type WSEventType `json:"type"`
	Data any         `json:"data"`
}

type WSCreateDomainData struct {
	Domain      string   `json:"domain"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	Target      string   `json:"target,omitempty"`
}

type WSUpdateDomainData struct {
	Domain      string   `json:"domain"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	Target      string   `json:"target,omitempty"`
	Checked     int64    `json:"checked,omitempty"`
}

type WSDeleteDomainData struct {
	Domain string `json:"domain"`
}

type WSCreateURLData struct {
	URL         string   `json:"url"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	Target      string   `json:"target,omitempty"`
}

type WSUpdateURLData struct {
	URL         string   `json:"url"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	Target      string   `json:"target,omitempty"`
	Checked     int64    `json:"checked,omitempty"`
}

type WSDeleteURLData struct {
	URL string `json:"url"`
}

// This will connect to the FishFish API's WebSocket Stream for real-time updates.
// It will block and write events to the specified channel.
// It is not recommended to use this function directly, as you will have to manually parse events.
// If you want to keep an updated database of domains and urls, use the AutoSync client.
func (c *RawClient) ConnectWS(ctx context.Context, ch chan WSEvent) error {
	if c.defaultAuthType == authTypeNone {
		return fmt.Errorf("authentication is required to use the websocket")
	}

	headers := http.Header{}
	headers.Add("Authorization", c.sessionToken.Token)
	headers.Add("User-Agent", "fishfish-go")

	conn, res, err := websocket.Dial(ctx, "wss://api.fishfish.gg/v1/stream", &websocket.DialOptions{
		HTTPHeader: headers,
	})

	if err != nil {
		fmt.Println(res)
		return fmt.Errorf("could not connect to websocket: %w", err)
	}

	go keepAlive(conn, ctx)

	for {
		var eventData WSEvent
		err := wsjson.Read(ctx, conn, &eventData)

		if err != nil {
			fmt.Println(err)
			// Context was Cancelled
			if errors.Is(err, ctx.Err()) {
				conn.Close(websocket.StatusNormalClosure, "")
				return nil
			}
			if errors.Is(err, &json.UnmarshalTypeError{}) || errors.Is(err, &json.SyntaxError{}) {
				// Invalid Data
				return nil
			}
			// Unexpected error
			return err
		}
		ch <- eventData
	}
}

// Send a byte every 10 seconds to keep the connection from closing
func keepAlive(conn *websocket.Conn, ctx context.Context) {
	for range time.Tick(time.Second * 10) {
		select {
		case <-ctx.Done():
			// Websocket Cancelled
			return
		default:
			// Ping the server
			conn.Write(ctx, websocket.MessageBinary, []byte{1})
		}
	}
}
