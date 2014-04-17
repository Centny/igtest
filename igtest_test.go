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

func TestRun(t *testing.T) {
	Run([]string{"-l", "-m", "YOE", "test/sub.ig"})
	Run([]string{"-l", "-m", "EYOE", "test/sub.ig"})
	Run([]string{"-l", "-m", "NONE", "test/sub.ig"})
	Run([]string{"-l", "-m", "YOE"})
	Run([]string{"-l", "-m", "NN", "test/sub.ig"})
	Run([]string{"test/subb.ig"})
	Run([]string{"test/err.ig"})
}
