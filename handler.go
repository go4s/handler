package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

type Register func(gin.IRouter)

var (
	once      = new(sync.Once)
	registers []Register
)

func init() {
	once.Do(func() {
		registers = make([]Register, 0)
	})
}

func Hook(engine gin.IRouter) {
	for _, register := range registers {
		register(engine)
	}
}

type Router interface {
	Version
	Group
}

type Grouper interface {
	Singular() string
	Plural() string
}

func Add(r Router) {
	registers = append(registers, func(g gin.IRouter) {
		api := g.Group(fmt.Sprintf("api/%s", r.Version()))
		{
			{
				// create
				if handle, found := r.(Creatable); found {
					api.POST(fmt.Sprintf("%s", r.Group().Plural()), handle.Create)
					api.POST(fmt.Sprintf("%s/:id", r.Group().Singular()), handle.Create)
				}
			}
			{
				// delete
				if handle, found := r.(Deletable); found {
					api.DELETE(fmt.Sprintf("%s/:id", r.Group().Singular()), handle.Delete)
				}
			}
			{
				// list
				if handle, found := r.(Listable); found {
					api.GET(fmt.Sprintf("%s", r.Group().Plural()), handle.List)
					api.GET(fmt.Sprintf("%s/:id", r.Group().Singular()), handle.List)
				}
			}
			{
				// watch
				if handle, found := r.(Watchable); found {
					api.GET(fmt.Sprintf("%s/events", r.Group().Plural()), handle.Watch)
					api.GET(fmt.Sprintf("%s/:id/events", r.Group().Singular()), handle.Watch)
				}
			}
			// use defiled
			if handle, found := r.(Customize); found {
				handle.Raw(g)
			}
		}
	})
}

type Version interface {
	Version() string
}
type Group interface {
	Group() Grouper
}
type Watchable interface {
	Watch(*gin.Context)
}
type Listable interface {
	List(*gin.Context)
}
type Creatable interface {
	Create(*gin.Context)
}
type Deletable interface {
	Delete(*gin.Context)
}
type Customize interface {
	Raw(router gin.IRouter)
}
