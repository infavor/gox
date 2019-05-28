package httpx

import (
	"bytes"
	"compress/gzip"
	"errors"
	json "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	METHOD_GET                         = "GET"
	METHOD_POST                        = "POST "
	METHOD_TRACE                       = "TRACE"
	METHOD_DELETE                      = "DELETE"
	METHOD_PUT                         = "PUT"
	METHOD_OPTIONS                     = "OPTIONS"
	METHOD_HEAD                        = "HEAD"
	METHOD_CONNECT                     = "CONNECT"
	CONTENT_TYPE_X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded"
	CONTENT_TYPE_JSON                  = "application/json"
)

var (
	httpClient     *http.Client
	defaultHeaders = make(map[string]string, 10)
)

func init() {
	httpClient = &http.Client{
		Timeout: time.Second * 20,
	}
	defaultHeaders["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
}

type mock struct {
	url               string
	method            string
	headers           map[string]string
	contentType       string
	parameterMap      map[string][]string
	body              []byte
	request           http.Request
	response          http.Response
	responseContainer interface{}
	callback          func(status int, response []byte)
}

func Mock() *mock {
	return &mock{
		method:      METHOD_GET,
		headers:     map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"},
		contentType: CONTENT_TYPE_X_WWW_FORM_URLENCODED,
	}
}

func (m *mock) URL(url string) *mock {
	m.url = url
	return m
}

func (m *mock) Header(name, value string) *mock {
	m.headers[name] = value
	return m
}

func (m *mock) Headers(headers map[string]string) *mock {
	for k, v := range m.headers {
		m.headers[k] = v
	}
	return m
}

func (m *mock) ContentType(contentType string) *mock {
	if strings.HasPrefix(contentType, CONTENT_TYPE_X_WWW_FORM_URLENCODED) ||
		strings.HasPrefix(contentType, CONTENT_TYPE_JSON) {
		m.contentType = contentType
	} else {
		panic(errors.New("not supported contentType: '" + contentType +
			"', contentType is currently only support " + "'" + CONTENT_TYPE_X_WWW_FORM_URLENCODED +
			" and '" + CONTENT_TYPE_JSON + "'"))
	}
	return m
}

// parameter will be generated on url
func (m *mock) Parameters(params map[string]string) *mock {
	if params != nil && len(params) > 0 {
		for k, v := range params {
			oldV := m.parameterMap[k]
			if oldV == nil {
				oldV = make([]string, len(oldV)+1)
			}
			oldV[len(oldV)] = v
			m.parameterMap[k] = oldV
		}
	}
	return m
}

// parameter will be generated on url
func (m *mock) Parameter(name, value string) *mock {
	oldV := m.parameterMap[name]
	if oldV == nil {
		oldV = make([]string, len(oldV)+1)
	}
	oldV[len(oldV)] = value
	m.parameterMap[name] = oldV
	return m
}

// Body must be a type like map[string][]string or custom struct
// if contentType is 'application/x-www-form-urlencoded' then type of body must be type map[string][]string,
// if contentType is 'application/json' then type of body could be any.
func (m *mock) Body(body interface{}) *mock {
	if body == nil {
		m.body = nil
		return m
	}
	bodyType := reflect.TypeOf(body).String()
	if m.contentType == CONTENT_TYPE_JSON {
		jv, _ := json.Marshal(body)
		m.body = jv
	} else {
		if bodyType != "map[string][]string" {
			panic(errors.New("body type must be 'map[string][]string' if contentType is 'application/x-www-form-urlencoded'"))
		}
		m.body = encodeParameters(body.(map[string][]string))
	}
	return m
}

func (m *mock) Success(response interface{}) *mock {
	m.responseContainer = response
	return m
}
func (m *mock) Error(callback func(status int, response []byte)) *mock {
	m.callback = callback
	return m
}

func (m *mock) Do() error {
	paramsBytes := encodeParameters(m.parameterMap)
	req, err := http.NewRequest(m.method, strings.Join([]string{m.url + "?" + string(paramsBytes)}, ""), bytes.NewReader(m.body))
	if err != nil {
		return err
	}
	for k, v := range m.headers {
		req.Header.Add(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		m.callback(resp.StatusCode, bs)
		return nil
	}
	decodeUseGzip := false
	if resp.Header != nil {
		for k, v := range resp.Header {
			// fmt.Println(k, v[0])
			if strings.ToLower(k) == "content-encoding" {
				if v != nil && len(v) > 0 {
					if v[0] == "gzip" {
						decodeUseGzip = true
						break
					}
				}
			}
		}
	}
	var body []byte
	if decodeUseGzip {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		bs, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		body = bs
	} else {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		body = bs
	}
	return json.Unmarshal(body, m.responseContainer)
}

func (m *mock) Get() *mock {
	m.method = METHOD_GET
	return m
}
func (m *mock) Post() *mock {
	m.method = METHOD_POST
	return m
}
func (m *mock) Options() *mock {
	m.method = METHOD_OPTIONS
	return m
}
func (m *mock) Head() *mock {
	m.method = METHOD_HEAD
	return m
}
func (m *mock) Put() *mock {
	m.method = METHOD_PUT
	return m
}
func (m *mock) Delete() *mock {
	m.method = METHOD_DELETE
	return m
}
func (m *mock) Connect() *mock {
	m.method = METHOD_CONNECT
	return m
}
func (m *mock) Trace() *mock {
	m.method = METHOD_TRACE
	return m
}

func encodeParameters(params map[string][]string) []byte {
	if params == nil || len(params) == 0 {
		return []byte{}
	}
	var buffer bytes.Buffer
	for k, vl := range params {
		if vl == nil || len(vl) == 0 {
			buffer.WriteString(k)
			buffer.WriteString("=")
			continue
		}
		for i, v := range vl {
			buffer.WriteString(k)
			buffer.WriteString("=")
			buffer.WriteString(v)
			if i != len(vl)-1 {
				buffer.WriteString("&")
			}
		}
	}
	return buffer.Bytes()
}
