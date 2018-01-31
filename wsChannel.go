package discordtranscriber

import (
	"fmt"

	"github.com/Necroforger/discordtranscriber/wsmux"
	"github.com/gorilla/websocket"
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
		return
	}
	writeEvent(conn, "channel", stringify(channel))
}

func (s *Server) wsVoiceChannel(conn *websocket.Conn, e wsmux.Event) {
	channel, err := s.Client.State.Channel(e.Data)
	if err != nil {
		writeEvent(conn, "voice_channel", `""`)
		return
	}
	writeEvent(conn, "voice_channel", stringify(channel))
}
