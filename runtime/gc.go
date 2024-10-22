package runtime

import (
	"github.com/gogim1/goscript/file"
)

func (s *state) new(value Value) int {
	location := len(s.heap)
	s.heap = append(s.heap, value)
	return location
}

func (s *state) gc() int {
	s.mark()
	s.sweep()
	s.relocate()
	return s.removed
}

type collector struct {
	*state
	values     map[int64]struct{}
	locations  map[int]struct{}
	relocation map[int]int
	removed    int
}

func (c *collector) traverse(value Value, visitor func(Value)) {
	if closure, ok := value.(*Closure); ok {
		if _, visited := c.values[value.GetId()]; !visited {
			if visitor != nil {
				visitor(value)
			}
			c.values[value.GetId()] = struct{}{}
			for _, item := range closure.Env {
				if _, visited := c.locations[item.location]; !visited {
					c.locations[item.location] = struct{}{}
					c.traverse(c.heap[item.location], visitor)
				}
			}
		}
	} else if continuation, ok := value.(*Continuation); ok {
		if _, visited := c.values[value.GetId()]; !visited {
			if visitor != nil {
				visitor(value)
			}
			c.values[value.GetId()] = struct{}{}
			for _, layer := range continuation.Stack {
				if layer.frame {
					for _, item := range *layer.env {
						if _, visited := c.locations[item.location]; !visited {
							c.locations[item.location] = struct{}{}
							c.traverse(c.heap[item.location], visitor)
						}
					}
				}
				if len(layer.args) != 0 {
					for _, v := range layer.args {
						c.traverse(v, visitor)
					}
				}
				if layer.callee != nil {
					c.traverse(layer.callee, visitor)
				}
			}
		}
	}
}

func (c *collector) mark() {
	clear(c.values)
	clear(c.locations)

	c.traverse(NewContinuation(file.SourceLocation{Line: -1, Col: -1}, c.stack), nil)
	if c.value != nil {
		c.traverse(c.value, nil)
	}
}

func (c *collector) sweep() {
	clear(c.relocation)

	if len(c.locations) == len(c.heap) {
		return
	}
	i := 0
	for j := 0; j < len(c.heap); j++ {
		if _, ok := c.locations[j]; ok {
			c.heap[i] = c.heap[j]
			c.relocation[j] = i
			i++
		}
	}
	c.removed = len(c.heap) - i
	c.heap = c.heap[:i]
}

func (c *collector) relocate() {
	clear(c.values)
	clear(c.locations)

	if len(c.relocation) == 0 {
		return
	}
	patcher := func(value Value) {
		if closure, ok := value.(*Closure); ok {
			for i, item := range closure.Env {
				closure.Env[i] = envItem{
					name:     item.name,
					location: c.relocation[item.location],
				}
			}
		} else if continuation, ok := value.(*Continuation); ok {
			for _, layer := range continuation.Stack {
				if layer.frame {
					for i, item := range *layer.env {
						(*layer.env)[i] = envItem{
							name:     item.name,
							location: c.relocation[item.location],
						}
					}
				}
			}
		}
	}

	c.traverse(NewContinuation(file.SourceLocation{Line: -1, Col: -1}, c.stack), patcher)
	if c.value != nil {
		c.traverse(c.value, patcher)
	}
}
