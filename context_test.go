package igtest

import (
	"bufio"
	"fmt"
	"github.com/Centny/Cny4go/util"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"testing"
)

func TestBc(t *testing.T) {
	defer func() {
		C_BC = "bc"
	}()
	ctx := NewCtx(nil)
	ctx.Kvs["abc"] = util.Map{
		"ab": 123,
	}
	ctx.SET("ab", "11")
	ctx.SET("abdd", 11)
	fmt.Println(ctx.GET("ab"))
	ctx.Kvs["a"] = 11
	fmt.Println(ctx.Kvs.StrValP("/abc/ab"))
	fmt.Println(ctx.Compile("$a+$(/abc/ab)"))
	v, err := ctx.BC("$a+$(/abc/ab)")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if v != 134 {
		t.Error("not right")
		return
	}
	_, err = ctx.BC("$annn+$(/abc/ab)")
	if err == nil {
		t.Error("not error")
		return
	}
	C_BC = ",kk"
	_, err = ctx.BC("$a+$(/abc/ab)")
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err)
	C_BC = "echo"
	_, err = ctx.BC("$a+$(/abc/ab)")
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err)
	ctx.BC("")
	//
	fmt.Println(ctx.Join("$a", "$(/abc/ab)"))
	ctx.ShowLog = true
	ctx.log("kkkk")
}
func TestHr(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.FormValue("a"))
		fmt.Println(r.FormValue("b"))
		fmt.Println(r.Header.Get("kk"))
		fmt.Println(r.FormFile("file"))
		w.Write([]byte("OKK"))
	}))
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`
			{"abc":11}
			`))

	}))
	ctx := NewCtx(nil)
	ctx.HG(ts.URL, "a=1", "b=2", "^kk=11", "#code=aa", "#data=bb", "#err=cc")
	ctx.P("$aa", "$bb", "$cc")
	util.FWrite("/tmp/11122.txt", "testing")
	defer os.Remove("/tmp/11122.txt")
	ctx.HP(ts.URL, "a=2", "b=3", "^kk=kkk", "%file=/tmp/11122.txt", "#code=aa", "#data=bb", "#err=cc")
	ctx.P("$aa", "$bb", "$cc")
	ctx.HP("kkk", "a=2", "b=3", "^kk=kkk", "%file=/tmp/11122.txt", "#code=aa", "#data=bb", "#err=cc")
	ctx.HP(ts.URL, "a=2", "b=3", "^kk", "#code=aa", "#data=bb", "#err=cc")
	ctx.HP()
	ctx.HR("MM", "kkkk")
	ctx.HP(ts2.URL, "a=2", "b=3", "^kk", "#code=aa", "#data=bb", "#err=cc")
}
func TestEx(t *testing.T) {
	ctx := NewCtx(nil)
	ctx.EX("/bin/echo", `{\"abc\":123}`, "#aa", "#data=kk", "#err=ekk")
	ctx.P("$kk")
	ctx.EX("/bin/echo", `abbb`, "#aa", "#data=kk", "#err=ekk")
	ctx.EX("/bin/eco", `{\"abc\":123}`, "#aa", "#data=kk", "#err=ekk")
	NewCtx(ctx)
}
func TestRead(t *testing.T) {
	fmt.Println(os.Getwd())
	cmd := exec.Command("bc", "-l")
	in, err := cmd.StdinPipe()
	fmt.Println(err)
	stdout, err := cmd.StdoutPipe()
	fmt.Println(err)
	cmd.Start()
	in.Write([]byte("1+2\n"))
	in.Close()
	r := bufio.NewReader(stdout)
	line, _, err := r.ReadLine()
	fmt.Println(string(line), err)
}

func TestExecCmd(t *testing.T) {
	data, err := exec.Command("/bin/sh", "-c", "'/bin/echo 1+2 | /bin/echo'").Output()
	fmt.Println(string(data), err)
}

func TestUrlEncode(t *testing.T) {
	urls, _ := url.ParseQuery("ab=1")
	urls.Add("kk", "1111")
	fmt.Println(urls.Encode())
}
