package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zsulocal/upbit-go/types"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type RequestOptions struct {
	Url     string
	Method  string
	Body    map[string]string
	Json    map[string]string
	Query   map[string]string
	Headers map[string]string
}

func Request(options *RequestOptions, result interface{}) (
	err error,
) {
	client := &http.Client{}
	var rawbody io.Reader
	if options.Body != nil {
		body := url.Values{}
		for k, v := range options.Body {
			body.Set(k, v)
		}
		rawbody = strings.NewReader(body.Encode())
	}

	if options.Json != nil {
		body, _ := json.Marshal(options.Json)
		rawbody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(options.Method, options.Url, rawbody)
	if err != nil {
		return
	}

	if options.Query != nil {
		q := req.URL.Query()
		for index, value := range options.Query {
			q.Add(index, value)
		}

		req.URL.RawQuery = q.Encode()
	}

	if options.Headers != nil {
		for prop, value := range options.Headers {
			req.Header.Add(prop, value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	Body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(Body, result)
	var upbitErr types.ResponseError
	if err != nil {
		err = json.Unmarshal(Body, &upbitErr)
		if err == nil {
			err = errors.New(upbitErr.Err.Message)
			return
		}
	}
	if err != nil {
		return
	}
	return
}
