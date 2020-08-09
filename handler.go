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

func Add(r interface {
	Version
	Group
}) {
	registers = append(registers, func(g gin.IRouter) {
		api := g.Group(fmt.Sprintf("api/%s/%s", r.Version(), r.Group()))
		{
			// standard
			if handle, found := r.(Watchable); found {
				api.GET("/events", handle.Watch)
			}
			if handle, found := r.(Listable); found {
				api.GET("/", handle.List)
			}
			if handle, found := r.(Creatable); found {
				api.POST("/", handle.Create)
			}
			if handle, found := r.(Deletable); found {
				api.DELETE("/", handle.Delete)
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
	Group() string
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
