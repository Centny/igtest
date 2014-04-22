package igtest

import (
	"encoding/json"
	"errors"
	"github.com/Centny/Cny4go/util"
	"io"
	"os"
	"strings"
)

type JsonMarker struct {
	pre  util.Map
	ms   []util.Map
	Path string
}

func (j *JsonMarker) new(l *Line, ctx *Ctx) util.Map {
	name := ""
	if len(l.Args) > 0 {
		name = ctx.Compile(l.Args[0])
	}
	f := l.C.F
	return util.Map{
		"file": f,
		"line": l.Num,
		"type": l.T,
		"name": name,
		"desc": ctx.Compile(strings.Join(l.Args[1:], ",")),
	}
}
func (j *JsonMarker) add(m util.Map) {
	j.ms = append(j.ms, m)
	if j.pre != nil {
		j.pre["subs"] = j.ms
	}
}
func (j *JsonMarker) Add(l *Line, ctx *Ctx) {
	j.add(j.new(l, ctx))
}
func (j *JsonMarker) Sub(l *Line, ctx *Ctx) Marker {
	sub := NewJsonMarker(j.Path)
	sl := j.new(l, ctx)
	sub.pre = sl
	j.add(sl)
	return sub
}

func (j *JsonMarker) StoreW(w io.Writer) (int, error) {
	bys, err := json.Marshal(j.ms)
	if err != nil {
		return 0, err
	}
	return w.Write(bys)
}

func (j *JsonMarker) StoreF(fname string) error {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = j.StoreW(f)
	return err
}

func (j *JsonMarker) Store() error {
	if len(j.Path) < 1 {
		return errors.New("Path not set")
	}
	os.Remove(j.Path)
	return j.StoreF(j.Path)
}

func NewJsonMarker(path string) *JsonMarker {
	jm := &JsonMarker{}
	jm.pre = nil
	jm.ms = []util.Map{}
	jm.Path = path
	return jm
}
