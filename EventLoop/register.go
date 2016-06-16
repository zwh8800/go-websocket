package EventLoop

import "golang.org/x/net/websocket"

type Register struct {
}

type Handler func(bus *EventBus, r *Register, ws *websocket.Conn, args ...interface{})

func (r *Register) Subscribe(event string, handler Handler) {

}

func (r *Register) SubscribeOnce(event string, handler Handler) {

}

func (r *Register) SubscribeAsync(event string, handler Handler) {

}

func (r *Register) SubscribeOnceAsync(event string, handler Handler) {

}
