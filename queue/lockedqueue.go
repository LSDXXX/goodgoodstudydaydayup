package queue

import "sync"

type lNode struct {
	next *lNode
	data interface{}
}

type lockedQueue struct {
	hLock sync.Mutex
	tLock sync.Mutex

	head *lNode
	tail *lNode
}

func (l *lockedQueue) Enqueue(data interface{}) {
	n := lNode{
		next: nil,
		data: data,
	}
	l.tLock.Lock()
	l.tail.next = &n
	l.tail = &n
	l.tLock.Unlock()
}

func (l *lockedQueue) Dequeue() interface{} {
	l.hLock.Lock()
	defer l.hLock.Unlock()
	if l.head.next == nil {
		return nil
	}
	data := l.head.next.data
	l.head = l.head.next
	return data
}
