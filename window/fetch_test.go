package window

import (
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"testing"
)

func TestFetchNoBody(t *testing.T) {
	vm := goja.New()
	new(require.Registry).Enable(vm)
	Enable(vm)
	if c := vm.Get("window"); c == nil {
		t.Fatal("window not found")
	}
	script := `data = window.fetch("http://api.example.com");`
	v, err := vm.RunString(script)
	if err != nil {
		t.Fatal("Failed to process url script.", err)
	}
	res, ok := v.Export().(*fetchResponse)
	if !ok {
		t.Logf("res not fetchResponse:%t", ok)
	} else {
		t.Logf("res:%v", res)
		t.Logf("res:%s", res.Body)
	}
}

func TestFetchWithBody(t *testing.T) {
	// document: https://developer.mozilla.org/zh-CN/docs/Web/API/Fetch_API/Using_Fetch
	vm := goja.New()
	new(require.Registry).Enable(vm)
	Enable(vm)
	if c := vm.Get("window"); c == nil {
		t.Fatal("window not found")
	}
	body := map[string]string{}
	body["name"] = "abc.test"
	body["hence"] = "20240123"
	bodyBytes, _ := json.Marshal(body)
	script := fmt.Sprintf(`data = window.fetch("http://api.example.com", 
		{
		"method":"POST", 
		"body":%s, 
		"headers": {"Content-Type":"application/json"}});`, bodyBytes)
	v, err := vm.RunString(script)
	if err != nil {
		t.Fatal("Failed to process url script.", err)
	}
	res, ok := v.Export().(*fetchResponse)
	if !ok {
		t.Logf("res not fetchResponse:%t", ok)
	} else {
		t.Logf("res:%v", res)
		t.Logf("res:%s", res.Body)
	}
}
