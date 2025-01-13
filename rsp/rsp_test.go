package rsp_test

import (
	"fmt"
	"testing"

	"github.com/go-json-experiment/json"

	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/rsp"
)

func jsonString(v any) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}
func TestMsg(t *testing.T) {
	fmt.Println(jsonString(rsp.M("test")))
	fmt.Println(jsonString(rsp.OK()))
	fmt.Println(jsonString(rsp.M("aaaaaaaaaaaaaaaa")))
	fmt.Println(jsonString(rsp.OK()))
	fmt.Println(jsonString(rsp.E(code.RecordUnknown, "RecordUnknown")))
	fmt.Println(jsonString(rsp.OK()))
}
