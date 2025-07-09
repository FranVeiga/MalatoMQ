package queues

import (
	"testing"
	"time"
)

func TestQueueing(t *testing.T) {
	queue := newQueue("test_queue")

	queue.Queue(NewQItem("foo"))
	queue.Queue(NewQItem("bar"))
	queue.Queue(NewQItem("baz"))

	count := queue.Count()
	if count != 3 {
		t.Errorf("Queue message count should be 3, is %v", count)
	}
}

func TestDequeueing(t *testing.T) {

	queue := newQueue("test_queue")

	msgs := [6]string{"foo", "bar", "baz", "test1", "test2", "test3"}
	for _, msg := range msgs {
		queue.Queue(NewQItem(msg))
	}

	for _, msg := range msgs {
		msg_recv := queue.Dequeue()
		if msg != msg_recv.Message {
			t.Errorf("Expected %v, received %v", msg, msg_recv)
		}
	}
}

func TestDequeueBlocking(t *testing.T) {
	queue := newQueue("test_queue")
	// var wg sync.WaitGroup

	go func() {
		time.Sleep(time.Second)
		queue.Queue(NewQItem("should unblock"))
	}()

	item := queue.Dequeue() // should block, function panics if not?
	if item.Message != "should unblock" {
		t.Errorf("Received other message: %v", item.Message)
	}

}
