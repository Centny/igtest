package igtest

import (
	"errors"
	"fmt"
	"github.com/Centny/Cny4go/util"
	"math"
)

type OnExecedFunc func(*Line, bool, interface{}, error) (interface{}, error)

func Err(f string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(f, args...))
}

func ValY(v interface{}) bool {
	if v == nil {
		return false
	}
	iv := util.IntVal(v)
	if iv != math.MaxInt64 {
		return iv != 0
	}
	sv := util.StrVal(v)
	if len(sv) > 0 {
		return true
	}
	return false
}

func YesOnExeced(l *Line, left bool, v interface{}, e error) (interface{}, error) {
	if e != nil {
		return nil, e
	}
	if !left {
		return v, e
	}
	if ValY(v) {
		return v, nil
	} else {
		return v, Err("Excepted YES but NO(%v) in line(%v)", v, l.L)
	}
}

func ExpYesOnExeced(l *Line, left bool, v interface{}, e error) (interface{}, error) {
	if e != nil {
		return nil, e
	}
	if !left {
		return v, e
	}
	if l == nil {
		return nil, Err("Line is nil")
	}
	if l.T != "EXP" {
		return v, e
	}
	if ValY(v) {
		return v, nil
	} else {
		return v, Err("Excepted YES but NO(%v) in line(%v)", v, l.L)
	}
}

func Escape(bys []byte) []byte {
	return escape.ReplaceAllFunc(bys, func(sby []byte) []byte {
		return sby[1:]
	})
}
