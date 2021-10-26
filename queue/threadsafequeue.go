package queue

import "unsafe"

type ThreadSafeQueue interface {
	Enqueue(interface{})
	Dequeue() interface{}
}

func NewLockfreeQueue() ThreadSafeQueue {
	q := &lockfreeQueue{}
	n := lfNode{}
	n.pNext.Store(nodeCounter{})
	q.head.Store(nodeCounter{
		pNode: unsafe.Pointer(&n),
	})
	q.tail.Store(nodeCounter{
		pNode: unsafe.Pointer(&n),
	})
	return q
}

func NewLockedQueue() ThreadSafeQueue {
	n := lNode{}
	q := &lockedQueue{
		head: &n,
		tail: &n,
	}
	return q
}
