package igtest

import (
	"os"
	"testing"
)

func TestExec(t *testing.T) {
	jmk := NewJsonMarker("/tmp/exec.json")
	ctx := NewCtx(nil)
	ctx.ShowLog = true
	ctx.Mark = jmk
	err := ExecCtx("test/exec.ig", ctx)
	if err != nil {
		t.Error(err.Error())
		return
	}
	os.Remove("/tmp/exec.json")
	jmk.Store()
	Exec("test/inc.ig")
	ExecCtx("kkk", NewCtx(nil))
}

func TestRun(t *testing.T) {
	Run([]string{"-h"})
	Run([]string{"-l", "-m", "YOE", "test/sub.ig"})
	Run([]string{"-l", "-m", "YOE", "-R", "JSON", "-r", "sub.json", "test/sub.ig"})
	Run([]string{"-l", "-m", "YOE", "-R", "JSON", "-r", "/t/k/sub.json", "test/sub.ig"})
	Run([]string{"-l", "-m", "YOE", "-R", "JSO", "-r", "sub.json", "test/sub.ig"})
	Run([]string{"-l", "-m", "YOE", "test/sub.ig", "ab=1111"})
	Run([]string{"-l", "-m", "EYOE", "test/sub.ig"})
	Run([]string{"-l", "-m", "NONE", "test/sub.ig"})
	Run([]string{"-l", "-m", "YOE"})
	Run([]string{"-l", "-m", "NN", "test/sub.ig"})
	Run([]string{"test/subb.ig"})
	Run([]string{"test/err.ig"})
}

func TestJme(t *testing.T) {
	err := NewJME("/tmp/tt.json", "test/exec.ig").Exec()
	if err != nil {
		t.Error(err.Error())
		return
	}
	NewJME("/tmp/tt.json", "test/ex.ig").Exec()
	NewJME("/tmp/tt.json", "test/err.ig").Exec()
	NewJME("/kk/k/tt.json", "test/exec.ig").Exec()
	jme := NewJME("/tmp/tt.json", "test/exec.ig")
	err = jme.E("SUB test/sub.ig")
	if err != nil {
		t.Error(err.Error())
		return
	}
	jme.E("SUBB test/sub.ig")
	jme.E("SUB test/subb.ig")
}
