package edago

import (
	"fmt"
	"reflect"
	"sync"
)

type EventBus interface {
	Subscriber
	Publisher
	Controller
}

// fn must be func
type Subscriber interface {
	Subscribe(topic string, fn any) error
	Unsubscribe(topic string, fn any)
}

type Publisher interface {
	Publish(topic string, args ...any)
}

type Controller interface {
	Close(topic string)
}

type defaultEventBus struct {
	handlers map[string][]*handler
	mutex    sync.RWMutex
}

type handler struct {
	fn reflect.Value
	ch chan []reflect.Value
}

func NewEventBus() EventBus {
	return &defaultEventBus{
		handlers: make(map[string][]*handler),
	}
}

func (eb *defaultEventBus) Subscribe(topic string, fn any) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%v must be reflect.Func", fn)
	}

	h := &handler{
		fn: reflect.ValueOf(fn),
		ch: make(chan []reflect.Value),
	}

	go func() {
		for args := range h.ch {
			h.fn.Call(args)
		}
	}()

	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.handlers[topic] = append(eb.handlers[topic], h)

	return nil
}

func (eb *defaultEventBus) Unsubscribe(topic string, fn any) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	if _, ok := eb.handlers[topic]; !ok {
		return
	}
	for i, handler := range eb.handlers[topic] {
		if reflect.DeepEqual(handler.fn, fn) {
			eb.handlers[topic] = append(
				eb.handlers[topic][:i], eb.handlers[topic][i+1:]...,
			)
			close(handler.ch)
			break
		}
	}
	if len(eb.handlers[topic]) == 0 {
		delete(eb.handlers, topic)
	}
}

func (eb *defaultEventBus) Publish(topic string, args ...any) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()
	rargs := make([]reflect.Value, len(args))
	for i, v := range args {
		rargs[i] = reflect.ValueOf(v)
	}
	for _, handler := range eb.handlers[topic] {
		handler.ch <- rargs
	}
}

func (eb *defaultEventBus) Close(topic string) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	for _, handler := range eb.handlers[topic] {
		close(handler.ch)
	}
	delete(eb.handlers, topic)
}
