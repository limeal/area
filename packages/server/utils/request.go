package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// `RequestParams` is a struct that contains a `Method` string, a `Body` string, a `QueryParams` map of
// strings to strings, a `UrlParams` map of strings to strings, and a `Headers` map of strings to
// strings.
// @property {string} Method - The HTTP method to use for the request.
// @property {string} Body - The body of the request.
// @property QueryParams - These are the query parameters that you want to pass in the URL.
// @property UrlParams - These are the parameters that are part of the URL. For example, if you have a
// URL like `/users/{userId}/posts/{postId}`, then `userId` and `postId` are the URL parameters.
// @property Headers - This is a map of key-value pairs that will be added to the request headers.
type RequestParams struct {
	Method      string
	Body        string
	QueryParams map[string]string
	UrlParams   map[string]string
	Headers     map[string]string
}

type RequestDescriptor struct {
	BaseURL           string                                             `json:"url"`    // Base URL of the request
	Params            func(params []interface{}) *RequestParams          `json:"-"`      // Params to send to the request
	ExpectedStatus    []int                                              `json:"status"` // Expected status code of the response
	TransformResponse func(response any) (map[string]interface{}, error) `json:"-"`      // Transform the response
}

// Calling the `DoRequest` function and returning the response.
func (rd *RequestDescriptor) Call(params []interface{}) (map[string]interface{}, *http.Response, error) {
	resp, _, res, err := DoRequest(rd.BaseURL, rd.Params(params), rd.ExpectedStatus, true)
	if err != nil {
		return nil, res, err
	}
	if rd.TransformResponse == nil {
		return resp.(map[string]interface{}), res, nil
	}
	m, err := rd.TransformResponse(resp)
	return m, res, err
}

// Calling the `DoRequest` function and returning the response.
func (rd *RequestDescriptor) CallPure(params []interface{}) (any, *http.Response, error) {
	resp, _, res, err := DoRequest(rd.BaseURL, rd.Params(params), rd.ExpectedStatus, true)
	if err != nil {
		return resp, res, err
	}
	return resp, res, nil
}

// Calling the `DoRequest` function and returning the response.
func (rd *RequestDescriptor) CallEncode(params []interface{}) ([]byte, *http.Response, error) {
	_, enc, res, err := DoRequest(rd.BaseURL, rd.Params(params), rd.ExpectedStatus, false)
	if err != nil {
		return enc, res, err
	}
	return enc, res, nil
}

// -----------------------------------UTILS--------------------------------------------

// It takes a base URL, a method, a body, a map of URL parameters, a map of query parameters, and a map
// of headers, and returns a response
func MakeRequest(baseurl string, p *RequestParams) (*http.Response, error) {
	if p == nil {
		return nil, fmt.Errorf("Request params is nil")
	}
	client := &http.Client{}
	for urlKey, urlParam := range p.UrlParams {
		baseurl = strings.Replace(baseurl, "${"+urlKey+"}", urlParam, -1)
	}
	fmt.Println("URL :> ", baseurl)
	req, err := http.NewRequest(p.Method, baseurl, bytes.NewBufferString(p.Body))
	if err != nil {
		return nil, err
	}
	queryParams := url.Values{}
	for key := range p.QueryParams {
		queryParams.Add(key, p.QueryParams[key])
	}
	req.URL.RawQuery = queryParams.Encode()
	for key := range p.Headers {
		req.Header.Add(key, p.Headers[key])
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Read the body of the response and return it as a byte array.
func RequestReadBody(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// It makes a request, reads the body, and decodes the body into a struct if the status code is in the
// expectedStatus array
func DoRequest(baseurl string, p *RequestParams, expectedStatus []int, decode bool) (any, []byte, *http.Response, error) {
	resp, err := MakeRequest(baseurl, p)
	if err != nil {
		return nil, nil, nil, err
	}
	body, err := RequestReadBody(resp)
	if err != nil {
		return nil, nil, resp, err
	}
	if !decode {
		for _, eS := range expectedStatus {
			if resp.StatusCode == eS {
				return nil, body, resp, nil
			}
		}
		return nil, body, resp, fmt.Errorf("%d - Unexpected status code expected one of %v", resp.StatusCode, expectedStatus)
	}
	var retResponse any
	if err = json.Unmarshal(body, &retResponse); err != nil {
		return nil, nil, resp, err
	}
	for _, eS := range expectedStatus {
		if resp.StatusCode == eS {
			return retResponse, nil, resp, nil
		}
	}
	return nil, nil, resp, fmt.Errorf("%d - Unexpected status code expected one of %v", resp.StatusCode, expectedStatus)
}
