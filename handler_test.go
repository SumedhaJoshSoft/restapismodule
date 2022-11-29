package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type HttpHandlerMock struct {
	mock.Mock
}

func (m *HttpHandlerMock) requestHandler(method string, request string, url string, handle http.HandlerFunc, body string) (rr *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(method, request, bytes.NewBuffer([]byte(body)))

	rr = httptest.NewRecorder()
	router := mux.NewRouter()

	router.HandleFunc(url, handle).Methods(method)
	router.ServeHTTP(rr, req)
	return
}

func TestDefaultHandler(t *testing.T) {
	mockHandler := &HttpHandlerMock{}
	res := mockHandler.requestHandler(http.MethodGet, "/", "/", defaultHandler, "")
	assert.Equal(t, res.Code, http.StatusOK)

}

/*
Expecting valid json response having all loaded websites as key and their respetive statuses as value
Sample response:
		{
			"http://www.facebook.com": "UP",
			"http://www.fakewebsite1.com": "DOWN",
			"http://www.google.com": "UP"
		}
*/
func TestCheckAllSiteStatusHandler(t *testing.T) {
	mockHandler := &HttpHandlerMock{}
	res := mockHandler.requestHandler(http.MethodGet, "/websites", "/websites", checkSiteStatusHandler, "")
	assert.Equal(t, res.Code, http.StatusOK)

	var websites = make(map[string]string)
	err := json.NewDecoder(res.Body).Decode(&websites)
	assert.Equal(t, err, nil)
}

/*
Expecting valid json response with value as "DOWN" for website "cdbcb"
Expected Response:
		{
			"cdbcb": "DOWN"
		}
*/
func TestCheckSiteStatusHandler(t *testing.T) {

	mockHandler := &HttpHandlerMock{}
	res := mockHandler.requestHandler(http.MethodGet, "/websites?name=cdbcb", "/websites", checkSiteStatusHandler, "")
	assert.Equal(t, res.Code, http.StatusOK)

	var response = map[string]string{
		"cdbcb": "DOWN",
	}
	var websites = make(map[string]string)

	err := json.NewDecoder(res.Body).Decode(&websites)
	assert.Equal(t, response, websites)
	assert.Equal(t, err, nil)
}

/*
Request body should be a valid json not a empty string or nil
Sample Input:
		{"websites":["http://www.google.com","http://www.facebook.com","http://www.fakewebsite1.com"]}

*/
func TestLoadWebsitesHandler(t *testing.T) {
	mockHandler := &HttpHandlerMock{}
	res := mockHandler.requestHandler(http.MethodPost, "/websites", "/websites", loadWebsitesHandler, "")
	assert.Equal(t, res.Code, http.StatusBadRequest)
}

/*
Sending valid json in request body
Sample Input:
		{"websites":["http://www.google.com","http://www.facebook.com","http://www.fakewebsite1.com"]}
*/
func TestLoadWebsitesWithBodyHandler(t *testing.T) {
	var websites = `{"websites":["http://www.google.com","http://www.facebook.com","http://www.fakewebsite1.com"]}`

	mockHandler := &HttpHandlerMock{}
	res := mockHandler.requestHandler(http.MethodPost, "/websites", "/websites", loadWebsitesHandler, websites)
	assert.Equal(t, res.Code, http.StatusOK)
}
