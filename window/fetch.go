package window

import (
	"bytes"
	"encoding/json"
	"github.com/dop251/goja"
	"io"
	"net/http"
	"strings"
	"time"
)

type initModel struct {
	Method         string                 `json:"method"`      // *GET, POST, PUT, DELETE, etc.
	Mode           string                 `json:"mode"`        // no-cors, *cors, same-origin
	Cache          string                 `json:"cache"`       // *default, no-cache, reload, force-cache, only-if-cached
	Credentials    string                 `json:"credentials"` // include, *same-origin, omit
	Headers        map[string]string      `json:"headers"`
	Redirect       string                 `json:"redirect"`       // manual, *follow, error
	ReferrerPolicy string                 `json:"referrerPolicy"` // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
	Body           map[string]interface{} `json:"body"`
	// more for extend
}

type fetchResponse struct {
	Headers    map[string][]string `json:"headers"`
	Ok         bool                `json:"ok"`
	Status     int                 `json:"status"`
	StatusText string              `json:"statusText"`
	Body       []byte              `json:"body"`
}

type fetchMode struct {
	r *goja.Runtime
}

func (f *fetchMode) Fetch(call goja.FunctionCall) goja.Value {
	arg0 := call.Argument(0)
	arg1 := call.Argument(1)
	var init initModel
	if !goja.IsUndefined(arg1) {
		outdata := arg1.Export()
		data, _ := json.Marshal(outdata)
		_ = json.Unmarshal(data, &init)
	}
	out := f.fetch(arg0.String(), &init)
	oo := f.r.ToValue(out)
	return oo
}

func (f *fetchMode) fetch(reqUrl string, init *initModel) *fetchResponse {
	method := "GET"
	if init != nil && init.Method != "" {
		method = init.Method
	}
	var body bytes.Buffer
	if init != nil && len(init.Body) > 0 {
		_ = json.NewEncoder(&body).Encode(init.Body)
	}
	req, err := http.NewRequest(strings.ToUpper(method), reqUrl, &body)
	if err != nil {
		return nil
	}
	// make header
	f.headers(init, req)
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer func() { _ = res.Body.Close() }()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	return &fetchResponse{
		Headers:    res.Header,
		Ok:         res.StatusCode >= 200 && res.StatusCode < 300,
		Status:     res.StatusCode,
		StatusText: res.Status,
		Body:       buf,
	}
}

func (f *fetchMode) headers(init *initModel, req *http.Request) {
	if init == nil || len(init.Headers) == 0 {
		return
	}
	for h, v := range init.Headers {
		req.Header.Set(h, v)
	}
}

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        1,
		MaxIdleConnsPerHost: 1,
		MaxConnsPerHost:     5,
		IdleConnTimeout:     time.Minute * 1,
	},
	Timeout: time.Second * 30,
}
