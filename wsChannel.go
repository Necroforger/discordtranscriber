package discordtranscriber

import (
	"fmt"

	"github.com/gorilla/websocket"
	"gitlab.com/koishi/discordtranscriber/wsmux"
)

// ValidChannel returns true if a channel is valid
func (s *Server) wsValidChannel(conn *websocket.Conn, e wsmux.Event) {
	_, err := s.Client.State.Channel(e.Data)
	writeEvent(conn, "valid_channel", fmt.Sprint(err == nil))
}

// ValidChannel returns true if a channel is valid
func (s *Server) wsChannel(conn *websocket.Conn, e wsmux.Event) {
	channel, err := s.Client.State.Channel(e.Data)
	if err != nil {
		writeEvent(conn, "channel", `""`)
	}
	writeEvent(conn, "channel", stringify(channel))
}
