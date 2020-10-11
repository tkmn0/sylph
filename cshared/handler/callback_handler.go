package handler

import (
	"unsafe"

	"github.com/tkmn0/sylph"
	"github.com/tkmn0/sylph/cshared/converter"
	"github.com/tkmn0/sylph/cshared/queue"
	"github.com/tkmn0/sylph/pkg/channel"
)

const uintptr_zero = 0

type Event byte

const (
	Call Event = iota
)

type CallbackHandler struct {
	converter                    *converter.Converter
	onTransportEventQueues       map[uintptr]*queue.Queue
	onTransportClosedEventQueues map[uintptr]*queue.Queue
	onChannelEventQueues         map[uintptr]*queue.Queue
	onChannelClosedEventQueues   map[uintptr]*queue.Queue
	onChannelErrorEventQueues    map[uintptr]*queue.Queue
	onChannelMessageEventQueues  map[uintptr]*queue.Queue
	onChannelDataEventQueues     map[uintptr]*queue.Queue
	onTransportHandler           func(t sylph.Transport, p uintptr)
	onChannelHandler             func(c channel.Channel, p uintptr)
}

func NewCallbackHandler() *CallbackHandler {
	return &CallbackHandler{
		converter:                    converter.NewConverter(),
		onTransportEventQueues:       map[uintptr]*queue.Queue{},
		onTransportClosedEventQueues: map[uintptr]*queue.Queue{},
		onChannelEventQueues:         map[uintptr]*queue.Queue{},
		onChannelClosedEventQueues:   map[uintptr]*queue.Queue{},
		onChannelErrorEventQueues:    map[uintptr]*queue.Queue{},
		onChannelMessageEventQueues:  map[uintptr]*queue.Queue{},
		onChannelDataEventQueues:     map[uintptr]*queue.Queue{},
	}
}

func (h *CallbackHandler) OnTransport(handler func(t sylph.Transport, p uintptr)) {
	h.onTransportHandler = handler
}

func (h *CallbackHandler) OnChannel(handler func(c channel.Channel, p uintptr)) {
	h.onChannelHandler = handler
}

func (h *CallbackHandler) SetupServerEvents(server *sylph.Server) {
	onTransportQueue := queue.NewQueue()
	h.onTransportEventQueues[uintptr(unsafe.Pointer(server))] = onTransportQueue

	server.OnTransport(func(t sylph.Transport) {
		if h.onTransportHandler != nil {
			h.onTransportHandler(t, uintptr(unsafe.Pointer(&t)))
		}
		h.setupTransportEvents(&t)
		onTransportQueue.Enqueue(uintptr(uintptr(unsafe.Pointer(&t))))
	})
}

func (h *CallbackHandler) SetupClientEvents(client *sylph.Client) {
	onTransportQueue := queue.NewQueue()
	h.onTransportEventQueues[uintptr(unsafe.Pointer(client))] = onTransportQueue

	client.OnTransport(func(t sylph.Transport) {
		if h.onTransportHandler != nil {
			h.onTransportHandler(t, uintptr(unsafe.Pointer(&t)))
		}
		h.setupTransportEvents(&t)
		onTransportQueue.Enqueue(uintptr(unsafe.Pointer(&t)))
	})
}

func (h *CallbackHandler) setupTransportEvents(transport *sylph.Transport) {
	onTransportClosedQueue := queue.NewQueue()
	onChannelQueue := queue.NewQueue()
	h.onTransportClosedEventQueues[uintptr(unsafe.Pointer(transport))] = onTransportClosedQueue
	h.onChannelEventQueues[uintptr(unsafe.Pointer(transport))] = onChannelQueue

	(*transport).OnClose(func() {
		onTransportClosedQueue.Enqueue(Call)
	})

	(*transport).OnChannel(func(channel channel.Channel) {
		if h.onChannelHandler != nil {
			h.onChannelHandler(channel, uintptr(unsafe.Pointer(&channel)))
		}
		h.setupChannelEvents(&channel)
		onChannelQueue.Enqueue(uintptr(unsafe.Pointer(&channel)))
	})
}

func (h *CallbackHandler) setupChannelEvents(channel *channel.Channel) {
	ptr := uintptr(unsafe.Pointer(channel))
	onChannelClosedQueue := queue.NewQueue()
	onChannelErrorQueue := queue.NewQueue()
	onChannelMessageQueue := queue.NewQueue()
	onchannelDataQueue := queue.NewQueue()
	h.onChannelClosedEventQueues[ptr] = onChannelClosedQueue
	h.onChannelErrorEventQueues[ptr] = onChannelErrorQueue
	h.onChannelMessageEventQueues[ptr] = onChannelMessageQueue
	h.onChannelDataEventQueues[ptr] = onchannelDataQueue

	(*channel).OnClose(func() {
		onChannelClosedQueue.Enqueue(Call)
	})

	(*channel).OnError(func(err error) {
		onChannelErrorQueue.Enqueue(err.Error())
	})

	(*channel).OnMessage(func(message string) {
		onChannelMessageQueue.Enqueue(message)
	})

	(*channel).OnData(func(data []byte) {
		onchannelDataQueue.Enqueue(data)
	})
}

func (h *CallbackHandler) ReadOnTransport(p unsafe.Pointer) uintptr {
	q, exists := h.onTransportEventQueues[uintptr(p)]
	if !exists {
		return uintptr_zero
	}

	t := q.Dequeue()
	if t == nil {
		return uintptr_zero
	}

	return t.(uintptr)
}

func (h *CallbackHandler) ReadOnTransportClosed(p unsafe.Pointer) uintptr {
	q, exists := h.onTransportClosedEventQueues[uintptr(p)]
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
	q, exists := h.onChannelEventQueues[uintptr(p)]
	if !exists {
		return uintptr_zero
	}

	c := q.Dequeue()
	if c == nil {
		return uintptr_zero
	}

	return c.(uintptr)
}

func (h *CallbackHandler) ReadOnChannelClosed(p unsafe.Pointer) uintptr {
	q, exists := h.onChannelClosedEventQueues[uintptr(p)]
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
	q, exists := h.onChannelErrorEventQueues[uintptr(p)]
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
	q, exists := h.onChannelMessageEventQueues[uintptr(p)]
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
	q, exists := h.onChannelDataEventQueues[uintptr(p)]
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
