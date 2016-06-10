package main

import (
	"sync"
	"fmt"
)

/**
 * Graph Node object
 */
type Node struct {
	value string
	depth int
	parents []*Node
}

/**
 * Add a parent to the node
 */
func (self *Node) addParent(parent *Node) {
	self.parents = append(self.parents, parent)
}

/**
 * Remove all the node's parents
 */
func (self *Node) clearParents() {
	self.parents = make([]*Node, 0)
}

/**
 * Create a new Node
 */
func newNode(value string, depth int) (n *Node) {
	n = &Node {
		value: value,
		depth: depth,
		parents: make([]*Node, 0),
	}
	return
}

/**
 * Graph object
 *
 * The graph is not a real graph, it's more like a tree since it is
 * an acyclic oriented graph, but we do not allow duplication, so multiple
 * branches might join at some point, but only at the same depth. If two
 * nodes have a common child but are not the same depth, the child is only
 * linked to the lowest depth one.
 *
 * The graph also manages the Queue in a way that make the search "breadth-first"-like
 */
type Graph struct {
	stop chan bool         // Stop workers
	done chan bool         // Notify main loop
	mut *sync.Mutex        // Protects the following
	index map[string]*Node // Node index
	queue *Queue           // Link queue to parse
	endLabel string        // Link we're lokking for
	endNode *Node          // Result
}

/**
 * Add a new node to the grapg and return it.
 * Checks whether it already exists before adding it.
 */
func (self *Graph) add(value string, depth int) *Node {

	n, ok := self.index[value]

	if !ok {
		fmt.Println(value)
		node := newNode(value, depth)
		self.index[value] = node
		self.queue.push(value)

		if value == self.endLabel {
			self.endNode = node
			// Stop all workers
			for i := 0; i < nWorker; i++ {
				self.stop <- true
			}
		}

		return node

	} else if depth < n.depth {
		n.depth = depth
		n.clearParents()
	}
	
	return n
}

/**
 * Add children to the given node. The children were retrieved by
 * a working by parsing the respective page.
 */
func (self *Graph) addChildren(value string, children []string) {

	self.mut.Lock()
	defer self.mut.Unlock()

	n, ok := self.index[value]
	
	if ok {
		for _, v := range children {
			if c := self.add(v, n.depth + 1); n.depth < c.depth {
				c.addParent(n)
			}
		}
	}
}

/**
 * Pops the next page to fetch and parse.
 */
func (self *Graph) pop() (string, error) {
	self.mut.Lock()
	defer self.mut.Unlock()
	return self.queue.pop()
}

/**
 * Create a new Graph
 */
func newGraph() (g *Graph) {
	g = &Graph {
		mut: new(sync.Mutex),
		index: make(map[string]*Node),
		queue: newQueue(),
		stop: make(chan bool, nWorker),
		endLabel: "",
		endNode: nil,
	}
	return
}
