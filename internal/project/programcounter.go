package project

import (
	"sync"

	"github.com/microsoft/typescript-go/internal/compiler"
)

type programCounter struct {
	mu   sync.Mutex
	refs map[*compiler.Program]int32
}

func (c *programCounter) Add(program *compiler.Program) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.refs == nil {
		c.refs = make(map[*compiler.Program]int32)
	}
	c.refs[program] = 1
}

func (c *programCounter) Ref(program *compiler.Program) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	count, ok := c.refs[program]
	if !ok {
		return false
	}
	c.refs[program] = count + 1
	return true
}

func (c *programCounter) Deref(program *compiler.Program) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	count, ok := c.refs[program]
	if !ok {
		return false
	}
	count--
	if count < 0 {
		panic("program reference count went below zero")
	}
	if count == 0 {
		delete(c.refs, program)
		return true
	}
	c.refs[program] = count
	return false
}

func (c *programCounter) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.refs)
}
