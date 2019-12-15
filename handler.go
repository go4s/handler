package handler

import "github.com/gin-gonic/gin"

type Register func(*gin.IRouter) error

var (
	registers = [] Register{}
)

func Hook(engine *gin.IRouter) {
	var err error
	for _, register := range registers {
		err = register(engine)
		if err != nil {
			panic(err)
		}
	}
}

func Add(r Register) {
	registers = append(registers, r)
}
