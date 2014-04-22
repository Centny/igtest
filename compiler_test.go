package igtest

import (
	"fmt"
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	c := Compiler{}
	err := c.Load("test/compiler.ig")
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.ShowLine()
	ls, err := c.Compile(nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	ShowLine(ls, 0)
	//
	c.NewLine("BC ", 1, nil)
	c.NewLine("SET ", 1, nil)
	c.NewLine("GET ", 1, nil)
	c.NewLine("HR ", 1, nil)
	c.NewLine("HP ", 1, nil)
	c.NewLine("HG ", 1, nil)
	c.NewLine("EX ", 1, nil)
	c.NewLine("SUB ", 1, nil)
	c.NewLine("HRR ", 1, nil)
	c.NewLine("Y", 1, nil)
	c.NewLine("N", 1, nil)
	c.NewLine("R", 1, nil)
	c.NewLine("W", 1, nil)
	c.NewLine("D", 1, nil)
	c.NewLine("M", 1, nil)
	c = Compiler{}
	wd, _ := os.Getwd()
	err = c.PreCompile([]byte(fmt.Sprintf("INC %v/%v", wd, "test/sub.ig")))
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestLoad(t *testing.T) {
	r := Compiler{}
	err := r.Load("test/t.ig")
	if err != nil {
		t.Error(err.Error())
		return
	}
	// err = r.Load("test/t2.ig")
	// if err == nil {
	// 	t.Error("not error")
	// 	return
	// }
	err = r.Load("test/t2.igg")
	if err == nil {
		t.Error("not error")
		return
	}
	os.Mkdir("/tmp/jjj11", 0644)
	defer os.Remove("/tmp/jjj11")
	err = r.Load("/tmp/jjj11")
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(r.Lines)
	r.PreCompile([]byte("INC\n"))
}

func TestCompile(t *testing.T) {
	c := Compiler{}
	c.PreCompile([]byte(`
		FOR kk
		`))
	_, err := c.Compile(nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		FOR kk IN kkk ROF
		`))
	_, err = c.Compile(nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		FOR kk IN kkk kkk ROF
		`))
	_, err = c.Compile(nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		FOR kk IN kkk 
		kkk ROF
		`))
	_, err = c.Compile(nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		kkk
		`))
	_, err = c.Compile(nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		EX /kkkkkk kk
		`))
	err = c.CompileAndExec(NewCtx(nil), nil)
	if err == nil {
		t.Error("not error")
		return
	}

	c.PreCompile([]byte(`
		FOR kk IN kkk 
		EX kkkkkk
		ROF
		`))
	err = c.CompileAndExec(NewCtx(nil), nil)
	if err == nil {
		t.Error("not error")
		return
	}

	c.PreCompile([]byte(`
		FOR kk R 1
		EX kkkkkk
		ROF
		`))
	err = c.CompileAndExec(NewCtx(nil), nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		FOR kk R a~2
		EX kkkkkk
		ROF
		`))
	err = c.CompileAndExec(NewCtx(nil), nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		FOR kk R 1~a
		EX kkkkkk
		ROF
		`))
	err = c.CompileAndExec(NewCtx(nil), nil)
	if err == nil {
		t.Error("not error")
		return
	}
	c.PreCompile([]byte(`
		FOR kk R 1~5
		EX kkkkkk
		ROF
		`))
	err = c.CompileAndExec(NewCtx(nil), nil)
	if err == nil {
		t.Error("not error")
		return
	}
	//
}
