package sdk

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Client struct {
	AccessKeyId     string
	AccessKeySecret string
	EndPoint        string
}

type Request interface {
	GetHeaders() map[string]string
	SetHeader(key, value string)
	GetPath() string
	GetScheme() string
	GetMethod() string
	GetDomain() string
	SetDomain(domain string)
	GetQuery() map[string]string
	GetPort() string
	GetBodyReader() io.Reader
	BuildQueries() string
	GetTimeOut() time.Duration
	BuildUrl() string
	GetContent() []byte
}
type Response interface {
	IsSuccess() bool
	HttpStatus() int
	parseFromHttpResponse(httpResponse *http.Response) (err error)
	GetContentString() string
	GetContentBytes() []byte
}

func (c Client) Sign(rest Request) {
	signedStr := ShaHmac1(BuildStringToSign(rest), c.AccessKeySecret)
	rest.SetHeader("Authorization", fmt.Sprintf("OPENSEARCH %s:%s", c.AccessKeyId, signedStr))
}
func initRequestSign(rest Request, domain string) {
	rest.SetDomain(domain)
	rest.SetHeader("Date", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	rest.SetHeader("X-Opensearch-Nonce", fmt.Sprintf("%d", rand.Int()))
	content := rest.GetContent()
	if len(content) > 0 {
		rest.SetHeader("Content-MD5", GetMD5Base64(content))
	}
}
func (c Client) DoAction(rest Request, resp Response) error {
	initRequestSign(rest, c.EndPoint)
	c.Sign(rest)

	httpClient := http.Client{
		Timeout: rest.GetTimeOut(),
	}

	httpRest, err := http.NewRequest(rest.GetMethod(), rest.BuildUrl(), rest.GetBodyReader())
	if err != nil {
		return err
	}
	for key, value := range rest.GetHeaders() {
		httpRest.Header[key] = []string{value}
	}

	if host, containsHost := rest.GetHeaders()["Host"]; containsHost {
		httpRest.Host = host
	}
	httpResp, err := httpClient.Do(httpRest)
	if err != nil {
		return err
	}

	err = resp.parseFromHttpResponse(httpResp)
	if err != nil {
		return err
	}
	return nil
}

type OpenRequest struct {
	Headers    map[string]string
	Path       string
	Scheme     string
	Method     string
	Domain     string
	Port       string
	Query      map[string]string
	FormParams map[string]string
	Content    []byte
	TimeOut    time.Duration
}

func NewOpenRequest() *OpenRequest {
	rest := &OpenRequest{
		Headers: map[string]string{
			"Content-MD5":  "",
			"Content-Type": "application/json",
		},
		Query:      make(map[string]string),
		FormParams: make(map[string]string),
		TimeOut:    30 * time.Second,
		Scheme:     "http",
		Method:     MethodGet,
		Port:       "80",
	}
	return rest
}
func (rest *OpenRequest) GetHeaders() map[string]string {
	return rest.Headers
}
func (rest *OpenRequest) GetPath() string {
	return rest.Path
}
func (rest *OpenRequest) GetScheme() string {
	return rest.Scheme
}
func (rest *OpenRequest) GetMethod() string {
	return rest.Method
}
func (rest *OpenRequest) GetDomain() string {
	return rest.Domain
}
func (rest *OpenRequest) GetPort() string {
	return rest.Port
}
func (rest *OpenRequest) GetBodyReader() io.Reader {
	if rest.FormParams != nil && len(rest.FormParams) > 0 {
		formString := GetUrlFormedMap(rest.FormParams)
		return strings.NewReader(formString)
	} else if len(rest.Content) > 0 {
		return bytes.NewReader(rest.Content)
	} else {
		return nil
	}
}
func (rest *OpenRequest) BuildQueries() string {
	// append urlBuilder
	urlBuilder := bytes.Buffer{}
	urlBuilder.WriteString(rest.Path)
	if len(rest.Query) > 0 {
		urlBuilder.WriteString("?")
		urlBuilder.WriteString(rest.buildQueryString())
	}
	result := urlBuilder.String()
	return result
}
func (rest *OpenRequest) GetQuery() map[string]string {
	return rest.Query
}
func (rest *OpenRequest) GetTimeOut() time.Duration {
	return rest.TimeOut
}
func (rest *OpenRequest) BuildUrl() string {
	// for network trans, need url encoded
	scheme := strings.ToLower(rest.Scheme)
	domain := rest.Domain
	port := rest.Port
	path := rest.Path
	u := fmt.Sprintf("%s://%s:%s%s", scheme, domain, port, path)
	querystring := rest.buildQueryString()
	if len(querystring) > 0 {
		u = fmt.Sprintf("%s?%s", u, querystring)
	}
	return u
}
func (rest *OpenRequest) buildQueryString() string {
	queryParams := rest.Query
	// sort QueryParams by key
	q := url.Values{}
	for key, value := range queryParams {
		q.Add(key, value)
	}
	return strings.ReplaceAll(q.Encode(), "+", "%20")
}
func (rest *OpenRequest) SetDomain(domain string) {
	rest.Domain = domain
}
func (rest *OpenRequest) SetHeader(key, value string) {
	rest.Headers[key] = value
}
func (rest *OpenRequest) GetContent() []byte {
	return rest.Content
}

type OpenResponse struct {
	httpStatus         int
	httpHeaders        http.Header
	httpContentBytes   []byte
	httpContentString  string
	originHttpResponse *http.Response
}

func (resp *OpenResponse) HttpStatus() int {
	return resp.httpStatus
}
func (resp *OpenResponse) IsSuccess() bool {
	return resp.httpStatus == http.StatusOK
}
func (resp *OpenResponse) GetContentString() string {
	return resp.httpContentString
}
func (resp *OpenResponse) GetContentBytes() []byte {
	return resp.httpContentBytes
}
func (resp *OpenResponse) parseFromHttpResponse(httpResponse *http.Response) (err error) {
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	resp.httpStatus = httpResponse.StatusCode
	resp.httpHeaders = httpResponse.Header
	resp.httpContentBytes = body
	resp.httpContentString = string(body)
	resp.originHttpResponse = httpResponse
	return
}

func GetUrlFormedMap(source map[string]string) (urlEncoded string) {
	urlEncoder := url.Values{}
	for key, value := range source {
		urlEncoder.Add(key, value)
	}
	urlEncoded = urlEncoder.Encode()
	return
}
