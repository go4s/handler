package handler

import (
    "fmt"
    "sync"
    
    "github.com/gin-gonic/gin"
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

func AddWithMiddleware(r Router, middlewares ...gin.HandlerFunc) {
    registers = append(registers, func(engine gin.IRouter) {
        var (
            g = engine
        )
        if specifier, support := r.(PrefixSpecifier); support {
            g = g.Group(specifier.Prefix())
        }
        api := g.Group(fmt.Sprintf("api/%s", r.Version()))
        {
            {
                // update
                if handle, found := r.(Updatable); found {
                    chains := append(middlewares, handle.Update)
                    api.POST(fmt.Sprintf("%s", r.Group().Plural()), chains...)
                    api.POST(fmt.Sprintf("%s/:id", r.Group().Singular()), chains...)
                }
            }
            {
                // create
                if handle, found := r.(Creatable); found {
                    chains := append(middlewares, handle.Create)
                    api.PUT(fmt.Sprintf("%s", r.Group().Plural()), chains...)
                    api.PUT(fmt.Sprintf("%s/:id", r.Group().Singular()), chains...)
                }
            }
            {
                // delete
                if handle, found := r.(Deletable); found {
                    chains := append(middlewares, handle.Delete)
                    api.DELETE(fmt.Sprintf("%s", r.Group().Singular()), chains...)
                    api.DELETE(fmt.Sprintf("%s/:id", r.Group().Singular()), chains...)
                }
            }
            {
                // list
                if handle, found := r.(Listable); found {
                    chains := append(middlewares, handle.List)
                    api.GET(fmt.Sprintf("%s", r.Group().Plural()), chains...)
                    api.GET(fmt.Sprintf("%s/:id", r.Group().Singular()), chains...)
                }
            }
        }
    })
}

func Add(r Router) {
    registers = append(registers, func(g gin.IRouter) {
        if specifier, support := r.(PrefixSpecifier); support {
            g = g.Group(specifier.Prefix())
        }
        api := g.Group(fmt.Sprintf("api/%s", r.Version()))
        {
            {
                // update
                if handle, found := r.(Updatable); found {
                    api.POST(fmt.Sprintf("%s", r.Group().Plural()), handle.Update)
                    api.POST(fmt.Sprintf("%s/:id", r.Group().Singular()), handle.Update)
                }
            }
            {
                // create
                if handle, found := r.(Creatable); found {
                    api.PUT(fmt.Sprintf("%s", r.Group().Plural()), handle.Create)
                    api.PUT(fmt.Sprintf("%s/:id", r.Group().Singular()), handle.Create)
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

type (
    Watchable interface {
        Watch(*gin.Context)
    }
    Listable interface {
        List(*gin.Context)
    }
    Creatable interface {
        Create(*gin.Context)
    }
    Updatable interface {
        Update(*gin.Context)
    }
    Deletable interface {
        Delete(*gin.Context)
    }
    Customize interface {
        Raw(router gin.IRouter)
    }
)

type (
    PrefixSpecifier interface {
        Prefix() string
    }
    Grouper interface {
        Singular() string
        Plural() string
    }
    Version interface {
        Version() string
    }
    Group interface {
        Group() Grouper
    }
)
