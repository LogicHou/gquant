package gquant

import (
	"math"

	"github.com/LogicHou/gquant/indicator"
)

const abortIndex int8 = math.MaxInt8 >> 1

type Context struct {
	Ticker   *indicator.Ticker
	handlers []HandlerFunc
	index    int8
}

func newContext(t *indicator.Ticker) *Context {
	return &Context{
		Ticker: t,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}
