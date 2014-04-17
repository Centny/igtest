package main

import (
	"fmt"
	"github.com/Centny/igtest"
)

func main() {
	ctx := igtest.NewCtx(nil)
	ctx.ShowLog = true
	c := igtest.Compiler{}
	err := c.Load("script file")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ctx.SET("CWD", c.Cwd)
	err = c.CompileAndExec(ctx, igtest.YesOnExeced)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
