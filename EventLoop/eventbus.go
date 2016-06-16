package EventLoop

import (
	"fmt"
	"runtime/debug"
	"sync"

	"golang.org/x/net/websocket"
)

type EventBus struct {
	publisher Publisher

	registers     []*Register
	registersLock sync.RWMutex
}

type Publisher func(bus *EventBus, ws *websocket.Conn)

func NewEventLoop(publisher Publisher) *EventBus {
	return &EventBus{
		publisher: publisher,
		registers: make([]*Register, 0),
	}
}

func (bus *EventBus) Add(ws *websocket.Conn) *Register {
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

	bus.registersLock.Lock()
	defer bus.registersLock.Unlock()
	register := newRegister(bus, ws)
	bus.registers = append(bus.registers, register)

	return register
}

func (bus *EventBus) Remove(register *Register) {
	bus.registersLock.Lock()
	defer bus.registersLock.Unlock()
	for i := 0; i < len(bus.registers); i++ {
		if bus.registers[i] == register {
			bus.registers[i] = bus.registers[len(bus.registers)-1]
			bus.registers = bus.registers[:len(bus.registers)-1]
			return
		}
	}
}

func (bus *EventBus) Publish(event string, args ...interface{}) {
	bus.registersLock.RLock()
	defer bus.registersLock.RUnlock()
	for _, r := range bus.registers {
		r.emit(event, args)
	}
}
