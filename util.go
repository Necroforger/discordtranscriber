package discordtranscriber

import (
	"encoding/json"

	"github.com/Necroforger/discordtranscriber/wsmux"
	"github.com/gorilla/websocket"
)

func stringify(v interface{}) string {
	b, _ := json.Marshal(v)
	if b == nil {
		return ""
	}
	return string(b)
}

func writeEvent(conn *websocket.Conn, name string, data string) error {
	return conn.WriteJSON(wsmux.Event{
		Name: name,
		Data: data,
	})
}
