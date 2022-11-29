package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8000/", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

}

/*
Expecting valid json response having all loaded websites as key and their respetive statuses as value
*/
func TestCheckAllSiteStatusHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8000/websites", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	var websites = make(map[string]string)

	err = json.NewDecoder(res.Body).Decode(&websites)
	assert.Equal(t, err, nil)
}

/*
Expecting valid json response with value as "DOWN" for website "cdbcb"
*/
func TestCheckSiteStatusHandler(t *testing.T) {

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8000/websites?name=cdbcb", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	var response = map[string]string{
		"cdbcb": "DOWN",
	}

	var websites = make(map[string]string)

	err = json.NewDecoder(res.Body).Decode(&websites)
	assert.Equal(t, response, websites)
	assert.Equal(t, err, nil)
}

/*
Request body should be a valid json not a empty string or nil
*/
func TestLoadWebsitesHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8000/websites?fff=cdbcb", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

/*
Sending valid json in request body
*/
func TestLoadWebsitesWithBodyHandler(t *testing.T) {
	var websites = map[string][]string{
		"websites": []string{"http://www.google.com", "http://www.facebook.com", "http://www.fakewebsite1.com"},
	}
	body, _ := json.Marshal(websites)
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8000/websites", bytes.NewBuffer([]byte(body)))
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
