package EventLoop

import (
	"sync"

	"golang.org/x/net/websocket"
)

type internalHandler struct {
	IsOnce      bool
	IsAsync     bool
	Handler     Handler
	OnceHandler sync.Once
	Event       string
}

type Handler func(bus *EventBus, r *Register, ws *websocket.Conn, args ...interface{})

type Register struct {
	bus *EventBus
	ws  *websocket.Conn

	// TODO(zhangweihu): 使用更细粒度的锁
	handlers     map[string][]*internalHandler
	handlersLock sync.RWMutex
}

func newRegister(bus *EventBus, ws *websocket.Conn) *Register {
	return &Register{
		bus:      bus,
		ws:       ws,
		handlers: make(map[string][]*internalHandler),
	}
}

func (r *Register) emit(event string, args ...interface{}) {
	r.handlersLock.RLock()
	defer r.handlersLock.RUnlock()
	handlers, ok := r.handlers[event]
	if !ok {
		return
	}
	for _, h := range handlers {
		closure := func() {
			if h.IsAsync {
				go h.Handler(r.bus, r, r.ws, args)
			} else {
				h.Handler(r.bus, r, r.ws, args)
			}
		}

		if h.IsOnce {
			h.OnceHandler.Do(closure)
		} else {
			closure()
		}
	}
}

func (r *Register) Publish(event string, args ...interface{}) {
	r.bus.registersLock.RLock()
	defer r.bus.registersLock.RUnlock()
	for _, register := range r.bus.registers {
		if r != register {
			register.emit(event, args)
		}
	}
}

func (r *Register) subscribe(event string, handler *internalHandler) {
	r.handlersLock.Lock()
	defer r.handlersLock.Unlock()
	handlers, ok := r.handlers[event]
	if !ok {
		handlers = make([]*internalHandler, 0)
	}

	handlers = append(handlers, handler)

	r.handlers[event] = handlers
}

func (r *Register) Subscribe(event string, handler Handler) {
	r.subscribe(event, &internalHandler{
		IsOnce:  false,
		IsAsync: false,
		Handler: handler,
		Event:   event,
	})
}

func (r *Register) SubscribeOnce(event string, handler Handler) {
	r.subscribe(event, &internalHandler{
		IsOnce:  true,
		IsAsync: false,
		Handler: handler,
		Event:   event,
	})
}

func (r *Register) SubscribeAsync(event string, handler Handler) {
	r.subscribe(event, &internalHandler{
		IsOnce:  false,
		IsAsync: true,
		Handler: handler,
		Event:   event,
	})
}

func (r *Register) SubscribeOnceAsync(event string, handler Handler) {
	r.subscribe(event, &internalHandler{
		IsOnce:  true,
		IsAsync: true,
		Handler: handler,
		Event:   event,
	})
}
