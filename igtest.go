package igtest

import (
	"fmt"
)

const Version string = "Beta v0.1.0"

func Exec(fname string) error {
	return ExecCtx(fname, NewCtx(nil))
}

func ExecCtx(fname string, ctx *Ctx) error {
	c := Compiler{}
	err := c.Load(fname)
	if err != nil {
		return err
	}
	ctx.SET("CWD", c.Cwd)
	err = c.CompileAndExec(ctx, YesOnExeced)
	if err != nil {
		return err
	}
	return nil
}

func usage() {
	fmt.Println(`Usage:igr [-l -m mode] file ...
 -l show log
 -m mode in YOE/EYOE/NONE
  YOE assert on execute all command
  EYOE assert on expression
  NONE no assert in command(ony Y/N)`)
}
func Run(args []string) int {
	var log bool = false
	var mode string = "YOE" //EYOE
	var fset []string = []string{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-l" {
			log = true
		} else if a == "-m" {
			if i != len(args)-1 {
				mode = args[i+1]
				i++
			}
		} else {
			fset = append(fset, a)
		}
	}
	if len(fset) < 1 {
		usage()
		return 1
	}
	var oexec OnExecedFunc = YesOnExeced
	switch mode {
	case "YOE":
		oexec = YesOnExeced
	case "EYOE":
		oexec = ExpYesOnExeced
	case "NONE":
		oexec = nil
	default:
		fmt.Println("invalid mode:", mode)
		usage()
		return 1
	}
	ctx := NewCtx(nil)
	ctx.ShowLog = log
	c := Compiler{}
	for _, f := range fset {
		err := c.Load(f)
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
		ctx.SET("CWD", c.Cwd)
		err = c.CompileAndExec(ctx, oexec)
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
	}
	return 0
}
