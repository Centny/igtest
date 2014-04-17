package igtest

import (
	"path/filepath"
	"strconv"
	"strings"
)

type Line struct {
	T        string
	L        string
	Pre      *Line
	Args     []string
	Lines    []*Line
	C        *Compiler
	OnExeced OnExecedFunc
}

func (l *Line) addl(line *Line) {
	line.Pre = l
	l.Lines = append(l.Lines, line)
}

func (l *Line) Exec(ctx *Ctx, left bool) (interface{}, error) {
	v, err := l.exec(ctx, left)
	if l.OnExeced == nil {
		return v, err
	} else {
		return l.OnExeced(l, left, v, err)
	}
}

func (l *Line) exec(ctx *Ctx, left bool) (interface{}, error) {
	switch l.T {
	case "BC":
		if len(l.Args) < 1 {
			return nil, Err("Usage:BC $a*$b,but:%v", l.L)
		}
		return ctx.BC(l.Args[0])
	case "SET":
		if len(l.Args) < 2 {
			return nil, Err("Usage:SET path value,but:%v", l.L)
		}
		return true, ctx.SET(l.Args[0], l.Args[1])
	case "GET":
		if len(l.Args) < 1 {
			return nil, Err("Usage:GET path,but:%v", l.L)
		}
		return ctx.GET(l.Args[0]), nil
	case "HR":
		if len(l.Args) < 2 {
			return nil, Err("Usage:HR method url key=val ... ,but:%v", l.L)
		}
		_, v, err := ctx.HR(l.Args...)
		return v, err
	case "HP":
		if len(l.Args) < 1 {
			return nil, Err("Usage:HP url key=val ... ,but:%v", l.L)
		}
		_, v, err := ctx.HP(l.Args...)
		return v, err
	case "HG":
		if len(l.Args) < 1 {
			return nil, Err("Usage:HG url key=val ... ,but:%v", l.L)
		}
		_, v, err := ctx.HG(l.Args...)
		return v, err
	case "EX":
		if len(l.Args) < 1 {
			return nil, Err("Usage:EX cmd args ... ,but:%v", l.L)
		}
		return ctx.EX(l.Args...)
	case "P":
		ctx.P(l.Args...)
		return true, nil
	case "SUB":
		return l.Sub(ctx, left)
	case "FOR":
		return l.For(ctx, left)
	case "EXP":
		return l.Exp(ctx, left)
	case "Y":
		return l.Y(ctx, left)
	case "N":
		return l.N(ctx, left)
	case "ASSIGN":
		return l.Assign(ctx, left)
	default:
		return nil, Err("invalid line type(%v):%v", l.T, l.L)
	}
}

func (l *Line) for_r(ctx *Ctx, left bool) (interface{}, error) {
	rgs := strings.Split(l.Args[2], "~")
	if len(rgs) < 2 {
		return nil, Err("invalid for(%v) range:%v", l.L, l.Args[2])
	}
	beg, err := strconv.ParseInt(ctx.Compile(rgs[0]), 10, 32)
	if err != nil {
		return nil, Err("for(%v) range(%v) error:%v", l.L, l.Args[2], err.Error())
	}
	end, err := strconv.ParseInt(ctx.Compile(rgs[1]), 10, 32)
	if err != nil {
		return nil, Err("for(%v) range(%v) error:%v", l.L, l.Args[2], err.Error())
	}
	for i := beg; i < end; i++ {
		ctx.SET(l.Args[0], i)
		for _, line := range l.Lines {
			_, err := line.Exec(ctx, true)
			if err != nil {
				return nil, err
			}
		}
	}
	return true, nil
}

func (l *Line) for_in(ctx *Ctx, left bool) (interface{}, error) {
	ins := strings.Split(l.Args[2], "~")
	for _, in := range ins {
		ctx.SET(l.Args[0], ctx.Compile(in))
		for _, line := range l.Lines {
			_, err := line.Exec(ctx, true)
			if err != nil {
				return nil, err
			}
		}
	}
	return true, nil
}

func (l *Line) For(ctx *Ctx, left bool) (interface{}, error) {
	if len(l.Args) < 3 {
		return nil, Err("Usage:FOR k R|IN 1~5|A~b~c,but:%v", l.L)
	}
	if l.Args[1] == "R" {
		return l.for_r(ctx, left)
	} else if l.Args[1] == "IN" {
		return l.for_in(ctx, left)
	} else {
		return nil, Err("invalid for(%v) option:%v", l.L, l.Args[1])
	}
}

func (l *Line) Sub(ctx *Ctx, left bool) (interface{}, error) {
	if len(l.Args) < 1 {
		return nil, Err("Usage:SUB path k=v... ,but:%v", l.L)
	}
	var ignore_err bool = false
	nctx := NewCtx(ctx)
	for _, arg := range l.Args[0:] {
		if arg == "-cookie" { //setting cookie to sub
			nctx.H = ctx.H
			continue
		}
		if arg == "-ig" {
			ignore_err = true
			continue
		}
		kvs := strings.SplitN(arg, "=", 2)
		if len(kvs) < 2 {
			continue
		}
		nctx.Kvs[ctx.Compile(kvs[0])] = ctx.Compile(kvs[1])
	}
	c := Compiler{}
	tpath := ctx.Compile(l.Args[0])
	if !filepath.IsAbs(tpath) {
		tpath = filepath.Join(l.C.Cwd, tpath)
	}
	err := c.Load(tpath)
	if err != nil {
		return nil, err
	}
	err = c.CompileAndExec(nctx, l.OnExeced)
	if ignore_err {
		return true, nil
	} else {
		return true, err
	}
}
func (l *Line) Assign(ctx *Ctx, left bool) (interface{}, error) {
	if len(l.Args) < 2 {
		return nil, Err("Invalid assign:%v", l.L)
	}
	var v interface{} = nil
	var err error = nil
	if c_reg_EXP.MatchString(l.Args[1]) {
		nl, _ := l.C.NewLine(l.Args[1], l.OnExeced)
		v, err = nl.Exec(ctx, false)
		if err != nil {
			return nil, err
		}
	} else {
		v = l.Args[1]
	}
	err = ctx.SET(strings.Trim(l.Args[0], " \t$()"), v)
	return true, err
}

func (l *Line) Exp(ctx *Ctx, left bool) (interface{}, error) {
	if len(l.Args) < 1 {
		return nil, Err("Invalid expression:%v", l.L)
	}
	exp := l.Args[0]
	if e_reg_c.MatchString(exp) {
		nl, err := l.C.NewLine(strings.Trim(exp, " \t @[]"), l.OnExeced)
		if err != nil {
			return nil, err
		}
		return nl.Exec(ctx, left)
	} else if e_reg_m.MatchString(exp) {
		return ctx.BC(strings.Trim(exp, " \t@{}"))
	} else if e_reg_s.MatchString(exp) {
		return ctx.Join(strings.Split(strings.Trim(exp, " \t @()"), "+")...), nil
	} else if e_reg_v.MatchString(exp) {
		return ctx.Compile(exp), nil
	} else {
		return nil, Err("invalid express(%v)", exp)
	}
}

func (l *Line) Y(ctx *Ctx, left bool) (interface{}, error) {
	if len(l.Args) < 1 {
		return nil, Err("Usage:Y $a,but:%v", l.L)
	}
	nl, err := l.C.NewLine(l.Args[0], l.OnExeced)
	if err != nil {
		return nil, err
	}
	v, err := nl.Exec(ctx, false)
	if err != nil {
		return nil, err
	}
	if !ValY(v) {
		return false, Err("line(%v) expected YES but NO", l.L)
	}
	return true, nil
}

func (l *Line) N(ctx *Ctx, left bool) (interface{}, error) {
	if len(l.Args) < 1 {
		return nil, Err("Usage:N $a,but:%v", l.L)
	}
	nl, err := l.C.NewLine(l.Args[0], l.OnExeced)
	if err != nil {
		return nil, err
	}
	v, err := nl.Exec(ctx, false)
	if err != nil {
		return nil, err
	}
	if ValY(v) {
		return false, Err("line(%v) expected NO but YES", l.L)
	}
	return true, nil
}
