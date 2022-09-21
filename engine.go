package gquant

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	dialect "github.com/LogicHou/gquant/dialect"
	"github.com/LogicHou/gquant/indicator"
)

type HandlerFunc func(*Context)

type Engine struct {
	middlewares []HandlerFunc
	strategy    *Strategy
}

func New() *Engine {
	engine := &Engine{
		strategy: newStrategy(),
	}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (engine *Engine) Use(middlewares ...HandlerFunc) {
	engine.middlewares = append(engine.middlewares, middlewares...)
}

func (engine *Engine) AddHandle(handler HandlerFunc) {
	engine.strategy.addHandle(handler)
}

func (engine *Engine) Run(platform dialect.Platform) {
	log.Println("engine start")

	err := ListenTicker(engine, platform)
	if err != nil {
		log.Println(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server exiting")
}

func (engine *Engine) Serve(t *indicator.Ticker) {
	c := newContext(t)
	c.handlers = engine.middlewares
	engine.strategy.handle(c)
}
