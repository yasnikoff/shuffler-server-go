package api

// TODO: refactor all api tests to use new (api2_test) system

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"fmt"
	"reflect"
	"strings"
	"encoding/json"
	"bytes"
	"bitbucket.org/yasnikoff/shuffler/app/app_fixtures"
	"bitbucket.org/yasnikoff/shuffler/app"
	"net/url"
	"github.com/pkg/errors"
	"log"
)

type apiTestFx struct {
	*app_fixtures.AppFixture
	CaseName           string
	Expected           interface{}
	ExpectedReturnCode int
	SubPath            string
	URL                string
	Resp               *http.Response
	Body               []byte
	app                app.App
	handler            *handler
	server             *httptest.Server
	appName            string
	writeLogs          bool

	currentCase        TestCase
}

func NewRunner(t *testing.T, name string) *apiTestFx {
	r := apiTestFx{
		AppFixture:app_fixtures.New(t),
		appName:name,
		handler:NewHandler(),
		ExpectedReturnCode:http.StatusOK,

	}
	err := r.handler.registerApp(r.App, r.appName)
	if err != nil {
		panic(err)
	}
	r.server = httptest.NewServer(r.handler)
	//fmt.Printf("Listening on %q\n", r.server.URL)
	return &r
}

func (r *apiTestFx)Cleanup() {
	r.AppFixture.Cleanup()
	r.server.CloseClientConnections()
	r.server.Close()
}

func (r *apiTestFx)Errorf(msg string, args ...interface{}) {
	r.T.Errorf("%q: path %q, case %q: %s", r.appName, r.subPathUrl(r.SubPath), strings.ToUpper(r.CaseName), fmt.Sprintf(msg, args...))
}

func (r *apiTestFx)Logf(format string, args ...interface{}) {
	if !r.writeLogs {
		return
	}

	log.Printf(format, args...)
}

// Check causes test to fail if expected criteria are not met.
func (r *apiTestFx)Check() {

	expected := bytes.NewBuffer([]byte{})

	if r.Expected != nil {
		err := json.NewEncoder(expected).Encode(r.Expected)
		if err != nil {
			r.Errorf("error marshalling test.Expected: %s", err)
			return
		}
	}

	if r.Resp.StatusCode != r.ExpectedReturnCode {
		r.Errorf("status code is %d, expected %d", r.Resp.StatusCode, r.ExpectedReturnCode)
	}

	if !reflect.DeepEqual(expected.Bytes(), r.Body) {
		r.Errorf("wrong data recieved:\ndata is:\t%q\nexpect:\t\t%q\n", r.Body, expected)
	}
}

// subPathUrl returns full url string with path input constructed from subPath string slice.
func (r *apiTestFx)subPathUrl(subPath ...string) string {

	parts := []string{r.server.URL, ApiPrefix, r.appName}
	if len(subPath) > 0 {
		parts = append(parts, subPath...)
	}
	u := url.URL{
		Path:toPath(parts...),
	}
	return strings.TrimSpace(u.EscapedPath())
}

// function toPath returns path url input constructed from string slice
func toPath(parts ...string) string {
	if len(parts) > 0 {
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if part != "" {
				result = append(result, part)
			}
		}
		return strings.TrimSuffix(strings.Join(result, "/"), "/") + "/"
	}
	return ""
}

func (r *apiTestFx)Get(subPath string) {
	r.URL = r.subPathUrl(subPath)
	r.SubPath = subPath
	var err error
	//fmt.Printf("Getting url %q\n", r.URL)
	r.Resp, err = http.Get(r.URL)
	if err != nil {
		r.Fatalf("error getting url %q: %s", r.URL, err)

	}
	if r.Resp == nil {
		r.Fatalf("nil response from %q", r.URL)
	}
	r.Body = r.ReadBody(r.Resp)
}

func (r *apiTestFx)Post(subPath string, data interface{}) {
	r.URL = r.subPathUrl(subPath)
	r.SubPath = subPath

	json, err := json.Marshal(data)
	if err != nil {
		r.Fatal(errors.Wrapf(err, "can't marshal %#v", data))
	}

	r.Logf("sending json: %q", json)
	r.Resp, err = http.Post(r.URL, "application/json", bytes.NewReader(json))
	if err != nil {
		r.Fatalf("error getting url %q: %s", r.URL, err)

	}
	if r.Resp == nil {
		r.Fatalf("nil response from %q", r.URL)
	}
	r.Body = r.ReadBody(r.Resp)
}

func (r *apiTestFx)ReadBody(resp *http.Response) []byte {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.T.Fatalf("error reading response body: %s", err)
	}
	defer resp.Body.Close()

	return body

}

func (r *apiTestFx)RunGet(subPathParts ...string) {
	r.Get(toPath(subPathParts...))
	r.Check()
}

