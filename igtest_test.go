package igtest

import (
	"testing"
)

func TestExec(t *testing.T) {
	// c := Compiler{}
	// c.Load("test/tcmt.ig")
	// c.ShowLine()
	// ls, _ := c.Compile(nil)
	// ShowLine(ls, 0)
	err := Exec("test/exec.ig")
	if err != nil {
		t.Error(err.Error())
	}

	ExecCtx("kkk", NewCtx(nil))
}
