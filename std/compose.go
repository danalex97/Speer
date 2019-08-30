package std

import (
	. "github.com/danalex97/Speer/interfaces"
)


type Composable interface {
	ComposeOnTimeout(timeout int) bool
	ComposeOnPredicate(predicate func () bool) // [TODO]: think how to change interface
}

const predicateRecheck = 10

// Composable implentation that allows composing in a 'chain'. That is, the composable's 
// routine runs until a timeout or a predicate occurs. Afterwards, the parent struct can use the 
// composed node's utilities. 
type ChainComposer struct {
	rc   RoutineCapabilities
	node Node

	r Routine

	ready bool
	fire  bool
}

func NewChainComposer(node Node, rc RoutineCapabilities) *ChainComposer {
	return &ChainComposer{
		rc: rc,
		node: node,

		r: nil,

		ready: false,
		fire: false,
	}
}

func (c *ChainComposer) ComposeOnTimeout(timeout int) bool {
	if c.ready {
		return true
	}

	if !c.fire {
		c.fire = true
		c.rc.Callback(timeout, func() {
			c.ready = true
		})
	}

	c.node.OnNotify()
	return false
}

func (c *ChainComposer) ComposeOnPredicate(pred func() bool) {
	if c.ready {
		c.r.Stop()
		return
	}

	if !c.fire {
		c.fire = true
		c.r = c.rc.Routine(predicateRecheck, func() {
			if pred() {
				c.ready = true
			}
		})
	}

	c.node.OnNotify()
}
