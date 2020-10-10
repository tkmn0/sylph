package queue

import "sync"

type Queue struct {
	queue []interface{}
	mu    sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{
		queue: make([]interface{}, 0),
	}
}

func (q *Queue) Enqueue(value interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, value)
}

func (q *Queue) Dequeue() interface{} {
	q.mu.RLock()
	defer func() {
		if len(q.queue) > 0 {
			q.queue = q.queue[1:]
		} else {
			q.queue = q.queue[0:]
		}
		q.mu.RUnlock()
	}()

	if len(q.queue) > 0 {
		return q.queue[0]
	} else {
		return nil
	}
}
