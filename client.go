/*
 * Gate API v4
 *
 * APIv4 provides spot, margin and futures trading operations. There are public APIs to retrieve the real-time market statistics, and private APIs which needs authentication to trade on user's behalf.
 *
 * Contact: support@mail.gate.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gateapi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	jsonCheck = regexp.MustCompile("(?i:[application|text]/json)")
	xmlCheck  = regexp.MustCompile("(?i:[application|text]/xml)")
)

// APIClient manages communication with the Gate API v4 API v4.8.0
// In most cases there should be only one, shared, APIClient.
type APIClient struct {
	cfg    *Configuration
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// API Services

	FuturesApi *FuturesApiService

	MarginApi *MarginApiService

	SpotApi *SpotApiService

	WalletApi *WalletApiService
}

type service struct {
	client *APIClient
}

// NewAPIClient creates a new API client. Requires a userAgent string describing your application.
// optionally a custom http.Client to allow for advanced features such as caching.
func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.FuturesApi = (*FuturesApiService)(&c.common)
	c.MarginApi = (*MarginApiService)(&c.common)
	c.SpotApi = (*SpotApiService)(&c.common)
	c.WalletApi = (*WalletApiService)(&c.common)

	return c
}

func atoi(in string) (int, error) {
	return strconv.Atoi(in)
}

// selectHeaderContentType select a content type from the available list.
func selectHeaderContentType(contentTypes []string) string {
	if len(contentTypes) == 0 {
		return ""
	}
	if contains(contentTypes, "application/json") {
		return "application/json"
	}
	return contentTypes[0] // use the first content type specified in 'consumes'
}

// selectHeaderAccept join all accept types and return
func selectHeaderAccept(accepts []string) string {
	if len(accepts) == 0 {
		return ""
	}

	if contains(accepts, "application/json") {
		return "application/json"
	}

	return strings.Join(accepts, ",")
}

// contains is a case insenstive match, finding needle in a haystack
func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.ToLower(a) == strings.ToLower(needle) {
			return true
		}
	}
	return false
}

// Verify optional parameters are of the correct type.
func typeCheckParameter(obj interface{}, expected string, name string) error {
	// Make sure there is an object.
	if obj == nil {
		return nil
	}

	// Check the type is as expected.
	if reflect.TypeOf(obj).String() != expected {
		return fmt.Errorf("Expected %s to be of type %s but received %s.", name, expected, reflect.TypeOf(obj).String())
	}
	return nil
}

// parameterToString convert interface{} parameters to string, using a delimiter if format is provided.
func parameterToString(obj interface{}, collectionFormat string) string {
	var delimiter string

	switch collectionFormat {
	case "pipes":
		delimiter = "|"
	case "ssv":
		delimiter = " "
	case "tsv":
		delimiter = "\t"
	case "csv":
		delimiter = ","
	}

	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		return strings.Trim(strings.Replace(fmt.Sprint(obj), " ", delimiter, -1), "[]")
	} else if t, ok := obj.(time.Time); ok {
		return t.Format(time.RFC3339)
	}

	return fmt.Sprintf("%v", obj)
}

// callAPI do the request.
func (c *APIClient) callAPI(request *http.Request) (*http.Response, error) {
	return c.cfg.HTTPClient.Do(request)
}

// Change base path to allow switching to mocks
func (c *APIClient) ChangeBasePath(path string) {
	c.cfg.BasePath = path
}

func (c *APIClient) SetKeySecret(key string, secret string) {
	c.cfg.Key = key
	c.cfg.Secret = secret
}

// prepareRequest build the request
func (c *APIClient) prepareRequest(
	ctx context.Context,
	path string, method string,
	postBody interface{},
	headerParams map[string]string,
	queryParams url.Values,
	formParams url.Values,
	formFileName string,
	fileName string,
	fileBytes []byte,
	authRequired bool) (localVarRequest *http.Request, err error) {

	var body *bytes.Buffer

	// Detect postBody type and post.
	if postBody != nil {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}

		body, err = setBody(postBody, contentType)
		if err != nil {
			return nil, err
		}
	}

	// add form parameters and file if available.
	if strings.HasPrefix(headerParams["Content-Type"], "multipart/form-data") && len(formParams) > 0 || (len(fileBytes) > 0 && fileName != "") {
		if body != nil {
			return nil, errors.New("Cannot specify postBody and multipart form at the same time.")
		}
		body = &bytes.Buffer{}
		w := multipart.NewWriter(body)

		for k, v := range formParams {
			for _, iv := range v {
				if strings.HasPrefix(k, "@") { // file
					err = addFile(w, k[1:], iv)
					if err != nil {
						return nil, err
					}
				} else { // form value
					w.WriteField(k, iv)
				}
			}
		}
		if len(fileBytes) > 0 && fileName != "" {
			w.Boundary()
			//_, fileNm := filepath.Split(fileName)
			part, err := w.CreateFormFile(formFileName, filepath.Base(fileName))
			if err != nil {
				return nil, err
			}
			_, err = part.Write(fileBytes)
			if err != nil {
				return nil, err
			}
			// Set the Boundary in the Content-Type
			headerParams["Content-Type"] = w.FormDataContentType()
		}

		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
		w.Close()
	}

	if strings.HasPrefix(headerParams["Content-Type"], "application/x-www-form-urlencoded") && len(formParams) > 0 {
		if body != nil {
			return nil, errors.New("Cannot specify postBody and x-www-form-urlencoded form at the same time.")
		}
		body = &bytes.Buffer{}
		body.WriteString(formParams.Encode())
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
	}

	// Setup path and query parameters
	requestUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Adding Query Param
	query := requestUrl.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	// Encode the parameters.
	requestUrl.RawQuery = query.Encode()

	// generate signature if needed
	if authRequired {
		h := sha512.New()
		if body != nil {
			h.Write(body.Bytes())
		}
		hashedPayload := hex.EncodeToString(h.Sum(nil))

		t := strconv.FormatInt(time.Now().Unix(), 10)
		rawQuery, err := url.QueryUnescape(requestUrl.RawQuery)
		if err != nil {
			return nil, err
		}
		msg := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, requestUrl.Path, rawQuery, hashedPayload, t)
		mac := hmac.New(sha512.New, []byte(c.cfg.Secret))
		mac.Write([]byte(msg))
		sign := hex.EncodeToString(mac.Sum(nil))
		headerParams["KEY"] = c.cfg.Key
		headerParams["SIGN"] = sign
		headerParams["Timestamp"] = t
	}

	// Generate a new request
	if body != nil {
		localVarRequest, err = http.NewRequest(method, requestUrl.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, requestUrl.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers.Set(h, v)
		}
		localVarRequest.Header = headers
	}

	// Override request host, if applicable
	if c.cfg.Host != "" {
		localVarRequest.Host = c.cfg.Host
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.cfg.UserAgent)

	if ctx != nil {
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)
	}

	for header, value := range c.cfg.DefaultHeader {
		localVarRequest.Header.Add(header, value)
	}

	return localVarRequest, nil
}

func (c *APIClient) decode(v interface{}, b []byte, contentType string) (err error) {
		if strings.Contains(contentType, "application/xml") {
			if err = xml.Unmarshal(b, v); err != nil {
				return err
			}
			return nil
		} else if strings.Contains(contentType, "application/json") {
			if err = json.Unmarshal(b, v); err != nil {
				return err
			}
			return nil
		}
	return errors.New("undefined response type")
}

// Add a file to the multipart request
func addFile(w *multipart.Writer, fieldName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := w.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	return err
}

// Prevent trying to import "fmt"
func reportError(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// Set request body from an interface{}
func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer, err error) {
	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if jsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	} else if xmlCheck.MatchString(contentType) {
		xml.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("Invalid body type %s\n", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

// detectContentType method is used to figure out `Request.Body` content type for request header
func detectContentType(body interface{}) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}

// Ripped from https://github.com/gregjones/httpcache/blob/master/httpcache.go
type cacheControl map[string]string

func parseCacheControl(headers http.Header) cacheControl {
	cc := cacheControl{}
	ccHeader := headers.Get("Cache-Control")
	for _, part := range strings.Split(ccHeader, ",") {
		part = strings.Trim(part, " ")
		if part == "" {
			continue
		}
		if strings.ContainsRune(part, '=') {
			keyval := strings.Split(part, "=")
			cc[strings.Trim(keyval[0], " ")] = strings.Trim(keyval[1], ",")
		} else {
			cc[part] = ""
		}
	}
	return cc
}

// CacheExpires helper function to determine remaining time before repeating a request.
func CacheExpires(r *http.Response) time.Time {
	// Figure out when the cache expires.
	var expires time.Time
	now, err := time.Parse(time.RFC1123, r.Header.Get("date"))
	if err != nil {
		return time.Now()
	}
	respCacheControl := parseCacheControl(r.Header)

	if maxAge, ok := respCacheControl["max-age"]; ok {
		lifetime, err := time.ParseDuration(maxAge + "s")
		if err != nil {
			expires = now
		} else {
			expires = now.Add(lifetime)
		}
	} else {
		expiresHeader := r.Header.Get("Expires")
		if expiresHeader != "" {
			expires, err = time.Parse(time.RFC1123, expiresHeader)
			if err != nil {
				expires = now
			}
		}
	}
	return expires
}

func strlen(s string) int {
	return utf8.RuneCountInString(s)
}

// GenericOpenAPIError Provides access to the body, error and model on returned errors.
type GenericOpenAPIError struct {
	body  []byte
	error string
	model interface{}
}

// Error returns non-empty string if there was an error.
func (e GenericOpenAPIError) Error() string {
	return e.error
}

// Body returns the raw bytes of the response
func (e GenericOpenAPIError) Body() []byte {
	return e.body
}

// Model returns the unpacked model of the error
func (e GenericOpenAPIError) Model() interface{} {
	return e.model
}
