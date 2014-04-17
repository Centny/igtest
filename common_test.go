package igtest

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"
)

func TestMatch(t *testing.T) {
	val := b_reg_cmt.ReplaceAll(
		[]byte(`

/*

*/
kkdkfs

/* 
kkk*/`), []byte(""))
	fmt.Println(string(val))
	fmt.Println(regexp.MatchString("^INC[\\ \\t]+", "INC abcc"))
	fmt.Println(regexp.MatchString("^INC[\\ \\t]+", "INCabcc"))
	fmt.Println(regexp.MustCompile("^INC[\\ \\t]+").ReplaceAllString("INC abcc", "ss"))
}

func TestValY(t *testing.T) {
	ValY(nil)
	ValY(1)
	ValY(0)
	ValY("ss")
	ValY("")
}

func TestYesOnExeced(t *testing.T) {
	YesOnExeced(&Line{T: "EXP"}, true, nil, Err("kk"))
	YesOnExeced(&Line{T: "EXP"}, true, 1, nil)
	YesOnExeced(&Line{T: "EXP"}, true, 0, nil)
	YesOnExeced(&Line{T: "EXP"}, false, 0, nil)
}

func TestExpYesOnExeced(t *testing.T) {
	ExpYesOnExeced(nil, true, nil, nil)
	ExpYesOnExeced(&Line{T: "EXP"}, true, nil, Err("kk"))
	ExpYesOnExeced(&Line{T: "EXP"}, true, nil, nil)
	ExpYesOnExeced(&Line{T: "EXP"}, true, 1, nil)
	ExpYesOnExeced(&Line{T: "EXP"}, true, 0, nil)
	ExpYesOnExeced(&Line{T: "ABC"}, true, 1, nil)
	ExpYesOnExeced(&Line{T: "EXP"}, false, 1, nil)
}

func TestFpath(t *testing.T) {
	fmt.Println(filepath.Join("a", "a/"))
	fmt.Println(filepath.Dir("test/inc.ig"))
}

func TestComment(t *testing.T) {
	fmt.Println(regexp.MustCompile("(?U)//.*\\n").ReplaceAllString(`
		//kkk
		sdfs
		kfdf //kk
		`, "\n"))
}
