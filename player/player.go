/*Package player helps with playing audio*/
package player

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Error values
var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrAlreadyRunning  = errors.New("Session already running")
)

// play result constants
const (
	resNothing = iota // do nothing on finish
	resCleanup        // Remove session from map after playing
)

// Player is an audio player
type Player struct {
	smu sync.Mutex // Sessions mutex
	// Sessions is a map of guild ids to audio sessions
	// You can only have one audio player per guild
	Sessions map[string]*Session
}

// NewPlayer returns a new Player
func NewPlayer() *Player {
	return &Player{
		Sessions: map[string]*Session{},
	}
}

// Session returns the audio Session for a guild
func (p *Player) session(guildID string) (*Session, error) {
	if s, ok := p.Sessions[guildID]; ok {
		return s, nil
	}
	return nil, ErrSessionNotFound
}

// removes a session from the map
func (p *Player) removeSession(guildID string) {
	delete(p.Sessions, guildID)
}

func (p *Player) addSession(session *Session) {
	p.Sessions[session.GuildID()] = session
}

// Play plays an audio session
func (p *Player) Play(dst *discordgo.VoiceConnection, src io.Reader) *Session {
	session := NewSession(dst, src)
	go func() {
		res, _ := session.Play()
		p.cleanup(res, session)
	}()

	p.smu.Lock()
	oldses, err := p.session(dst.GuildID)
	if err == nil { // Stop currently playing session
		oldses.Stop()
		p.removeSession(oldses.GuildID())
	}
	p.addSession(session)
	p.smu.Unlock()

	return session
}

// Resume resumes a player
func (p *Player) Resume(guildID string) error {
	p.smu.Lock()
	defer p.smu.Unlock()

	session, err := p.session(guildID)
	if err != nil {
		return err
	}

	go func() {
		n, _ := session.Play()
		p.cleanup(n, session)
	}()

	return nil
}

func (p *Player) cleanup(res int, session *Session) {
	if res == resCleanup {
		p.smu.Lock()
		p.removeSession(session.GuildID())
		p.smu.Unlock()
	}
}

// Pause pauses a playing Session
func (p *Player) Pause(guildID string) error {
	p.smu.Lock()
	defer p.smu.Unlock()

	session, err := p.session(guildID)
	if err != nil {
		return err
	}
	session.paused = true

	return nil
}

// Stop stops a Session from playing and removes it from the
// Sessions map.
func (p *Player) Stop(guildID string) error {
	p.smu.Lock()
	defer p.smu.Unlock()

	session, err := p.session(guildID)
	if err != nil {
		return err
	}
	session.Stop() // Stop session if it is still running

	// If the session is paused, play it again so it can be cleaned up
	session.mu.Lock()
	if session.paused {
		session.mu.Unlock()
		p.Resume(session.GuildID())
	}

	return nil
}

// Session is an audio session
type Session struct {
	mu      sync.Mutex
	running bool
	paused  bool
	stopped bool
	src     io.Reader                  // Audio source
	dst     *discordgo.VoiceConnection // Audio destination
}

// NewSession returns a pointer to a new session
func NewSession(dst *discordgo.VoiceConnection, src io.Reader) *Session {
	return &Session{
		src: src,
		dst: dst,
	}
}

// ReadNext reads the next opus frame
func (s *Session) ReadNext() ([]byte, error) {
	var size int16
	if err := binary.Read(s.src, binary.LittleEndian, &size); err != nil {
		return nil, err
	}
	var buf = make([]byte, size)
	err := binary.Read(s.src, binary.LittleEndian, buf)
	return buf, err
}

// GuildID returns the Session GuildID
func (s *Session) GuildID() string {
	return s.dst.GuildID
}

// Play plays audio
func (s *Session) Play() (int, error) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return resNothing, ErrAlreadyRunning
	}
	s.running = true
	s.paused = false
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	for {
		s.mu.Lock()
		if s.stopped { // Player needs to stop
			s.mu.Unlock()
			return resCleanup, nil
		}
		if s.paused { // Pause the player
			s.mu.Unlock()
			return resNothing, nil
		}
		s.mu.Unlock()

		frame, err := s.ReadNext()
		if err != nil {
			log.Println(err)
			return resCleanup, err
		}
		s.dst.OpusSend <- frame
	}
}

// Stop stops the session
func (s *Session) Stop() {
	s.mu.Lock()
	s.stopped = true
	s.mu.Unlock()
}

// WaitForFinish waits for the session to finish playing
func (s *Session) WaitForFinish() {
	for {
		s.mu.Lock()
		if !s.running && !s.paused {
			s.mu.Unlock()
			break
		}
		s.mu.Unlock()
	}
}
