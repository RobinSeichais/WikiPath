package main

import "sync"


type Item struct {
	Value string
	Depth int
	Score float64
}

func NewItem(value string, depth int, score float64) Item {
	return Item{
		value,
		depth,
		score,
	}
}

type PriorityQueue struct {
	tree []Item
	mut *sync.Mutex
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		make([]Item, 0),
		new(sync.Mutex),
	}
}

func (pq *PriorityQueue) Pop() (string, int) {

	pq.mut.Lock()
	defer pq.mut.Unlock()

	L := pq.Len()
	
	if L == 0 {
		return "", -1
	}
	
	value, depth := pq.tree[0].Value, pq.tree[0].Depth

	pq.tree[0] = pq.tree[L-1]
	pq.tree = pq.tree[:L-1]
	L -= 1
	
	i := 0
	ic := 0

	if 2*i+1 < L && pq.tree[2*i+1].Score > pq.tree[i].Score {
		ic = 2*i+1
	}
	if 2*i+2 < L && pq.tree[2*i+2].Score > pq.tree[i].Score {
		ic = 2*i+2
	}

	for i != ic {
		pq.tree[i], pq.tree[ic] = pq.tree[ic], pq.tree[i]

		i = ic
		if 2*i+1 < L && pq.tree[2*i+1].Score > pq.tree[i].Score {
			ic = 2*i+1
		}
		if 2*i+2 < L && pq.tree[2*i+2].Score > pq.tree[i].Score {
			ic = 2*i+2
		}
	}

	return value, depth
}

func (pq PriorityQueue) Top() float64 {
	return pq.tree[0].Score
}

func (pq *PriorityQueue) Push(it Item) {

	pq.mut.Lock()
	defer pq.mut.Unlock()

	i := pq.Len()
	pq.tree = append(pq.tree, it)

	if i == 0 {
		return
	}

	ip := int(float64(i-1)/2)

	for pq.tree[i].Score > pq.tree[ip].Score {
		pq.tree[i], pq.tree[ip] = pq.tree[ip], pq.tree[i]
		if ip == 0 {
			break
		}
		ip, i = int(float64(i-1)/2), ip
	}

	return
}

func (pq PriorityQueue) Len() int {
	return len(pq.tree)
}