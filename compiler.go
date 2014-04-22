package igtest

import (
	"bufio"
	"bytes"
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
	Lnum  []int
	Cwd   string
	F     string
}

func (c *Compiler) addl(line string, num int) {
	c.Lines = append(c.Lines, line)
	c.Lnum = append(c.Lnum, num)
}

func (c *Compiler) Load(spath string) error {
	c.F = spath
	cwd, _ := os.Getwd()
	c.Cwd = filepath.Join(cwd, filepath.Dir(spath))
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
	// sbys = b_reg_cmt.ReplaceAll(sbys, []byte(""))
	// sbys = l_reg_cmt.ReplaceAll(sbys, []byte("\n"))
	// sbys = Escape(sbys)
	var iscmt bool = false
	R := bufio.NewReader(bytes.NewReader(sbys))
	for i := 1; true; i++ {
		bys, err := util.ReadLine(R, 1024000, false)
		if err != nil {
			break
		}
		if b_reg_cmt_beg.Match(bys) {
			iscmt = true
		}
		if iscmt && b_reg_cmt_end.Match(bys) {
			iscmt = false
			continue
		}
		if iscmt {
			continue
		}
		bys = l_reg_cmt.ReplaceAll(bys, []byte(""))
		line := string(Escape(bys))
		line = strings.Trim(line, "\t \n")
		if len(line) < 1 {
			continue
		}
		// if c_reg_INC.MatchString(line) {
		// 	args := empty.Split(line, -1)
		// 	if len(args) < 2 {
		// 		return errors.New(fmt.Sprintf(
		// 			"INC(%v) error:%v", line, "argument is empty"))
		// 	}
		// 	tpath := ""
		// 	if filepath.IsAbs(args[1]) {
		// 		tpath = args[1]
		// 	} else {
		// 		tpath = filepath.Join(c.Cwd, args[1])
		// 	}
		// 	err := c.load(tpath)
		// 	if err != nil {
		// 		return errors.New(fmt.Sprintf("INC(%v) load error:%v", line, err.Error()))
		// 	}
		// 	continue
		// }
		c.addl(line, i)
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
		num := c.Lnum[cidx]
		if c_reg_FOR.MatchString(line) {
			args := empty.Split(line, -1)
			if len(args) < 4 {
				return nil, Err("invalid FOR in %v:%v %v", c.F, num, line)
			}
			if c_reg_ROF.MatchString(line) { //if FOR ROF in the same line
				args := empty.Split(line, -1)
				if len(args) < 6 {
					return nil, Err("invalid FOR in %v:%v %v", c.F, num, line)
				}
				l, err := c.NewLine(strings.Join(args[4:len(args)-1], " "), num, onexeced)
				if err != nil {
					return nil, err
				}
				nl := c.cline("FOR", line, num, args[1:4], onexeced)
				nl.addl(l)
				clines = append(clines, nl)
			} else {
				l := c.cline("FOR", line, num, args[1:4], onexeced)
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
				l, err := c.NewLine(cline, num, onexeced)
				if err != nil {
					return nil, err
				}
				for_l[for_i].addl(l)
			}
			for_l = for_l[:len(for_l)-1]
			for_i--
		} else if for_i > -1 {
			l, err := c.NewLine(line, num, onexeced)
			if err != nil {
				return nil, err
			}
			for_l[for_i].addl(l)
		} else {
			l, err := c.NewLine(line, num, onexeced)
			if err != nil {
				return nil, err
			}
			clines = append(clines, l)
		}
	}
	return clines, nil
}

func (c *Compiler) NewLine(line string, num int, onexeced OnExecedFunc) (*Line, error) {
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
			return nil, Err("BC error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "BC"
		args = args[1:]
	} else if c_reg_SET.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 3 {
			return nil, Err("SET error:%v in %v:%v", "argument is less 2", c.F, num)
		}
		T = "SET"
		args = args[1:]
	} else if c_reg_GET.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("GET error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "GET"
		args = args[1:]
	} else if c_reg_HR.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 3 {
			return nil, Err("HR error:%v in %v:%v", "argument is less 2", c.F, num)
		}
		T = "HR"
		args = args[1:]
	} else if c_reg_HG.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("HG error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "HG"
		args = args[1:]
	} else if c_reg_HP.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("HP error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "HP"
		args = args[1:]
	} else if c_reg_EX.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("EX error:%v in %v:%v", "argument is empty", c.F, num)
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
			return nil, Err("SUB error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "SUB"
		args = args[1:]
	} else if c_reg_EXP.MatchString(line) {
		T = "EXP"
		args = []string{line}
	} else if c_reg_Y.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("Y error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "Y"
		args = args[1:]
	} else if c_reg_N.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("N error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "N"
		args = args[1:]
	} else if c_reg_R.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("R error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "R"
		args = args[1:]
	} else if c_reg_W.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 3 {
			return nil, Err("W error:%v in %v:%v", "argument is less 2", c.F, num)
		}
		T = "W"
		args = args[1:]
	} else if c_reg_D.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("D error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "D"
		args = args[1:]
	} else if c_reg_M.MatchString(line) {
		args = empty.Split(line, -1)
		if len(args) < 2 {
			return nil, Err("M error:%v in %v:%v", "argument is empty", c.F, num)
		}
		T = "M"
		args = args[1:]
	} else {
		return nil, Err("invalid line(%v) in %v:%v", line, c.F, num)
	}
	return c.cline(T, line, num, args, onexeced), nil
}

func (c *Compiler) cline(t, l string, num int, args []string, onexeced OnExecedFunc) *Line {
	return &Line{
		T:        t,
		L:        l,
		Num:      num,
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

func NewCompiler() *Compiler {
	return &Compiler{}
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
