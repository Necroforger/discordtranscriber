package discordtranscriber

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"gitlab.com/koishi/discordtranscriber/wsmux"
)

// stringify is a shortcut for marshalling structs to a string
func stringify(v interface{}) string {
	b, _ := json.Marshal(v)
	if b == nil {
		return ""
	}
	return string(b)
}

// writeEvent provides a shortcut for sending websocket events
//     conn : websocket connection
//     data : data to send
func writeEvent(conn *websocket.Conn, name string, data string) error {
	return conn.WriteJSON(wsmux.Event{
		Name: name,
		Data: data,
	})
}
