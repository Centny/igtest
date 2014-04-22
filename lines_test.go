package igtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSub(t *testing.T) {
	c := &Compiler{}
	c.Load("test/sub.ig")
	l := Line{}
	l.C = c
	l.T = "SUB"
	_, err := l.Sub(NewCtx(nil), false)
	if err == nil {
		t.Error("not error")
	}
	l.Args = []string{"kkkk.ig"}
	_, err = l.Sub(NewCtx(nil), false)
	if err == nil {
		t.Error("not error")
	}
	//
	l = Line{}
	l.C = c
	l.T = "EXP"
	_, err = l.Exp(NewCtx(nil), false)
	if err == nil {
		t.Error("not error")
	}
	l.Args = []string{"@[exx]"}
	_, err = l.Exp(NewCtx(nil), false)
	if err == nil {
		t.Error("not error")
	}
	l.Args = []string{"$[exx]"}
	_, err = l.Exp(NewCtx(nil), false)
	if err == nil {
		t.Error("not error")
	}
	l.Args = []string{"${exx}"}
	_, err = l.Exp(NewCtx(nil), false)
	if err == nil {
		t.Error("not error")
	}
	Exec("test/for.ig")
	//
	l = Line{}
	l.C = c
	l.T = "BC"
	l.exec(NewCtx(nil), false)
	l.T = "SET"
	l.exec(NewCtx(nil), false)
	l.T = "GET"
	l.exec(NewCtx(nil), false)
	l.T = "HR"
	l.exec(NewCtx(nil), false)
	l.T = "HP"
	l.exec(NewCtx(nil), false)
	l.T = "HG"
	l.exec(NewCtx(nil), false)
	l.T = "EX"
	l.exec(NewCtx(nil), false)
	l.T = "BC"
	l.exec(NewCtx(nil), false)
	l.T = "XX"
	l.exec(NewCtx(nil), false)
	l.T = "R"
	l.exec(NewCtx(nil), false)
	l.T = "D"
	l.exec(NewCtx(nil), false)
	l.T = "W"
	l.exec(NewCtx(nil), false)
	l.T = "M"
	l.exec(NewCtx(nil), false)
	//
	l = Line{}
	l.C = c
	l.Assign(NewCtx(nil), false)
	l.Args = []string{"$a", "@[kkkkkk]"}
	l.Assign(NewCtx(nil), false)
	l.Args = []string{"$a", "@[sssss"}
	l.Assign(NewCtx(nil), false)
	l.Args = []string{"$a", "$a<b"}
	l.Assign(NewCtx(nil), false)
	//
	l = Line{}
	l.C = c
	l.For(NewCtx(nil), false)
	l.Args = []string{"aaa", "KK", "kkk", "kkk"}
	l.For(NewCtx(nil), false)
	//
	l = Line{}
	l.C = c
	l.Y(NewCtx(nil), false)
	l.N(NewCtx(nil), false)
	l.Args = []string{"aaaa"}
	l.Y(NewCtx(nil), false)
	l.N(NewCtx(nil), false)
	l.Args = []string{"$a<b"}
	l.Y(NewCtx(nil), false)
	l.N(NewCtx(nil), false)
	l.Args = []string{"EX kkkk"}
	l.Y(NewCtx(nil), false)
	l.N(NewCtx(nil), false)
	l.Args = []string{"$a"}
	l.Y(NewCtx(nil), false)
	ctx := NewCtx(nil)
	ctx.SET("a", 111)
	l.N(ctx, false)
}

func TestASubNet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("P $bb\n"))
		fmt.Println("lllllll")
	}))
	c := &Compiler{}
	l := Line{}
	l.C = c
	l.T = "SUB"
	l.Args = []string{ts.URL, "bb=111111"}
	_, err := l.Sub(NewCtx(nil), true)
	if err != nil {
		t.Error(err.Error())
	}
	l.Args = []string{"http://kkk", "bb=111111"}
	_, err = l.Sub(NewCtx(nil), true)
	if err == nil {
		t.Error("not error")
	}
}
