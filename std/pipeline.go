package std

import (
	. "github.com/danalex97/Speer/interfaces"
)

type Pipeline interface {
	OnTimeout(node Node, timeout int) Pipeline
	Then(callback func() bool) Pipeline // return true when finished
	Once(callback func() bool) Pipeline // return true if we want to continue execution

	Step()
}

// A node in the Chain Pipeline is a decorated callback. The callback returns if 
// the function was executed or not.
type ChainNode struct {
	callback func() bool

	ready bool
	fire bool
}

// A Chain Pipeline is made out of Chain Nodes.
type ChainPipeline struct {
	RoutineCapabilities

	chain []ChainNode
	executed bool
}


func NewChainPipeline(rc RoutineCapabilities) *ChainPipeline {
	return &ChainPipeline{
		RoutineCapabilities: rc,
	}
}

func emptyNode() ChainNode {
	return ChainNode{
		callback: func() bool { return true },

		ready: false,
		fire: false,
	}
}

func (c *ChainPipeline) OnTimeout(node Node, timeout int) Pipeline {
	n := emptyNode()
	
	n.callback = func() bool {
		// If the routine finished executing, then just skip it.
		if n.ready {
			return false
		}

		// If the callback did not fired(i.e. we are in the first run), then
		// mark that it has fired and stop it when the timeout passed
		if !n.fire {
			n.fire = true
			c.Callback(timeout, func() {
				n.ready = true
			})
		}

		node.OnNotify()
		return true
	}

	c.chain = append(c.chain, n)
	return c
}

func (c *ChainPipeline) Once(callback func() bool) Pipeline {
	n := emptyNode()
	n.callback = func() bool {
		if n.ready {
			return false
		}

		n.ready = true
		return !callback()
	}
	
	c.chain = append(c.chain, n)
	return c
}

func (c *ChainPipeline) Then(callback func() bool) Pipeline {
	n := emptyNode()
	n.callback = func() bool {
		if n.ready {
			return false
		}
		
		n.ready = callback()
		return true
	}

	c.chain = append(c.chain, n)
	return c
}

func (c *ChainPipeline) Step() {
	for _, n := range c.chain {
		if n.callback() {
			return
		}
	}
}
