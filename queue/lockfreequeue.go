package queue

import (
	"sync/atomic"
	"unsafe"
)

type nodeCounter struct {
	pNode unsafe.Pointer
	count uint64
}

func (n *nodeCounter) node() *lfNode {
	return (*lfNode)(n.pNode)
}

type lfNode struct {
	pNext atomic.Value
	data  interface{}
}

func (n *lfNode) next() nodeCounter {
	return n.pNext.Load().(nodeCounter)
}

type lockfreeQueue struct {
	tail atomic.Value
	head atomic.Value
}

func (l *lockfreeQueue) Enqueue(data interface{}) {
	newNode := lfNode{
		data: data,
	}
	newNode.pNext.Store(nodeCounter{})
	var tail nodeCounter
	for {
		tail = l.tail.Load().(nodeCounter)
		next := tail.node().next()
		if next.pNode == nil {
			if tail.node().pNext.CompareAndSwap(next, nodeCounter{
				pNode: unsafe.Pointer(&newNode),
				count: next.count + 1,
			}) {
				break
			}
		} else {
			l.tail.CompareAndSwap(tail, nodeCounter{
				pNode: next.pNode,
				count: tail.count + 1,
			})
		}
	}
	l.tail.CompareAndSwap(tail, nodeCounter{
		pNode: unsafe.Pointer(&newNode),
		count: tail.count + 1,
	})
}

func (l *lockfreeQueue) Dequeue() interface{} {
	for {
		head := l.head.Load().(nodeCounter)
		tail := l.tail.Load().(nodeCounter)
		next := head.node().next()
		if head.pNode == tail.pNode {
			if next.pNode == nil {
				return nil
			}
			l.tail.CompareAndSwap(tail, nodeCounter{
				pNode: next.pNode,
				count: tail.count + 1,
			})
		} else {
			data := next.node().data
			if l.head.CompareAndSwap(head, nodeCounter{
				pNode: next.pNode,
				count: head.count + 1,
			}) {
				return data
			}
		}
	}
}
