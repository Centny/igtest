package igtest

import (
	"encoding/json"
	"fmt"
	"github.com/Centny/Cny4go/util"
	"math"
	"testing"
)

func TestAMap(t *testing.T) {
	bys, _ := json.Marshal(map[string]interface{}{
		"2-aaaa": 222,
		"1-aaa":  1111,
		"10-aaa": 1111,
	})
	fmt.Println(string(bys))
}

func TestMark(t *testing.T) {
	jm := NewJsonMarker("/tmp/t.json")
	jm.ms = []util.Map{
		util.Map{
			"kkk": math.NaN(),
		},
	}
	err := jm.Store()
	if err == nil {
		t.Error("not error")
		return
	}
	jm.StoreF("/jjjjj/kkss")
	NewJsonMarker("").Store()
}
