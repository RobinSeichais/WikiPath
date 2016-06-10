package main

import "errors"

/**
 * Queue Item object
 */
type QueueItem struct {
	value string
	next *QueueItem
}

/**
 * Create a new queue Item
 */
func newQueueItem(value string, next *QueueItem) *QueueItem {
	i := &QueueItem {
		value: value,
		next: next,
	}
	return i
}

/**
 * Queue object (FIFO linked list)
 */
type Queue struct {
	head *QueueItem
	tail *QueueItem
	size int
}

/**
 * Push an item to the queue
 */
func (self *Queue) push(value string) {
	i := newQueueItem(value, nil)
	if self.size == 0 {
		self.head = i
	} else {
		self.tail.next = i
	}
	self.tail = i
	self.size++
}

/**
 * Pop the head of the queue
 */
func (self *Queue) pop() (string, error) {
	if self.size == 0 {
		return "", errors.New("Empty queue.")
	}
	value := self.head.value
	self.head = self.head.next
	self.size--
	return value, nil
}

/**
 * Create a new queue
 */
func newQueue() *Queue {
	q := &Queue {
		head: nil,
		tail: nil,
		size: 0,
	}
	return q
}