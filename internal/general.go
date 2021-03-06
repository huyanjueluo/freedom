package internal

import (
	"sync"

	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// Initiator .
type Initiator interface {
	CreateParty(relativePath string, handlers ...context.Handler) iris.Party
	BindController(relativePath string, controller interface{}, handlers ...context.Handler)
	BindControllerByParty(party iris.Party, controller interface{})
	BindService(f interface{})
	InjectController(f interface{})
	BindRepository(f interface{})
	BindFactory(f interface{})
	GetService(ctx iris.Context, service interface{})
	// BindInfra if is a singleton, com is an object. if is multiton, com is a function
	BindInfra(single bool, com interface{})
	GetInfra(ctx iris.Context, com interface{})
	// Listen Event
	ListenEvent(eventName string, objectMethod string, appointInfra ...interface{})
	Start(f func(starter Starter))
	Iris() *iris.Application
}

// Starter .
type Starter interface {
	Iris() *iris.Application
	// Asynchronous cache warm-up
	AsyncCachePreheat(f func(repo *Repository))
	// Sync cache warm-up
	CachePreheat(f func(repo *Repository))
	GetSingleInfra(com interface{})
}

// SingleBoot 单例组件启动的接口 e.g. infra里面的kafka
type SingleBoot interface {
	Iris() *iris.Application
	EventsPath(infra interface{}) map[string]string
	Closeing(func())
}

// BeginRequest .
type BeginRequest interface {
	BeginRequest(Worker Worker)
}

var (
	globalApp     *Application
	globalAppOnce sync.Once
	prepares      []func(Initiator)
	starters      []func(starter Starter)
)

// Prepare app.BindController or app.BindControllerByParty.
func Prepare(f func(Initiator)) {
	prepares = append(prepares, f)
}
