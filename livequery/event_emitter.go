package livequery

import "unsafe"

// EventEmitter ...
type EventEmitter struct {
	events map[string][]HandlerType
}

// HandlerType ...
type HandlerType func(args ...string)

// NewEventEmitter ...
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: map[string][]HandlerType{},
	}
}

// Emit ...
func (e *EventEmitter) Emit(messageType string, args ...string) bool {
	if e.events == nil {
		e.events = map[string][]HandlerType{}
	}

	if handler, ok := e.events[messageType]; ok {
		if handler == nil {
			return false
		}
		for _, listener := range handler {
			listener(args...)
		}
		return true
	}
	return false
}

// AddListener ...
func (e *EventEmitter) AddListener(messageType string, listener HandlerType) *EventEmitter {
	if e.events == nil {
		e.events = map[string][]HandlerType{}
	}

	if handler, ok := e.events[messageType]; ok {
		handler = append(handler, listener)
		e.events[messageType] = handler
	} else {
		e.events[messageType] = []HandlerType{listener}
	}

	return e
}

// On ...
func (e *EventEmitter) On(messageType string, listener HandlerType) *EventEmitter {
	return e.AddListener(messageType, listener)
}

// Once ...
func (e *EventEmitter) Once(messageType string, listener HandlerType) *EventEmitter {
	fired := false

	var wrapListener HandlerType
	wrapListener = func(args ...string) {
		e.RemoveListener(messageType, wrapListener)
		if fired == false {
			fired = true
			listener(args...)
		}
	}
	e.On(messageType, wrapListener)
	return e
}

// RemoveListener ...
func (e *EventEmitter) RemoveListener(messageType string, listener HandlerType) *EventEmitter {
	if e.events == nil {
		return e
	}

	if handler, ok := e.events[messageType]; ok {
		position := -1
		for p, l := range handler {
			if equal(listener, l) {
				position = p
				break
			}
		}
		if position < 0 {
			return e
		}
		if len(handler) == 1 {
			handler[0] = nil
			delete(e.events, messageType)
		} else {
			handler[position] = handler[len(handler)-1]
			handler[len(handler)-1] = nil
			handler = handler[:len(handler)-1]
		}
	}

	return e
}

// RemoveAllListeners ...
func (e *EventEmitter) RemoveAllListeners(messageType string) *EventEmitter {
	if e.events == nil {
		return e
	}

	if messageType == "" {
		for key := range e.events {
			e.RemoveAllListeners(key)
		}
		e.events = map[string][]HandlerType{}
		return e
	}

	listeners := e.events[messageType]
	for key := range listeners {
		listeners[key] = nil
	}
	delete(e.events, messageType)

	return e
}

// Listeners ...
func (e *EventEmitter) Listeners(messageType string) []HandlerType {
	if e.events == nil {
		return []HandlerType{}
	}
	return e.events[messageType]
}

// ListenerCount ...
func (e *EventEmitter) ListenerCount(messageType string) int {
	if e.events == nil {
		return 0
	}
	return len(e.events[messageType])
}

func equal(f, s HandlerType) bool {
	addrF := *(*int)(unsafe.Pointer(&f))
	addrS := *(*int)(unsafe.Pointer(&s))
	return addrF == addrS
}
