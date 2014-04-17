package igtest

import (
	"bytes"
	"code.google.com/p/go.net/publicsuffix"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/util"
	"io/ioutil"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var v1_reg = regexp.MustCompile("\\$[a-zA-Z0-9]+")
var v2_reg = regexp.MustCompile("\\$\\([^\\)]+\\)")

type Ctx struct {
	Parent  *Ctx
	Kvs     util.Map
	H       util.HClient
	ShowLog bool
}

func (c *Ctx) log(f string, args ...interface{}) {
	if c.ShowLog {
		log.D(f, args...)
	}
}
func (c *Ctx) GET(path string) string {
	return c.Kvs.StrValP(c.Compile(path))
}

func (c *Ctx) SET(path string, val interface{}) error {
	if sval, ok := val.(string); ok {
		sval = c.Compile(sval)
		js, err := util.Json2Map(sval)
		if err == nil {
			return c.Kvs.SetValP(c.Compile(path), js)
		} else {
			return c.Kvs.SetValP(c.Compile(path), sval)
		}
	} else {
		return c.Kvs.SetValP(c.Compile(path), val)
	}
}

func (c *Ctx) Join(args ...string) string {
	var buf bytes.Buffer
	for _, arg := range args {
		// fmt.Println(arg)
		buf.WriteString(c.Compile(arg))
	}
	return buf.String()
}

func (c *Ctx) BC(exp string) (float64, error) {
	exp = strings.Trim(exp, "\t ")
	if len(exp) < 1 {
		return 0, errors.New("express is emtpy")
	}
	cexp := c.Compile(exp)
	c.log("BC %v", cexp)
	cmd := exec.Command(C_BC, "-l")
	w, _ := cmd.StdinPipe()
	r, _ := cmd.StdoutPipe()
	e, _ := cmd.StderrPipe()
	cmd.Start()
	w.Write(append([]byte(cexp), '\n'))
	w.Close()
	eres, _ := ioutil.ReadAll(e)
	ores, _ := ioutil.ReadAll(r)
	if len(eres) > 0 {
		return 0, errors.New(fmt.Sprintf("bc(%v) err(%v)", cexp, string(eres)))
	}
	if len(ores) < 1 {
		return 0, errors.New(fmt.Sprintf("bc(%v) err(not result)", cexp))
	}
	sres := strings.Trim(string(ores), "\n \t")
	fv, err := strconv.ParseFloat(sres, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("bc(%v) err(invalid result:%v)", cexp, sres))
	}
	return fv, err
}

func (c *Ctx) Compile(exp string) string {
	nbys := v1_reg.ReplaceAllFunc([]byte(exp), func(sby []byte) []byte {
		src := string(sby)
		src = strings.Trim(src, "$")
		return []byte(c.Kvs.StrVal(src))
	})
	pbys := v2_reg.ReplaceAllFunc(nbys, func(sby []byte) []byte {
		src := string(sby)
		src = strings.Trim(src, "$ ()")
		return []byte(c.Kvs.StrValP(src))
	})
	return string(pbys)
}
func (c *Ctx) HR(args ...string) (int, string, error) {
	if len(args) < 2 {
		return 0, "", errors.New(fmt.Sprintf("Usage:method url args"))
	}
	var turl = c.Compile(args[1])
	var method = c.Compile(args[0])
	fkey, fp := "", ""
	fields := map[string]string{}
	header := map[string]string{}
	exec_res := map[string]string{}
	for _, arg := range args[2:] {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) < 2 {
			continue
		}
		if strings.HasPrefix(kv[0], "^") {
			header[strings.TrimLeft(kv[0], "^")] = c.Compile(kv[1])
		} else if strings.HasPrefix(kv[0], "%") {
			fkey, fp = strings.TrimLeft(kv[0], "%"), c.Compile(kv[1])
		} else if strings.HasPrefix(kv[0], "#") {
			exec_res[strings.TrimLeft(kv[0], "#")] = c.Compile(kv[1])
		} else {
			fields[kv[0]] = c.Compile(kv[1])
		}
	}
	var code int
	var data string
	var err error
	if method == "POST" {
		c.log("POST(%v) fileds(%v) header(%v) fkey(%v) fpath(%v)", turl, fields, header, fkey, fp)
		code, data, err = c.H.HPostF_H(turl, fields, header, fkey, fp)
	} else if method == "GET" {
		vs := url.Values{}
		for k, v := range fields {
			vs.Add(k, v)
		}
		turl = fmt.Sprintf("%v?%v", turl, vs.Encode())
		c.log("GET(%v) header(%v)", turl, header)
		code, data, err = c.H.HGet_H(header, turl)
	} else {
		return 0, "", errors.New(fmt.Sprintf("unknow method(%v)", args[0]))
	}
	if k, ok := exec_res["code"]; ok {
		k = c.Compile(k)
		c.Kvs.SetVal(k, code)
	}
	if k, ok := exec_res["data"]; ok {
		k = c.Compile(k)
		js, err := util.Json2Map(data)
		if err == nil {
			c.Kvs.SetVal(k, js)
		} else {
			c.Kvs.SetVal(k, data)
		}
	}
	if k, ok := exec_res["err"]; ok && err != nil {
		c.Kvs.SetVal(k, err.Error())
	}
	return code, data, err
}

func (c *Ctx) HG(args ...string) (int, string, error) {
	nargs := make([]string, len(args)+1)
	nargs[0] = "GET"
	for i, a := range args {
		nargs[i+1] = a
	}
	return c.HR(nargs...)
}

func (c *Ctx) HP(args ...string) (int, string, error) {
	nargs := make([]string, len(args)+1)
	nargs[0] = "POST"
	for i, a := range args {
		nargs[i+1] = a
	}
	return c.HR(nargs...)
}
func (c *Ctx) P(args ...string) {
	nargs := make([]interface{}, len(args))
	for i, v := range args {
		nargs[i] = c.Compile(v)
	}
	fmt.Println(nargs...)
}

func (c *Ctx) EX(args ...string) (string, error) {
	nargs := []string{}
	exec_res := map[string]string{}
	for _, a := range args {
		if strings.HasPrefix(a, "#") {
			kv := strings.SplitN(a, "=", 2)
			if len(kv) < 2 {
				continue
			}
			exec_res[strings.Trim(kv[0], "#")] = kv[1]
		} else {
			nargs = append(nargs, c.Compile(a))
		}
	}
	bys, err := exec.Command(C_SH, "-c", strings.Join(nargs, " ")).Output()
	data := string(bys)
	if k, ok := exec_res["data"]; ok {
		js, err := util.Json2Map(data)

		if err == nil {
			c.Kvs.SetVal(k, js)
		} else {
			c.Kvs.SetVal(k, strings.Trim(data, "\t \n"))
		}
	}
	if k, ok := exec_res["err"]; ok && err != nil {
		c.Kvs.SetVal(k, err.Error())
	}
	return data, err
}

func (c *Ctx) W(args ...string) error {
	if len(args) < 1 {
		return Err("write file error:%v", "file path not set")
	}
	if len(args) < 2 {
		return Err("write file error:%v", "content not set")
	}
	fpath := c.Compile(args[0])
	f, err := os.OpenFile(args[0], os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	var val interface{} = nil
	if e_reg_v.MatchString(args[1]) {
		val, err = c.Kvs.ValP(strings.Trim(args[1], " \t$()"))
		if err != nil {
			return err
		}
	} else {
		val = args[1]
	}
	var l int = 0
	if mv := util.MapVal(val); mv != nil {
		jbys, _ := json.Marshal(mv)
		l, err = f.Write(jbys)
	} else {
		l, err = f.WriteString(util.StrVal(val))
	}
	if err == nil {
		c.log("wite %d data to file:%v", l, fpath)
	}
	return err
}

func (c *Ctx) R(args ...string) (interface{}, error) {
	if len(args) < 1 {
		return nil, Err("write file error:%v", "file path not set")
	}

	fpath := c.Compile(args[0])
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bys, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var val interface{} = nil
	val, err = util.Json2Map(string(bys))
	if err != nil {
		val = string(bys)
	}
	if len(args) > 1 && strings.HasPrefix(args[1], "#") {
		kv := strings.SplitN(args[1], "=", 2)
		if len(kv) == 2 && kv[0] == "#data" {
			c.log("setting data to %v", kv[1])
			c.SET(kv[1], val)
		}
	}
	c.log("read %d data from file:%v", len(bys), fpath)
	return val, nil
}

func (c *Ctx) D(args ...string) error {
	if len(args) < 1 {
		return Err("delete file error:%v", "file path not set")
	}
	fpath := c.Compile(args[0])
	err := os.RemoveAll(fpath)
	c.log("delete file:%v", fpath)
	return err
}
func NewCtx(p *Ctx) *Ctx {
	ctx := &Ctx{}
	ctx.Kvs = util.Map{}
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err == nil {
		ctx.H.Jar = jar
	}
	ctx.Parent = p
	if ctx.Parent != nil {
		ctx.Kvs["@PRE"] = p.Kvs
	}
	return ctx
}
