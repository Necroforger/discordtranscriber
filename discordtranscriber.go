package discordtranscriber

import (
	"log"
	"net/http"

	"github.com/Necroforger/discordtranscriber/wsmux"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Server server session
type Server struct {
	httpserver *http.Server
	Client     *discordgo.Session
	WSRouter   *wsmux.Router
	Log        bool
}

func (s *Server) log(i ...interface{}) {
	log.Println(i...)
}

// NewServer returns a new server
//     client : discord session to use for information retrieval.
func NewServer(client *discordgo.Session, p string, assets http.FileSystem) *Server {
	s := &Server{
		Client:   client,
		WSRouter: wsmux.NewRouter(),
		Log:      false,
	}

	servemux := http.NewServeMux()
	servemux.Handle("/", http.FileServer(assets))
	servemux.HandleFunc("/websocket/", s.websocketHandler)

	s.httpserver = &http.Server{
		Addr:    "127.0.0.1:" + p,
		Handler: servemux,
	}

	s.addHandlers()

	return s
}

// AddDefaultHandlers adds the default handlers
func (s *Server) addHandlers() {
	r := s.WSRouter
	r.On("send", s.wsSend)
	r.On("channel", s.wsChannel)
	r.On("valid_channel", s.wsValidChannel)
}

// ListenAndServe calls the underlying http server ListenAndServe
func (s *Server) ListenAndServe() error {
	return s.httpserver.ListenAndServe()
}

// ListenAndServeTLS calls the underlying http server ListenAndServeTLS
func (s *Server) ListenAndServeTLS(certfile, keyfile string) error {
	return s.httpserver.ListenAndServeTLS(certfile, keyfile)
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log("websocket err: ", err)
	}

	// Send things like avatar URL on connection
	if err := s.sendInitialData(conn); err != nil {
		s.log(err)
	}

	s.log("New websocket connection from: ", conn.RemoteAddr().String())
	s.readConnection(conn)
}

func (s *Server) sendInitialData(conn *websocket.Conn) error {
	var e error
	wr := func(name, data string) {
		err := writeEvent(conn, name, data)
		if err != nil {
			e = err
		}
	}

	wr("user", stringify(s.Client.State.User))
	wr("avatar", s.Client.State.User.AvatarURL(""))

	return e
}

func (s *Server) readConnection(conn *websocket.Conn) error {
	var e wsmux.Event
	for {
		err := conn.ReadJSON(&e)
		if err != nil {
			s.log("readConnection err: ", err)
			return err
		}

		err = s.execRequest(conn, e)
		if err != nil {
			s.log("readConnection execRequest err: ", err)
		}
	}
}

func (s *Server) execRequest(conn *websocket.Conn, e wsmux.Event) error {
	s.log("websocket request: ", e.Name, e.Data)
	return s.WSRouter.Execute(conn, e)
}
