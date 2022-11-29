package main

import (
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

func TestCheckAllSiteStatusHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8000/websites", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	var websites = make(map[string]string)

	err = json.NewDecoder(res.Body).Decode(&websites)
	assert.Equal(t, err, nil)
}

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

func TestLoadWebsitesHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8000/websites?fff=cdbcb", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
