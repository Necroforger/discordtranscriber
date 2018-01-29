package wsmux

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

// errors
var (
	ErrHandlerNotFound = errors.New("err: Handler not found")
)

// Event stores information sent from the WebSocket
type Event struct {
	Name string `json:"Name"`
	Data string `json:"Data"`
}

// UnmarshalInto unmarshals the Event Data field into interface 'v'
func (e Event) UnmarshalInto(v interface{}) error {
	return json.Unmarshal([]byte(e.Data), &v)
}

// Handler handles events sent via websocket
type Handler func(*websocket.Conn, Event)

// because functions cannot be compared directly
type handlerInstance struct {
	Handler Handler
}

// Router routes websocket events to functions
type Router struct {
	mu       sync.Mutex
	Handlers map[string][]*handlerInstance
}

// NewRouter returns a new Router
func NewRouter() *Router {
	return &Router{
		Handlers: map[string][]*handlerInstance{},
	}
}

// Route contains route information
type Route struct {
	Name    string
	Handler Handler
}

// On adds an event handler and returns a function that when called, removes the handler
func (w *Router) On(event string, fn Handler) func() {
	i := &handlerInstance{fn}
	return w.addEventListener(event, i)
}

// Execute handler handlers for event
func (w *Router) Execute(conn *websocket.Conn, e Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Handlers[e.Name] == nil || len(w.Handlers[e.Name]) == 0 {
		return ErrHandlerNotFound
	}

	for _, v := range w.Handlers[e.Name] {
		v.Handler(conn, e)
	}

	return nil
}

func (w *Router) addEventListener(event string, i *handlerInstance) func() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Handlers[event] = append(w.Handlers[event], i)
	return func() {
		w.removeEventListener(event, i)
	}
}

func (w *Router) removeEventListener(event string, hi *handlerInstance) {
	w.mu.Lock()
	defer w.mu.Unlock()

	handlers := w.Handlers[event]
	for i := range handlers {
		if handlers[i] == hi {
			handlers = append(handlers[:i], handlers[i+1:]...)
		}
	}
}
