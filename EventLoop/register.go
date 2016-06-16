package EventLoop

import "golang.org/x/net/websocket"

type internalHandler struct {
	IsOnce  bool
	IsAsync bool
	Handler Handler
	Event   string
}

type Handler func(bus *EventBus, r *Register, ws *websocket.Conn, args ...interface{})

type Register struct {
	ws *websocket.Conn

	handlers []*internalHandler
}

func newRegister(ws *websocket.Conn) *Register {
	return &Register{
		ws:       ws,
		handlers: make([]*internalHandler, 0),
	}
}

func (r *Register) publish(event string, args ...interface{}) {

}

func (r *Register) Subscribe(event string, handler Handler) {

}

func (r *Register) SubscribeOnce(event string, handler Handler) {

}

func (r *Register) SubscribeAsync(event string, handler Handler) {

}

func (r *Register) SubscribeOnceAsync(event string, handler Handler) {

}
