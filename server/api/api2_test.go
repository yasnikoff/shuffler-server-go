package api

import (
	"reflect"
	"fmt"
	"strings"
	"github.com/pkg/errors"
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"bitbucket.org/yasnikoff/shuffler/app/app_fixtures"
	"io/ioutil"
	"bytes"
	"log"
)

const ApiPrefix string = "api"

// base interface for all test cases
type TestCase interface {
	fmt.Stringer
	Path() []string
	Method() string
	CaseName() string
	SetCaseName(string)
	//T() *testing.T
	AppFx() *app_fixtures.AppFixture
}

type baseTestCase struct {
	//t        *testing.T
	testCase
	appFx    *app_fixtures.AppFixture
	path     []string
	method   string
	postData interface{}
	statusChecker
}

const pathInitCapacity = 8

func newBaseTestCase(t *testing.T) *baseTestCase {
	result := &baseTestCase{
		appFx:app_fixtures.New(t),
		path:make([]string, 1, pathInitCapacity),
		method:http.MethodGet,
		//t:t,
		statusChecker:statusChecker{expectedCode:http.StatusOK},
	}
	result.SetAppName("test")
	return result
}
/*
func (t *baseTestCase)T() *testing.T {
	return t.t
}*/
func (t *baseTestCase)AppFx() *app_fixtures.AppFixture {
	return t.appFx
}
func (t *baseTestCase)AppName() string {
	return t.path[0]
}

func (t *baseTestCase)SetAppName(name string) {
	t.path[0] = name
}
func (t *baseTestCase)Path() []string {
	return t.path
}
func (t *baseTestCase)Method() string {
	return t.method
}
func (t *baseTestCase)PostData() interface{} {
	return t.postData
}

func (t *baseTestCase)Setup(r *testRunner) {
	t.appFx.Setup()
}

func (t *baseTestCase)Cleanup() {
	t.appFx.Cleanup()
}

type testCase struct {
	Name string
}

func (tc *testCase)String() string {
	return tc.Name
}
func (tc *testCase)CaseName() string {
	return tc.String()
}
func (tc *testCase)SetCaseName(name string) {
	tc.Name = name
}

//func (tc *testCase)wrapFormat(format string)string{
//	return fmt.Sprintf("case %s:\n\t", strings.ToUpper(tc.CaseName())) + format
//}
func (tc *testCase)WrapError(err error) error {
	return errors.Wrapf(err, "case %s", strings.ToUpper(tc.CaseName()))
}
func (tc *testCase)WrapErrorf(err error, message string, args ...interface{}) error {
	return tc.WrapError(errors.Wrapf(err, message, args...))
}
//func (tc *testCase)WrapErrorf(err error, format string, args ...interface{})error{
//	return errors.Wrapf(err,  tc.wrapFormat(format), args...)
//}
// test case implement this interface if it needs some setup
type Setupper interface {
	Setup(*testRunner)
	Cleanup()
}

type ErrorWrapper interface {
	WrapError(error) error
	WrapErrorf(error, string, ...interface{})
}

// test case for GET method
type GetTestCase interface {
	TestCase
}

// test case for POST method
type PostTestCase interface {
	TestCase
	PostData() interface{}
}

// if the test case implements this method, it  response status code is checked to be equal to ExpectedStatusCode()
type StatusCodeChecker interface {
	ExpectedStatusCode() int
}

type statusChecker struct {
	expectedCode int
}

func (sch *statusChecker)ExpectedStatusCode() int {
	return sch.expectedCode
}

// if the test case implements this interface, response body is checked to be equal to ExpectedBody()
type BodyChecker interface {
	ExpectedBody() []byte
}

type JsonBodyChecker interface {
	ExpectedJsonSource() interface{}
}

type AppNamer interface {
	AppName() string
	SetAppName(string)
}

type testRunner struct {
	testCase  TestCase

	handler   *handler
	server    *httptest.Server

	Resp      *http.Response
	Body      []byte

	WriteLogs bool
}

func (r *testRunner)Setup(tc TestCase) error {

	r.testCase = tc

	var appName string
	if tc, ok := tc.(AppNamer); ok {
		appName = tc.AppName()
	} else {
		appName = "test"
	}

	r.handler = NewHandler()
	if err := r.handler.registerApp(tc.AppFx().App, appName); err != nil {
		return r.wrapErrorf(err, "can't register app in handler")
	}
	r.server = httptest.NewServer(r.handler)
	return nil
}

func (r *testRunner)Cleanup() {
	r.testCase = nil
	r.server.CloseClientConnections()
	r.server.Close()
}

func (r *testRunner)Run(tc TestCase, t *testing.T, f func(t *testing.T)) {

	if err := r.Setup(tc); err != nil {
		t.Fatal(r.wrapError(err))
	}
	defer r.Cleanup()

	if t, ok := tc.(Setupper); ok {
		t.Setup(r)
		defer t.Cleanup()
	}

	switch tc.Method(){
	case http.MethodGet:
		if err := r.Get(tc.Path()); err != nil {
			t.Fatal(err)
		}

	case http.MethodPost:
		if err := r.Post(tc.Path(), tc.(PostTestCase).PostData()); err != nil {
			t.Fatal(err)
		}
	}

	f(t)

}

func ( r *testRunner)Check() error {

	t := r.testCase

	if t, ok := t.(StatusCodeChecker); ok {
		if err := r.CheckStatusCode(t.ExpectedStatusCode()); err != nil {
			return err
		}
	}

	if t, ok := t.(BodyChecker); ok {
		if err := r.CheckBody(t.ExpectedBody()); err != nil {
			return err
		}
	}

	if t, ok := t.(JsonBodyChecker); ok {
		if err := r.CheckBodyJson(t.ExpectedJsonSource()); err != nil {
			return err
		}
	}
	return nil
}

func ( r *testRunner)Get(path []string) error {
	var err error
	url := r.URL(path...)
	r.Logf("getting url: %q", url)
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrapf(err, "error getting url %q: %s", url)
	}
	r.saveResp(resp)
	return nil
}

func ( r *testRunner)Post(path []string, data interface{}) error {
	var err error
	url := r.URL(path...)

	json, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "can't marshal %#v", data)
	}
	r.Logf("sending json: %q to %q", json, url)

	resp, err := http.Post(url, "application/json", bytes.NewReader(json))
	if err != nil {
		return errors.Wrapf(err, "error getting url %q: %s", url)
	}
	r.saveResp(resp)
	return nil
}

func (r *testRunner)saveResp(resp *http.Response) error {
	defer resp.Body.Close()

	r.Resp = resp

	var err error

	r.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return r.wrapErrorf(err, "error reading response body")
	}
	return nil
}

func (r *testRunner)Logf(format string, args ...interface{}) {
	if !r.WriteLogs {
		return
	}
	log.Printf(format, args...)
}

func (r *testRunner)wrapErrorf(err error, message string, args ...interface{}) error {
	return r.wrapError(errors.Wrapf(err, message, args...))
}
func (r *testRunner)wrapError(err error) error {
	if t, ok := r.testCase.(ErrorWrapper); ok {
		return t.WrapError(err)
	}
	return err
}

func ( r *testRunner)CheckStatusCode(code int) error {
	if code != r.Resp.StatusCode {
		return r.wrapError(fmt.Errorf("status code is %d, expected %d", r.Resp.StatusCode, code))
	}
	return nil
}

func ( r *testRunner)CheckBody(expectedBody []byte) error {
	if !reflect.DeepEqual(expectedBody, r.Body) {
		return r.wrapError(fmt.Errorf("wrong response body:\ndata is:\t%q\nexpect:\t\t%q\n", r.Body, expectedBody))
	}
	return nil
}

func ( r *testRunner)CheckBodyJson(reference interface{}) error {
	expected, err := json.Marshal(reference)
	if err != nil {
		return r.wrapError(errors.Wrap(err, "can't marshal reference object"))
	}
	if !reflect.DeepEqual(expected, r.Body) {
		return r.wrapError(fmt.Errorf("wrong data recieved:\ndata is:\t%q\nexpect:\t\t%q\n", r.Body, expected))
	}
	return nil
}

func ( r *testRunner)URL(path ...string) string {
	urlPathParts := make([]string, 2, len(path) + 2)
	urlPathParts[0] = r.server.URL
	urlPathParts[1] = ApiPrefix
	if len(path) > 0 {
		for _, part := range path {
			if part != "" {
				part = strings.Trim(part, "/")
				urlPathParts = append(urlPathParts, part)
			}
		}
	}
	return strings.Join(urlPathParts, "/") + "/"
}
