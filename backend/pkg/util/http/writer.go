package util_http

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Writer struct {
	baseUrl           string
	beforeRequestFunc func(req *fasthttp.Request) error
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) SetBaseUrl(url string) *Writer {
	w.baseUrl = url
	return w
}

func (w *Writer) OnBeforeRequest(beforeFunc func(req *fasthttp.Request) error) *Writer {
	w.beforeRequestFunc = beforeFunc
	return w
}

func (w *Writer) Do(request any, result any, queryParams map[string]string, method string, path string) error {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	url := w.baseUrl + path
	req.SetRequestURI(url)

	if request != nil {
		jsonBody, err := json.Marshal(request)
		if err != nil {
			return fmt.Errorf("error to marshal request body: %s", err.Error())
		}
		req.SetBody(jsonBody)
	}

	if queryParams != nil {
		for key, value := range queryParams {
			req.URI().QueryArgs().Add(key, value)
		}
	}

	if w.beforeRequestFunc != nil {
		err := w.beforeRequestFunc(req)
		if err != nil {
			return err
		}
	}

	res := fasthttp.AcquireResponse()
	err := fasthttp.Do(req, res)
	if err != nil {
		return fmt.Errorf("error to send request: %s", err.Error())
	}
	w.log(req, res)
	if res.StatusCode() != 200 {
		return fmt.Errorf("error to send request: bad status code")
	}
	fasthttp.ReleaseRequest(req)

	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		return fmt.Errorf("error to unmarshal response: %s", err.Error())
	}
	fasthttp.ReleaseResponse(res)
	return nil
}

func (w *Writer) log(req *fasthttp.Request, resp *fasthttp.Response) {
	fmt.Printf("\n------------------NEW REQ------------------\n%s\n-------------------------------------------\n%s\n", req.String(), resp.String())
}
