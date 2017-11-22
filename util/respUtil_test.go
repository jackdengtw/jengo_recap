package util

import (
	_ "errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
)

//
// mock http writer
//
type mockWriter struct {
	mockHeader http.Header
}

func (m *mockWriter) Header() http.Header {
	return m.mockHeader
}

func (m *mockWriter) Write(bytes []byte) (int, error) {
	log.Printf("Client will see body  : %v\n", string(bytes))
	return 0, nil
}

func (m *mockWriter) WriteHeader(code int) {
	log.Printf("Client will see status: [status code %v]\n", code)
}

var mw *mockWriter = &mockWriter{
	mockHeader: make(http.Header),
}

// mock logger
func mockErrorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Return Type
type commonResponse struct {
	Id  string
	Msg string
}

type realObj struct {
	Identity string
	Property string
}

//
// Example Unit Tests to show use cases
//
// Note,
// 1. Default http body is in json/utf8 format
// 2. Return Utility itself can fail, invoker SHOULD log failure event
//

// 500 Internal Server Error
// When,
// 1. a server side error leads to failed response,
// Note,
// 1. internal error detail won't return to http client
// 2.
func TestReturnError500(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnError500")
	generalMsg := "my func cannot open file."
	ReturnError(
		500,
		generalMsg,
		mw)

	os.Stderr.Sync()
}

func TestReturnErrorWithObj500(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnErrorWithObj500")
	generalMsg := "my func cannot open file."
	ReturnErrorWithObj(
		500,
		commonResponse{Id: "uuid1", Msg: generalMsg},
		mw)

	os.Stderr.Sync()
}

// 400 Bad Request
// When,
// any path/params/header error
func TestReturnError400(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnError400")
	generalMsg := "Param project_id in querystring is missing"
	ReturnError(
		400,
		generalMsg,
		mw)

	os.Stderr.Sync()
}

// 404 Not Found
// when,
// GET/DELETE non-existed ID:
func TestReturnError404(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnError404")
	generalMsg := "Id 123456 is NOT found."
	ReturnError(
		404,
		generalMsg,
		mw)

	os.Stderr.Sync()
}

// 405 Method Not Allowed
// When,
// PUT to anything only allow get
func TestReturnError405(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnError405")
	generalMsg := "Put /v0.1/run/run_id is NOT allowed."
	ReturnError(
		405,
		generalMsg,
		mw)

	os.Stderr.Sync()
}

// 409 Conflict
// When,
// Creating an existed user
func TestReturnError409(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnError409")
	generalMsg := "User alread existed."
	ReturnError(
		409,
		generalMsg,
		mw)

	os.Stderr.Sync()
}

// GET 200 OK
func TestReturnSuccessGet200(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnSuccessGet200")
	ReturnSuccess(
		200,
		fmt.Sprintf(
			`{
				"id": "I would be object requested", 
			 }`),
		"",
		mw)

	os.Stderr.Sync()
}

func TestReturnSuccessWithObjAndHeadersGet200(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnSuccessWithObjAndHeadersGet200")
	ReturnSuccessWithObjAndHeaders(
		200,
		realObj{Identity: "identity1", Property: "property1"},
		map[string]string{"content-encoding": "gzip"},
		mw)

	os.Stderr.Sync()
}

// PUT 200 OK
func TestReturnSuccessPut200(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnSuccessPut200")
	ReturnSuccessWithObj(
		200,
		CommonResponse{Url: "http://host:port/some_entity/some_id"},
		mw)

	os.Stderr.Sync()
}

// 201 Created
// Only POST to create successfully return this code
func TestReturnSuccess201(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnSuccess201")
	ReturnSuccessWithObj(
		201,
		CommonResponse{Url: "http://host:port/some_entity/some_id"},
		mw)

	os.Stderr.Sync()
}

// CALLBACK 202 Accepted
// Reserved for possible async call usage

// DELETE 204 No Content
func TestReturnSuccess204(t *testing.T) {
	origin := logger.Errorf
	defer func() { logger.Errorf = origin }()
	logger.Errorf = mockErrorf

	log.Println("\nTestReturnSuccess204")
	ReturnSuccess(
		204,
		"whatever",
		"whatever",
		mw)

	os.Stderr.Sync()
}
