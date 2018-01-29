package discordtranscriber

import (
	"github.com/gorilla/websocket"
	"gitlab.com/koishi/discordtranscriber/wsmux"
)

// SendRequest is data from a SendRequest
type SendRequest struct {
	Content        string `json:"Content"`
	ChannelID      string `json:"ChannelID"`
	VoiceChannelID string `json:"VoiceChannelID"`

	TTS bool `json:"TTS"` // Use discord TTS in text channels
}

func (s *Server) wsSend(conn *websocket.Conn, e wsmux.Event) {
	var sr SendRequest
	err := e.UnmarshalInto(&sr)
	if err != nil {
		s.log("send err: ", err)
		return
	}

	if sr.TTS {
		_, err = s.Client.ChannelMessageSendTTS(sr.ChannelID, sr.Content)
	} else {
		_, err = s.Client.ChannelMessageSend(sr.ChannelID, sr.Content)
	}

	if err != nil {
		s.log("send error: ", err)
	}
}
