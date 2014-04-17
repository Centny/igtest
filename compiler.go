package igtest

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/Centny/Cny4go/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Compiler struct {
	cidx  int
	Lines []string
	Cwd   string
}

func (c *Compiler) addl(line string) {
	c.Lines = append(c.Lines, line)
}

func (c *Compiler) Load(spath string) error {
	c.Cwd = filepath.Dir(spath)
	return c.load(spath)
}
func (c *Compiler) load(spath string) error {
	f, err := os.Open(spath)
	if err != nil {
		return err
	}
	sbys, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()
	return c.PreCompile(sbys)
}

//precompile will load all script and execute the INC cmd.
func (c *Compiler) PreCompile(sbys []byte) error {
	c.Lines = []string{}
	sbys = b_reg_cmt.ReplaceAll(sbys, []byte(""))
	sbys = l_reg_cmt.ReplaceAll(sbys, []byte("\n"))
	sbys = Escape(sbys)
	R := bufio.NewReader(bytes.NewReader(sbys))
	for {
		bys, err := util.ReadLine(R, 1024000, false)
		if err != nil {
			break
		}
		line := string(bys)
		line = strings.Trim(line, "\t \n")
		if len(line) < 1 {
			continue
		}
		if c_reg_INC.MatchString(line) {
			args := empty.Split(line, -1)
			if len(args) < 2 {
				return errors.New(fmt.Sprintf(
					"INC(%v) error:%v", line, "argument is empty"))
			}
			tpath := ""
			if filepath.IsAbs(args[1]) {
				tpath = args[1]
			} else {
				tpath = filepath.Join(c.Cwd, args[1])
			}
			err := c.load(tpath)
			if err != nil {
				return errors.New(fmt.Sprintf("INC(%v) load error:%v", line, err.Error()))
			}
			continue
		}
		c.addl(line)
	}
	return nil
}

func (c *Compiler) CompileAndExec(ctx *Ctx, onexeced OnExecedFunc) error {
	ls, err := c.Compile(onexeced)
	if err != nil {
		return err
	}
	for _, l := range ls {
		_, err := l.Exec(ctx, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) Compile(onexeced OnExecedFunc) ([]*Line, error) {
	lcount := len(c.Lines)
	clines := []*Line{}
	for_l := []*Line{}
	for_i := -1
	for cidx := 0; cidx < lcount; cidx++ {
		line := c.Lines[cidx]
		if c_reg_FOR.MatchString(line) {
			args := empty.Split(line, -1)
			if len(args) < 4 {
				return nil, Err("invalid FOR:%v", line)
			}
			if c_reg_ROF.MatchString(line) { //if FOR ROF in the same line
				args := empty.Split(line, -1)
				if len(args) < 6 {
					return nil, Err("invalid FOR:%v", line)
				}
				l, err := c.NewLine(strings.Join(args[4:len(args)-1], " "), onexeced)
				if err != nil {
					return nil, err
				}
				nl := c.cline("FOR", line, args[1:4], onexeced)
				nl.addl(l)
				clines = append(clines, nl)
			} else {
				l := c.cline("FOR", line, args[1:4], onexeced)
				for_l = append(for_l, l)
				for_i++
				if for_i == 0 {
					clines = append(clines, l)
				} else {
					for_l[for_i-1].addl(l)
				}
			}
		} else if c_reg_ROF.MatchString(line) {
			cline := c_reg_e_ROF.ReplaceAllString(line, "")
			cline = strings.Trim(cline, "\t ")
			if len(cline) > 0 { //check if have command in ROF
				l, err := c.NewLine(cline, onexeced)
				if err != nil {
					return nil, err
				}
				for_l[for_i].addl(l)
			}
			for_l = for_l[:len(for_l)-1]
			for_i--
		} else if for_i > -1 {
			l, err := c.NewLine(line, onexeced)
			if err != nil {
				return nil, err
			}
			for_l[for_i].addl(l)
		} else {
			l, err := c.NewLine(line, onexeced)
			if err != nil {
				return nil, err
			}
			clines = append(clines, l)
		}
	}
	return clines, nil
}

func (c *Compiler) NewLine(line string, onexeced OnExecedFunc) (*Line, error) {
	line = strings.Trim(line, " \t")
	var T = ""
	var args []string = []string{}
	if c_reg_Assign1.MatchString(line) {
		args = strings.SplitN(line, "=", 2)
		T = "ASSIGN"
	} else if c_reg_Assign2.MatchString(line) {
		args = strings.SplitN(line, "=", 2)
		T = "ASSIGN"
	} else if c_reg_BC.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("BC error:%v", "argument is empty")
		}
		T = "BC"
		args = args[1:]
	} else if c_reg_SET.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 3 {
			return nil, Err("SET error:%v", "argument is less 2")
		}
		T = "SET"
		args = args[1:]
	} else if c_reg_GET.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("GET error:%v", "argument is empty")
		}
		T = "GET"
		args = args[1:]
	} else if c_reg_HR.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 3 {
			return nil, Err("HR error:%v", "argument is less 2")
		}
		T = "HR"
		args = args[1:]
	} else if c_reg_HG.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("HG error:%v", "argument is empty")
		}
		T = "HG"
		args = args[1:]
	} else if c_reg_HP.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("HP error:%v", "argument is empty")
		}
		T = "HP"
		args = args[1:]
	} else if c_reg_EX.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("EX error:%v", "argument is empty")
		}
		T = "EX"
		args = args[1:]
	} else if c_reg_P.MatchString(line) {
		args = empty.Split(line, -1)
		T = "P"
		args = args[1:]
	} else if c_reg_SUB.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("SUB error:%v", "argument is empty")
		}
		T = "SUB"
		args = args[1:]
	} else if c_reg_EXP.MatchString(line) {
		T = "EXP"
		args = []string{line}
	} else if c_reg_Y.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("Y error:%v", "argument is empty")
		}
		T = "Y"
		args = args[1:]
	} else if c_reg_N.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("N error:%v", "argument is empty")
		}
		T = "N"
		args = args[1:]
	} else {
		// panic(11)
		return nil, Err("invalid line(%v)", line)
	}
	return c.cline(T, line, args, onexeced), nil
}

func (c *Compiler) cline(t, l string, args []string, onexeced OnExecedFunc) *Line {
	return &Line{
		T:        t,
		L:        l,
		Args:     args,
		Lines:    []*Line{},
		OnExeced: onexeced,
		C:        c,
	}
}

func (c *Compiler) ShowLine() {
	for _, l := range c.Lines {
		fmt.Println(l)
	}
}

func ShowLine(ls []*Line, inner int) {
	pre := ""
	for i := 0; i < inner; i++ {
		pre = fmt.Sprintf("%v ", pre)
	}
	for idx, l := range ls {
		fmt.Println(fmt.Sprintf("%v%v\t\t%v%v", pre, l.T, pre, idx))
		if len(l.Lines) > 0 {
			ShowLine(l.Lines, inner+1)
		}
	}
}
