package action

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"reflect"
	"strings"
)

type http struct {
	Uri       string            `json:"uri" bson:"uri"`
	Header    map[string]string `json:"header" bson:"header"`
	Method    string            `json:"method" bson:"method"`
	Arguments map[string]string `json:"arguments" bson:"arguments"`
	AuthToken string            `json:"auth_token"`
	Response  struct {
		Status int         `json:"status" bson:"status"`
		Error  interface{} `json:"error" bson:"error"`
	} `json:"response" bson:"response"`

	_call []func() error
}

func newHttp() *http {
	return &http{
		_call:  make([]func() error, 0),
		Method: "post", //current only support post
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
	for _, url := range urls {
		http._call = append(http._call,
			func() error {
				req := http.request()
				body, err := json.Marshal(http.Arguments)
				if err != nil {
					return err
				}
				req.SetBody(body)
				response, err := req.Post(url)
				if response != nil {
					http.Response.Status = response.StatusCode()
					http.Response.Error = response.Error()
				}
				if err != nil {
					return err
				}
				return nil
			})
	}
	return http
}

func (http *http) Params(p map[string]interface{}) HttpInterface {
	http.Arguments = make(map[string]string)
	for k, v := range p {
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.String:
			http.Arguments[k] = fmt.Sprintf("%s", v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			http.Arguments[k] = fmt.Sprintf("%d", v)
		default:
			http.Arguments[k] = fmt.Sprintf("%v", v)
		}
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
		if err := f(); err != nil {
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
