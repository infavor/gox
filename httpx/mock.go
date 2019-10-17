package httpx

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	json "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	METHOD_GET                         = "GET"
	METHOD_POST                        = "POST"
	METHOD_TRACE                       = "TRACE"
	METHOD_DELETE                      = "DELETE"
	METHOD_PUT                         = "PUT"
	METHOD_OPTIONS                     = "OPTIONS"
	METHOD_HEAD                        = "HEAD"
	METHOD_CONNECT                     = "CONNECT"
	CONTENT_TYPE_X_WWW_FORM_URLENCODED = "application/x-www-form-urlencoded"
	CONTENT_TYPE_JSON                  = "application/json"
	CONTENT_TYPE_MULTIPART             = "multipart/form-data"
)

var (
	httpClient           *http.Client
	defaultHeaders       = make(map[string]string, 10)
	allowedResponseTypes = make(map[string]bool)
)

func init() {
	allowedResponseTypes["int"] = true
	allowedResponseTypes["int64"] = true
	allowedResponseTypes["float32"] = true
	allowedResponseTypes["float64"] = true
	allowedResponseTypes["bool"] = true
	allowedResponseTypes["string"] = true
	allowedResponseTypes["struct"] = true
	allowedResponseTypes["map"] = true
	allowedResponseTypes["nil"] = true
	allowedResponseTypes["io.Writer"] = true

	httpClient = &http.Client{
		Timeout: time.Second * 20,
	}
}

// SetTTL sets http client request timeout value.
// Default timeout value is 20s.
// A Timeout of zero means no timeout.
func SetTTL(timeout time.Duration) {
	httpClient.Timeout = timeout
}

// mock is a fake http request instance.
type mock struct {
	url               string
	method            string
	headers           map[string]string
	contentType       string
	parameterMap      map[string][]string
	body              []byte
	request           http.Request
	response          http.Response
	multipartFiller   func(writer *multipart.Writer)
	responseContainer interface{}
	callback          func(status int, response []byte)
	successCodes      []int
}

// Mock returns an initialized mock.
func Mock() *mock {
	return &mock{
		method:            METHOD_GET,
		headers:           map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"},
		contentType:       CONTENT_TYPE_X_WWW_FORM_URLENCODED,
		successCodes:      []int{http.StatusOK},
		responseContainer: "",
		parameterMap:      make(map[string][]string),
	}
}

// URL sets the mock url.
func (m *mock) URL(url string) *mock {
	m.url = url
	return m
}

// Header adds an http header to the mock.
func (m *mock) Header(name, value string) *mock {
	m.headers[name] = value
	return m
}

// Headers adds many http headers to the mock.
func (m *mock) Headers(headers map[string]string) *mock {
	for k, v := range m.headers {
		m.headers[k] = v
	}
	return m
}

// ContentType sets ContentType of the mock.
func (m *mock) ContentType(contentType string) *mock {
	if strings.HasPrefix(contentType, CONTENT_TYPE_X_WWW_FORM_URLENCODED) ||
		strings.HasPrefix(contentType, CONTENT_TYPE_JSON) ||
		strings.HasPrefix(contentType, CONTENT_TYPE_MULTIPART) {
		m.contentType = contentType
	} else {
		panic(errors.New("not supported contentType: '" + contentType +
			"', contentType is currently only support " + "'" + CONTENT_TYPE_X_WWW_FORM_URLENCODED +
			"', '" + CONTENT_TYPE_MULTIPART + "' and '" + CONTENT_TYPE_JSON + "'"))
	}
	return m
}

// Parameters add parameters on the request url.
func (m *mock) Parameters(params map[string]string) *mock {
	if params != nil && len(params) > 0 {
		for k, v := range params {
			oldV := m.parameterMap[k]
			oldLen := len(oldV)
			if oldV == nil {
				oldV = make([]string, oldLen+1)
			}
			oldV[oldLen] = v
			m.parameterMap[k] = oldV
		}
	}
	return m
}

// Parameter add a parameter on the request url.
func (m *mock) Parameter(name, value string) *mock {
	oldV := m.parameterMap[name]
	oldLen := len(oldV)
	if oldV == nil {
		oldV = make([]string, oldLen+1)
	}
	oldV[oldLen] = value
	m.parameterMap[name] = oldV
	return m
}

// Body must be a type like map[string][]string or custom struct.
//
// if contentType is 'application/x-www-form-urlencoded' then type of body must be type map[string][]string,
//
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
	} else if m.contentType == CONTENT_TYPE_MULTIPART {
		if bodyType != "map[string][]string" {
			panic(errors.New("body type must be 'map[string][]string' if contentType is 'application/x-www-form-urlencoded'"))
		}
		m.body = encodeParameters(body.(map[string][]string))
	} else {
		if bodyType != "map[string][]string" {
			panic(errors.New("body type must be 'map[string][]string' if contentType is 'application/x-www-form-urlencoded'"))
		}
		m.body = encodeParameters(body.(map[string][]string))
	}
	return m
}

// Multipart add multipart content to request body.
//
// if Multipart() is called, the request content will set to 'multipart/form-data' automatically.
func (m *mock) Multipart(multipartFiller func(writer *multipart.Writer)) *mock {
	m.multipartFiller = multipartFiller
	m.method = "POST"
	return m
}

// Success defines the response type and tells what status codes should be recognized as success request,
//
// response type must be one of:
//
// int int64 float32 float64 bool string or pointer of a struct.
func (m *mock) Success(response interface{}, successCodes ...int) *mock {
	m.responseContainer = response
	m.successCodes = successCodes
	if !allowedResponseTypes[checkResponseType(response)] {
		panic("response type not allowed")
	}
	return m
}

// Error handles error during the mock request.
func (m *mock) Error(callback func(status int, response []byte)) *mock {
	m.callback = callback
	return m
}

// Do is the end of the mock chain,
// which will send the request and return the result.
func (m *mock) Do() (interface{}, int, error) {
	paramsStr := string(encodeParameters(m.parameterMap))

	isMultipart := false
	var mw *multipart.Writer
	var pipeReader *io.PipeReader
	var pipeWriter *io.PipeWriter
	if m.multipartFiller != nil {
		isMultipart = true
		pipeReader, pipeWriter = io.Pipe()
		mw = multipart.NewWriter(pipeWriter)
		go func() {
			defer pipeWriter.Close()
			defer mw.Close()
			m.multipartFiller(mw)
		}()
	}

	req, err := http.NewRequest(m.method, gox.TValue(paramsStr == "", m.url, m.url+"?"+paramsStr).(string),
		gox.TValue(isMultipart, pipeReader, bytes.NewReader(m.body)).(io.Reader))
	if err != nil {
		return m.responseContainer, 0, err
	}
	for k, v := range m.headers {
		req.Header.Add(k, v)
	}
	if isMultipart {
		req.Header.Set("Content-Type", mw.FormDataContentType())
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return m.responseContainer, 0, err
	}

	if !m.isSuccess(resp.StatusCode) {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return m.responseContainer, resp.StatusCode, err
		}
		if m.callback != nil {
			m.callback(resp.StatusCode, bs)
		}
		return m.responseContainer, resp.StatusCode, nil
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
			return m.responseContainer, resp.StatusCode, err
		}
		bs, err := ioutil.ReadAll(reader)
		if err != nil {
			return m.responseContainer, resp.StatusCode, err
		}
		body = bs
	} else {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return m.responseContainer, resp.StatusCode, err
		}
		body = bs
	}
	ret, err := convertResponse(checkResponseType(m.responseContainer), string(body), m.responseContainer)
	return ret, resp.StatusCode, err
}

// Get sets http method to GET.
func (m *mock) Get() *mock {
	m.method = METHOD_GET
	return m
}

// Post sets http method to Post.
func (m *mock) Post() *mock {
	m.method = METHOD_POST
	return m
}

// Options sets http method to Options.
func (m *mock) Options() *mock {
	m.method = METHOD_OPTIONS
	return m
}

// Head sets http method to Head.
func (m *mock) Head() *mock {
	m.method = METHOD_HEAD
	return m
}

// Put sets http method to Put.
func (m *mock) Put() *mock {
	m.method = METHOD_PUT
	return m
}

// Delete sets http method to Delete.
func (m *mock) Delete() *mock {
	m.method = METHOD_DELETE
	return m
}

// Connect sets http method to Connect.
func (m *mock) Connect() *mock {
	m.method = METHOD_CONNECT
	return m
}

// Trace sets http method to Trace.
func (m *mock) Trace() *mock {
	m.method = METHOD_TRACE
	return m
}

// encodeParameters encodes parameters to the pattern of 'a=xx&b=xx'.
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
		buffer.WriteString("&")
	}
	return buffer.Bytes()
}

// isSuccess determines whether the request is success.
func (m *mock) isSuccess(code int) bool {
	if m.successCodes != nil && len(m.successCodes) > 0 {
		for _, v := range m.successCodes {
			if v == code {
				return true
			}
		}
	}
	return code == http.StatusOK
}

// checkResponseType returns the type of response data container.
func checkResponseType(resp interface{}) string {
	if resp == nil {
		return "nil"
	}
	if _, c := resp.(io.Writer); c {
		return "io.Writer"
	}
	typ := reflect.TypeOf(resp)
	for {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			continue
		}
		break
	}
	return typ.Kind().String()
}

// convertResponse converts response to the type of response.
func convertResponse(typeName string, response string, responseContainer interface{}) (interface{}, error) {
	switch typeName {
	case "nil":
		return response, nil
	case "io.Writer":
		(responseContainer.(io.Writer)).Write([]byte(response))
		return response, nil
	case "int":
		return convert.StrToInt(response)
	case "int64":
		return convert.StrToInt64(response)
	case "float32":
		return convert.StrToFloat32(response)
	case "float64":
		return convert.StrToFloat64(response)
	case "bool":
		return convert.StrToBool(response)
	case "string":
		return response, nil
	case "map":
		err := json.UnmarshalFromString(response, responseContainer)
		return responseContainer, err
	case "struct":
		err := json.UnmarshalFromString(response, responseContainer)
		return responseContainer, err
	}
	return nil, errors.New("cannot convert response")
}
