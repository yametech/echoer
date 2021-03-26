package action

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	http2 "net/http"
	"strings"
)

type http struct {
	Uri       string                 `json:"uri" bson:"uri"`
	Header    map[string]string      `json:"header" bson:"header"`
	Method    string                 `json:"method" bson:"method"`
	Arguments map[string]interface{} `json:"arguments" bson:"arguments"`
	AuthToken string                 `json:"auth_token"`
	Response  struct {
		Status int         `json:"status" bson:"status"`
		Error  interface{} `json:"error" bson:"error"`
	} `json:"response" bson:"response"`

	_call []func() error
}

type https struct {
	Uri       string                 `json:"uri" bson:"uri"`
	Header    map[string]string      `json:"header" bson:"header"`
	Method    string                 `json:"method" bson:"method"`
	Arguments map[string]interface{} `json:"arguments" bson:"arguments"`
	AuthToken string                 `json:"auth_token"`
	CaPEM     string                 `json:"capem" bson:"capem"`
	Response  struct {
		Status int         `json:"status" bson:"status"`
		Error  interface{} `json:"error" bson:"error"`
	} `json:"response" bson:"response"`

	_call []func() error
}

func (h https) Post(urls []string) HttpsInterface {
	if len(h._call) > 0 {
		h._call = h._call[:0]
	}
	for index, url := range urls {
		_ = url
		_url := urls[index]
		_func := func() error {
			client, err := h.client()
			if err != nil {
				return err
			}
			body, err := json.Marshal(h.Arguments)
			if err != nil {
				return err
			}

			req, err := http2.NewRequest(http2.MethodPost, _url, bytes.NewReader(body))
			if err != nil{
				return err
			}
			response, err := client.Do(req)
			if response != nil {
				h.Response.Status = response.StatusCode
				h.Response.Error = err.Error()
				if h.Response.Status != 200 {
					return fmt.Errorf(
						"webhook post to (%s) response code (%d) data (%v) error (%s)",
						url,
						response.StatusCode,
						h.Arguments,
						err.Error(),
					)
				}
			}
			if err != nil {
				return err
			}
			return nil
		}
		h._call = append(h._call, _func)
	}
	return h
}

func (h https) Params(p map[string]interface{}) HttpsInterface {
	h.Arguments = make(map[string]interface{})
	for k, v := range p {
		h.Arguments[k] = v
	}
	h.CaPEM = p["capem"].(string)
	return h
}

func (h https) Do() error {
	var err error
	switch strings.ToLower(h.Method) {
	case "post":
	case "get":
	default:
		err = fmt.Errorf("not found method (%s)", h.Method)
		return err
	}
	for i, f := range h._call {
		_func := f
		if err := _func(); err != nil {
			if (i + 1) == len(h._call) {
				return err
			}
			continue
		}
		break
	}
	return nil
}

func (h https) client() (*http2.Client, error) {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(h.CaPEM))
	if !ok {
		return nil, fmt.Errorf("failed to parse root certificate")
	}

	tr := &http2.Transport{
		TLSClientConfig: &tls.Config{RootCAs: roots},
	}

	client := &http2.Client{Transport: tr}
	return client, nil
}

func newHttp() *http {
	return &http{
		_call:  make([]func() error, 0),
		Method: "post", //current only support post
	}
}

func newHttps() *https {
	return &https{
		_call:  make([]func() error, 0),
		Method: "post",
	}
}

func (http *http) request() *resty.Request {
	client := resty.New()
	req := client.R().SetHeader("Accept", "application/json") //default json
	if http.Header != nil {
		for k, v := range http.Header {
			req.SetHeader(k, v)
		}
	}
	if http.AuthToken != "" {
		req.SetAuthToken(http.AuthToken)
	}
	return req
}

func (http *http) HttpInterface() HttpInterface {
	return http
}

func (http *http) Post(urls []string) HttpInterface {
	if len(http._call) > 0 {
		http._call = http._call[:0]
	}
	for index, url := range urls {
		_ = url
		_url := urls[index]
		_func := func() error {
			req := http.request()
			body, err := json.Marshal(http.Arguments)
			if err != nil {
				return err
			}
			req.SetBody(body)
			response, err := req.Post(_url)
			if response != nil {
				http.Response.Status = response.StatusCode()
				http.Response.Error = response.Error()
				if http.Response.Status != 200 {
					return fmt.Errorf(
						"webhook post to (%s) response code (%d) data (%v) error (%s)",
						url,
						response.StatusCode(),
						http.Arguments,
						response.Error(),
					)
				}
			}
			if err != nil {
				return err
			}
			return nil
		}
		http._call = append(http._call, _func)
	}
	return http
}

func (http *http) Params(p map[string]interface{}) HttpInterface {
	http.Arguments = make(map[string]interface{})
	for k, v := range p {
		http.Arguments[k] = v
	}
	return http
}

func (http *http) Do() error {
	var err error
	switch strings.ToLower(http.Method) {
	case "post":
	case "get":
	default:
		err = fmt.Errorf("not found method (%s)", http.Method)
		return err
	}
	for i, f := range http._call {
		_func := f
		if err := _func(); err != nil {
			if (i + 1) == len(http._call) {
				return err
			}
			continue
		}
		break
	}
	return nil
}

type FakeHttp struct {
	http
}

func (http *FakeHttp) Post(urls []string) error {
	return nil
}

func (http *FakeHttp) Params(urls []string) error {
	return nil
}

func (http *FakeHttp) Do() error {
	return nil
}
