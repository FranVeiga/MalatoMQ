package queues

import (
	"log"
	"sync"
)

type genericQueue struct {
	Name         string
	Messages     []QItem
	mCount       int
	m            *sync.Mutex
	condNotEmpty sync.Cond
}

func newGenericQueue(name string) genericQueue {
	var mut sync.Mutex
	return genericQueue{
		Name:         name,
		Messages:     make([]QItem, 0, 10),
		condNotEmpty: *sync.NewCond(&mut),
		m:            &mut,
	}
}

func (q *genericQueue) Queue(item QItem) {

	q.m.Lock()
	q.Messages = append(q.Messages, item)
	q.mCount = len(q.Messages)
	q.condNotEmpty.Signal()
	q.m.Unlock()

}

func (q *genericQueue) Dequeue() QItem {
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

func (q *genericQueue) Count() int {
	q.m.Lock()
	count := q.mCount
	q.m.Unlock()
	return count
}

func (q *genericQueue) takeNext() QItem {
	if q.mCount == 0 {
		log.Fatal("Attempted to call takeNext() on on empty queue")
	}

	item := q.Messages[0]
	q.Messages = q.Messages[1:]
	q.mCount = len(q.Messages)

	return item

}
