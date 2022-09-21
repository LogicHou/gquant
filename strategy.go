package gquant

import (
	"log"
)

type Strategy struct {
	task HandlerFunc
}

func newStrategy() *Strategy {
	return &Strategy{}
}

func (s *Strategy) addHandle(handler HandlerFunc) {
	s.task = handler
}

func (s *Strategy) handle(c *Context) {
	if s.task == nil {
		log.Println("cannot find main handler")
		return
	}
	c.handlers = append(c.handlers, s.task)
	c.Next()
}
