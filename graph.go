package main

import (
	"strings"
	"sync"
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
func (n *Node) AddParent(parent *Node) {
	n.parents = append(n.parents, parent)
}

/**
 * Remove all the node's parents
 */
func (n *Node) ClearParents() {
	n.parents = make([]*Node, 0)
}

/**
 * Create a new Node
 */
func NewNode(value string, depth int) (n *Node) {
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
	index map[string]*Node
	mut *sync.Mutex
}

func (g *Graph) Set(value string, depth int, children []string) []string {

	g.mut.Lock()
	defer g.mut.Unlock()

	g.Add(value, depth)
	return g.AddChildren(value, children)
}

func (g *Graph) Add(value string, depth int) *Node { // LOCK

	n, ok := g.index[value]

	if !ok {
		n = NewNode(value, depth)
		g.index[value] = n
	} else if depth < n.depth {
		n.depth = depth
		n.ClearParents()
	}

	return n
}

/**
 * Add children to the given node. The children were retrieved by
 * a working by parsing the respective page.
 */
func (g *Graph) AddChildren(value string, children []string) []string {

	n, ok := g.index[value]
	newChildren := make([]string, 0)
	
	if ok {
		for _, c := range children {
			if child := g.Add(c, n.depth + 1); n.depth < child.depth {
				child.AddParent(n)
				newChildren = append(newChildren, c)
			}
		}
	}

	return newChildren
}

func (g *Graph) GetPath(head string) string {

	nodes := make([]string, 0)

	current, ok := g.index[head]
	for ok {
		nodes = append(nodes, current.value)
		if len(current.parents) == 0 {
			break
		}
		current = current.parents[0]
	}

	return strings.Join(nodes, "\n")
}

func (g Graph) Len() int {
	return len(g.index)
}

/**
 * Create a new Graph
 */
func NewGraph() *Graph {
	return &Graph{
		make(map[string]*Node),
		new(sync.Mutex),
	}
}
