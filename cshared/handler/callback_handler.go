package handler

import (
	"unsafe"

	"github.com/tkmn0/sylph/pkg/channel"

	"github.com/tkmn0/sylph/cshared/converter"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/cshared/queue"
)

const uintptr_zero = 0

type Event byte

const (
	Call Event = iota
)

type CallbackHandler struct {
	converter *converter.Converter
	// server/client events
	onTransportEventQueues map[unsafe.Pointer]*queue.Queue
	// transport events
	onTransportClosedEventQueues map[unsafe.Pointer]*queue.Queue
	onChannelEventQueues         map[unsafe.Pointer]*queue.Queue
	// channel events
	onChannelClosedEventQueues  map[unsafe.Pointer]*queue.Queue
	onChannelErrorEventQueues   map[unsafe.Pointer]*queue.Queue
	onChannelMessageEventQueues map[unsafe.Pointer]*queue.Queue
	onChannelDataEventQueues    map[unsafe.Pointer]*queue.Queue
}

func NewCallbackHandler() *CallbackHandler {
	return &CallbackHandler{
		converter:                    converter.NewConverter(),
		onTransportEventQueues:       map[unsafe.Pointer]*queue.Queue{},
		onTransportClosedEventQueues: map[unsafe.Pointer]*queue.Queue{},
		onChannelEventQueues:         map[unsafe.Pointer]*queue.Queue{},
		onChannelClosedEventQueues:   map[unsafe.Pointer]*queue.Queue{},
		onChannelErrorEventQueues:    map[unsafe.Pointer]*queue.Queue{},
		onChannelMessageEventQueues:  map[unsafe.Pointer]*queue.Queue{},
		onChannelDataEventQueues:     map[unsafe.Pointer]*queue.Queue{},
	}
}

func (h *CallbackHandler) SetupServerEvents(server *sylph.Server) {
	onTransportQueue := queue.NewQueue()
	h.onTransportEventQueues[unsafe.Pointer(server)] = onTransportQueue

	server.OnTransport(func(t sylph.Transport) {
		h.setupTransportEvents(t)
		onTransportQueue.Enqueue(unsafe.Pointer(&t))
	})
}

func (h *CallbackHandler) SetupClientEvents(client *sylph.Client) {
	onTransportQueue := queue.NewQueue()
	h.onTransportEventQueues[unsafe.Pointer(client)] = onTransportQueue

	client.OnTransport(func(t sylph.Transport) {
		h.setupTransportEvents(t)
		onTransportQueue.Enqueue(unsafe.Pointer(&t))
	})
}

func (h *CallbackHandler) setupTransportEvents(transport sylph.Transport) {
	ptr := unsafe.Pointer(&transport)
	onTransportClosedQueue := queue.NewQueue()
	onChannelQueue := queue.NewQueue()
	h.onTransportClosedEventQueues[ptr] = onTransportClosedQueue
	h.onChannelEventQueues[ptr] = onChannelQueue

	transport.OnClose(func() {
		onTransportClosedQueue.Enqueue(Call)
	})

	transport.OnChannel(func(channel channel.Channel) {
		h.setupChannelEvents(channel)
		onChannelQueue.Enqueue(unsafe.Pointer(&channel))
	})
}

func (h *CallbackHandler) setupChannelEvents(channel channel.Channel) {
	ptr := unsafe.Pointer(&channel)
	onChannelClosedQueue := queue.NewQueue()
	onChannelErrorQueue := queue.NewQueue()
	onChannelMessageQueue := queue.NewQueue()
	onchannelDataQueue := queue.NewQueue()
	h.onChannelClosedEventQueues[ptr] = onChannelClosedQueue
	h.onChannelErrorEventQueues[ptr] = onChannelErrorQueue
	h.onChannelMessageEventQueues[ptr] = onChannelMessageQueue
	h.onChannelDataEventQueues[ptr] = onchannelDataQueue

	channel.OnClose(func() {
		onChannelClosedQueue.Enqueue(Call)
	})

	channel.OnError(func(err error) {
		onChannelErrorQueue.Enqueue(err.Error())
	})

	channel.OnMessage(func(message string) {
		onChannelMessageQueue.Enqueue(message)
	})

	channel.OnData(func(data []byte) {
		onchannelDataQueue.Enqueue(data)
	})
}

func (h *CallbackHandler) ReadOnTransport(p unsafe.Pointer) uintptr {
	q, exists := h.onTransportEventQueues[p]
	if !exists {
		return uintptr_zero
	}

	t := q.Dequeue()
	if t == nil {
		return uintptr_zero
	}

	return uintptr(t.(unsafe.Pointer))
}

func (h *CallbackHandler) ReadOnTransportClosed(p unsafe.Pointer) uintptr {
	q, exists := h.onTransportClosedEventQueues[p]
	if !exists {
		return h.createEmptyEventPointer()
	}

	e := q.Dequeue()
	if e == nil {
		return h.createEmptyEventPointer()
	}

	return h.createByteEventPointer(nil)
}

func (h *CallbackHandler) ReadOnChannel(p unsafe.Pointer) uintptr {
	q, exists := h.onChannelEventQueues[p]
	if !exists {
		return uintptr_zero
	}

	c := q.Dequeue()
	if c == nil {
		return uintptr_zero
	}

	return uintptr(c.(unsafe.Pointer))
}

func (h *CallbackHandler) ReadOnChannelClosed(p unsafe.Pointer) uintptr {
	q, exists := h.onChannelClosedEventQueues[p]
	if !exists {
		return h.createEmptyEventPointer()
	}

	e := q.Dequeue()
	if e == nil {
		return h.createEmptyEventPointer()
	}

	return h.createByteEventPointer(nil)
}

func (h *CallbackHandler) ReadOnChannelError(p unsafe.Pointer) uintptr {
	q, exists := h.onChannelErrorEventQueues[p]
	if !exists {
		return h.createEmptyEventPointer()
	}

	errorMessage := q.Dequeue()
	if errorMessage == nil {
		return h.createEmptyEventPointer()
	}

	return h.createStringEventPointer(errorMessage.(string))
}

func (h *CallbackHandler) ReadOnChannelMessage(p unsafe.Pointer) uintptr {
	q, exists := h.onChannelMessageEventQueues[p]
	if !exists {
		return h.createEmptyEventPointer()
	}

	message := q.Dequeue()
	if message == nil {
		return h.createEmptyEventPointer()
	}

	return h.createStringEventPointer(message.(string))
}

func (h *CallbackHandler) ReadOnChannelData(p unsafe.Pointer) uintptr {
	q, exists := h.onChannelDataEventQueues[p]
	if !exists {
		return h.createEmptyEventPointer()
	}

	data := q.Dequeue()
	if data == nil {
		return h.createEmptyEventPointer()
	}

	return h.createByteEventPointer(data.([]byte))
}

func (h *CallbackHandler) createEmptyEventPointer() uintptr {
	return h.bytesToPointer(h.converter.ConvertBytes(converter.PayloadTypeNull, nil))
}

func (h *CallbackHandler) createByteEventPointer(body []byte) uintptr {
	return h.bytesToPointer(h.converter.ConvertBytes(converter.PlayloadTypeValue, body))
}

func (h *CallbackHandler) createStringEventPointer(body string) uintptr {
	return h.bytesToPointer(h.converter.ConvertString(converter.PlayloadTypeValue, body))
}

func (h *CallbackHandler) bytesToPointer(data []byte) uintptr {
	return uintptr(unsafe.Pointer(&data[0]))
}
