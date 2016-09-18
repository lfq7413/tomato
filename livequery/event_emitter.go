package livequery

import "unsafe"

// EventEmitter 事件发射器
type EventEmitter struct {
	events map[string][]HandlerType // TODO 增加并发锁
}

// HandlerType ...
type HandlerType func(args ...string)

// NewEventEmitter ...
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: map[string][]HandlerType{},
	}
}

// Emit 向指定通道中的所有订阅者发送事件消息
func (e *EventEmitter) Emit(messageType string, args ...string) bool {
	if e.events == nil {
		e.events = map[string][]HandlerType{}
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
			handler[len(handler)-1] = nil // 置空，防止内存泄漏
			handler = handler[:len(handler)-1]
		}
	}

	return e
}

// RemoveAllListeners 删除指定通道上类所有监听器，当不指定时，删除所有通道上的所有监听器
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
		listeners[key] = nil // 置空，防止内存泄漏
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

// equal 判断两个闭包函数地址是否相同
func equal(f, s HandlerType) bool {
	addrF := *(*int)(unsafe.Pointer(&f))
	addrS := *(*int)(unsafe.Pointer(&s))
	return addrF == addrS
}
