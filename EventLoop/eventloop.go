package EventLoop

import (
	"fmt"
	"runtime/debug"
	"sync"

	"golang.org/x/net/websocket"
)

type EventBus struct {
	publisher Publisher
	registers []*Register

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
	register := newRegister(ws)
	bus.registers = append(bus.registers, register)

	return register
}

func (bus *EventBus) find(register *Register) int {
	bus.registersLock.RLock()
	defer bus.registersLock.RUnlock()
	for i := 0; i < len(bus.registers); i++ {
		if bus.registers[i] == register {
			return i
		}
	}
	return -1
}

func (bus *EventBus) remove(i int) {
	bus.registersLock.Lock()
	defer bus.registersLock.Unlock()
	bus.registers = append(bus.registers[:i], bus.registers[i+1:]...)
}

func (bus *EventBus) Remove(register *Register) {
	i := bus.find(register)
	if i != -1 {
		bus.remove(i)
	}
}

func (bus *EventBus) Publish(event string, args ...interface{}) {
	for _, r := range bus.registers {
		r.publish(event, args)
	}
}
