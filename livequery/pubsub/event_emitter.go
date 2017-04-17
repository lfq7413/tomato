package pubsub

import (
	"sync"
	"unsafe"
)

// EventEmitter 事件发射器
type EventEmitter struct {
	mutex  sync.Mutex
	events map[string]map[int]HandlerType
}

// NewEventEmitter ...
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: map[string]map[int]HandlerType{},
	}
}

// Init ...
func (e *EventEmitter) Init() {
	e.events = map[string]map[int]HandlerType{}
}

// Emit 向指定通道中的所有订阅者发送事件消息
func (e *EventEmitter) Emit(messageType string, args ...string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.events == nil {
		e.events = map[string]map[int]HandlerType{}
	}

	if handler, ok := e.events[messageType]; ok {
		if handler == nil {
			return false
		}
		for _, listener := range handler {
			// TODO 完善多线程发送逻辑
			go listener(args...)
		}
		return true
	}
	return false
}

// AddListener 向指定通道添加订阅者的消息监听器
func (e *EventEmitter) AddListener(messageType string, listener HandlerType) *EventEmitter {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.events == nil {
		e.events = map[string]map[int]HandlerType{}
	}

	addr := *(*int)(unsafe.Pointer(&listener))
	if handler, ok := e.events[messageType]; ok {
		handler[addr] = listener
		e.events[messageType] = handler
	} else {
		e.events[messageType] = map[int]HandlerType{addr: listener}
	}

	return e
}

// On 向指定通道添加订阅者的消息监听器，同 AddListener
func (e *EventEmitter) On(messageType string, listener HandlerType) *EventEmitter {
	return e.AddListener(messageType, listener)
}

// Once 添加只执行一次的监听器
func (e *EventEmitter) Once(messageType string, listener HandlerType) *EventEmitter {
	fired := false

	// 包装订阅者的监听器，当包装监听器得到执行时，立即删除自身，并执行订阅者的监听器
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

// RemoveListener 删除指定通道上的指定监听器
func (e *EventEmitter) RemoveListener(messageType string, listener HandlerType) *EventEmitter {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.events == nil {
		return e
	}

	addr := *(*int)(unsafe.Pointer(&listener))
	if handler, ok := e.events[messageType]; ok {
		if _, ok := handler[addr]; ok == false {
			return e
		}
		if len(handler) == 1 {
			delete(handler, addr)
			delete(e.events, messageType)
		} else {
			delete(handler, addr)
		}
	}

	return e
}

// RemoveAllListeners 删除指定通道上类所有监听器，当不指定时，删除所有通道上的所有监听器
func (e *EventEmitter) RemoveAllListeners(messageType string) *EventEmitter {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.events == nil {
		return e
	}

	if messageType == "" {
		for key := range e.events {
			e.RemoveAllListeners(key)
		}
		e.events = map[string]map[int]HandlerType{}
		return e
	}

	listeners := e.events[messageType]
	for addr := range listeners {
		delete(listeners, addr)
	}
	delete(e.events, messageType)

	return e
}

// Listeners ...
func (e *EventEmitter) Listeners(messageType string) map[int]HandlerType {
	if e.events == nil {
		return map[int]HandlerType{}
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
