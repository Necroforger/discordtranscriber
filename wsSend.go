package discordtranscriber

import (
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/jonas747/dca"
	"gitlab.com/koishi/discordtranscriber/wsmux"
)

// SendRequest is data from a SendRequest
type SendRequest struct {
	Content   string `json:"Content"`
	GuildID   string `json:"GuildID"`
	ChannelID string `json:"ChannelID"`

	VoiceChannelID string       `json:"VoiceChannelID"`
	VoiceOptions   VoiceOptions `json:"VoiceOptions"`
	TTS            bool         `json:"TTS"`   // Use discord TTS in text channels
	Voice          bool         `json:"Voice"` // Use espeak to synthesize voice
	Text           bool         `json:"Text"`  // Use text chat
}

func (s *Server) wsSend(conn *websocket.Conn, e wsmux.Event) {
	var sr SendRequest
	err := e.UnmarshalInto(&sr)
	if err != nil {
		s.log("send err: ", err)
		return
	}

	if sr.Text {
		if sr.TTS {
			_, err = s.Client.ChannelMessageSendTTS(sr.ChannelID, sr.Content)
		} else {
			_, err = s.Client.ChannelMessageSend(sr.ChannelID, sr.Content)
		}
		if err != nil {
			s.log("send error: ", err)
		}
	}

	// If we are using a voice channel
	if sr.Voice && sr.VoiceChannelID != "" {
		s.PlayVoice(sr.GuildID, sr.VoiceChannelID, NewVoiceOptions(), sr.Content)
	}
}

// VoiceOptions ...
type VoiceOptions struct {
	Voice     string `json:"Voice"`     // espeak voice file to use
	Pitch     string `json:"Pitch"`     // Pitch of voice
	Speed     string `json:"Speed"`     // word rate in wpm
	Amplitude string `json:"Amplitude"` // Amplitude of voice, 0-200
}

// NewVoiceOptions returns default voice options
func NewVoiceOptions() *VoiceOptions {
	return &VoiceOptions{
		Voice:     "en-US+whisper",
		Pitch:     "50",
		Speed:     "130",
		Amplitude: "100",
	}
}

// PlayVoice uses espeak to play a voice in the given voice channel
func (s *Server) PlayVoice(guildID, channelID string, o *VoiceOptions, content string) error {
	if o == nil {
		o = NewVoiceOptions()
	}
	vc, err := s.Client.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	s.log("Playing voice: " + content)

	args := []string{
		"--stdout",
		"-v", o.Voice,
		"-p", o.Pitch,
		"-s", o.Speed,
		"-a", o.Amplitude,
		content,
	}

	espeak := exec.Command("espeak", args...)
	espeakOut, err := espeak.StdoutPipe()
	if err != nil {
		s.log("PlayVoice error getting espeak stdout pipe: ", err)
		return err
	}
	err = espeak.Start()
	if err != nil {
		s.log("espeak error: ", err)
	}

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	encodeSession, err := dca.EncodeMem(espeakOut, opts)
	if err != nil {
		return err
	}

	s.Player.Play(vc, encodeSession)

	return nil
}
