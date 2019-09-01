package std

import (
	. "github.com/danalex97/Speer/interfaces"
)

type Tree interface {
    Node

    IsRoot() bool
    Parent() string
    Children() []string
}

// [TODO] To implement the Tree interface for the order-of-join tree. 
type JoinTree struct {
}

// [TODO] To implement Bellman-Ford BFS tree, given a starting node.
type BfsTree struct {
}

func NewBfsTree(root string) *BfsTree {
}

// [TODO] To implement FloodEcho tree(using ACK/NACKS so we don't need a timeout), given a starting node.
type FloodTree struct {
}

func NewFloodTree(root string) *FloodTree {
}
