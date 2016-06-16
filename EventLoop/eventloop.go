package EventLoop

import (
	"fmt"
	"runtime/debug"
	"sync"

	"golang.org/x/net/websocket"
)

type EventBus struct {
	publisher  Publisher
	websockets []*websocket.Conn

	addLock sync.Mutex
}

type Publisher func(bus *EventBus, ws *websocket.Conn)

func NewEventLoop(publisher Publisher) *EventBus {
	return &EventBus{
		publisher:  publisher,
		websockets: make([]*websocket.Conn, 0),
	}
}

func (bus *EventBus) Add(ws *websocket.Conn) (*Register, error) {
	go func(bus *EventBus) {
		defer func(bus *EventBus) {
			if err := recover(); err != nil {
				publisherErr := fmt.Errorf("publisher panic: %+v\nStack Strace: %s",
					err, string(debug.Stack()))
				bus.Publish(ErrPublisher, publisherErr)
			}
		}(bus)

		if bus.publisher != nil {
			bus.publisher(bus, ws)
		}

	}(bus)
	bus.addLock.Lock()
	defer bus.addLock.Unlock()

	bus.websockets = append(bus.websockets, ws)

	return *Register{}, nil
}

func (bus *EventBus) Publish(event string, args ...interface{}) {

}
