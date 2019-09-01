package std

import (
	. "github.com/danalex97/Speer/interfaces"
)

type Coloring interface {
    Node

    Color() int
    MaxColor() int
}

// [TODO] Implement a tree 3-coloring in O(log *n).
type TreeColoring struct {
    Tree
}

// [TODO] Implement a graph delta-coloring.
type GraphColoring struct {
    Topology
}
