package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

const page_not_found = "404 page not found\n"

func testRequest(t *testing.T, method, path, data, expbody string, code int) {
	req := httptest.NewRequest(method, "http://localhost:9999/"+path, strings.NewReader(data))
	w := httptest.NewRecorder()
	HandleData(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != code {
		t.Errorf("%s %s - Expected status code %d isntead of %d", method, path, code, resp.StatusCode)
	}

	if string(body) != expbody {
		t.Errorf("%s %s - Expected %q as body instead of %q", method, path, expbody, string(body))
	}
}

func TestServer(t *testing.T) {

	testRequest(t, "GET", "nonexistent", "", page_not_found, 404)

	testRequest(t, "POST", "foo", "bar", "", 200)

	testRequest(t, "GET", "foo", "", "bar", 200)

	testRequest(t, "GET", "", "", "foo\n", 200)
	testRequest(t, "GET", "/", "", "foo\n", 200)

	testRequest(t, "DELETE", "foo", "", "bar", 200)
	testRequest(t, "DELETE", "foo", "", page_not_found, 404)

}
