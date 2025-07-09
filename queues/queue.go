package queues

import (
	"log"
	"sync"
)

type queue struct {
	Name         string
	Messages     []QItem
	mCount       int
	m            *sync.Mutex
	condNotEmpty sync.Cond
}

func newQueue(name string) queue {
	var mut sync.Mutex
	return queue{
		Name:         name,
		Messages:     make([]QItem, 0, 10),
		condNotEmpty: *sync.NewCond(&mut),
		m:            &mut,
	}
}

func (q *queue) Queue(item QItem) {

	q.m.Lock()
	q.Messages = append(q.Messages, item)
	q.mCount = len(q.Messages)
	q.condNotEmpty.Signal()
	q.m.Unlock()

}

func (q *queue) Dequeue() QItem {
	q.m.Lock()
	count := q.mCount

	for !(count > 0) {
		q.condNotEmpty.Wait() // Unlocks the mutex
		count = q.mCount
	}

	item := q.takeNext()
	q.m.Unlock()
	return item
}

func (q *queue) Count() int {
	q.m.Lock()
	count := q.mCount
	q.m.Unlock()
	return count
}

func (q *queue) takeNext() QItem {
	if q.mCount == 0 {
		log.Fatal("Attempted to call takeNext() on on empty queue")
	}

	item := q.Messages[0]
	q.Messages = q.Messages[1:]
	q.mCount = len(q.Messages)

	return item

}
