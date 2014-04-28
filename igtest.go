package igtest

import (
	"fmt"
	"strings"
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
	fmt.Println(`Usage:igr [-l -m mode -R type -r file] file ...
 -h show this
 -l show log
 -R the report data type
 -r the report store path
 -m mode in YOE/EYOE/NONE
  YOE assert on execute all command
  EYOE assert on expression
  NONE no assert in command(ony Y/N)`)
}
func Run(args []string) int {
	var log bool = false
	var mode string = "YOE" //EYOE or YOE
	var repo string = "JSON"
	var repof string = ""
	var fset []string = []string{}
	ctx := NewCtx(nil)
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-h" {
			usage()
			return 0
		} else if a == "-l" {
			log = true
		} else if a == "-m" {
			if i != len(args)-1 {
				mode = args[i+1]
				i++
			}
		} else if a == "-R" {
			if i != len(args)-1 {
				repo = args[i+1]
				i++
			}
		} else if a == "-r" {
			if i != len(args)-1 {
				repof = args[i+1]
				i++
			}
		} else if e_reg_kv.MatchString(a) {
			kv := strings.SplitN(a, "=", 2)
			ctx.SET(kv[0], kv[1])
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
	switch repo {
	case "JSON":
		ctx.Mark = NewJsonMarker(repof)
	default:
		fmt.Println("invalid report format:", repo)
		usage()
		return 1
	}
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
	if len(repof) > 0 {
		err := ctx.Mark.Store()
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
	}
	return 0
}

type JsonMarkerExe struct {
	*Ctx
	*Compiler
	Json   string
	Igs    string
	OnExec OnExecedFunc
}

func (j *JsonMarkerExe) Exec() error {
	err := j.Load(j.Igs)
	if err != nil {
		return err
	}
	j.SET("CWD", j.Cwd)
	err = j.CompileAndExec(j.Ctx, j.OnExec)
	if err != nil {
		return err
	}
	err = j.Mark.Store()
	if err != nil {
		return err
	}
	return nil
}

func (j *JsonMarkerExe) E(line string) error {
	l, err := j.NewLine(line, 0, j.OnExec)
	if err != nil {
		return err
	}
	_, err = l.Exec(j.Ctx, true)
	return err
}

func NewJME(json, igs string) *JsonMarkerExe {
	jme := &JsonMarkerExe{}
	jme.Json = json
	jme.Igs = igs
	jme.OnExec = YesOnExeced
	jme.Compiler = NewCompiler()
	jme.Ctx = NewCtx(nil)
	jme.Ctx.Mark = NewJsonMarker(jme.Json)
	return jme
}
