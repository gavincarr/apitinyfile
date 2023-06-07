package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	binary   = "apitinyfile"
	testuser = "test"
	testpass = "test"
)

func TestDefaults(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	env := Env{
		Read:      true,
		Write:     true,
		Delete:    true,
		Directory: dir,
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	env.setupRouter(r)
	data := "test data\n"
	b := strings.NewReader(data)

	// PUT
	req, _ := http.NewRequest("PUT", "/foo", b)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("PUT returned status %d, expected 204", w.Code)
		return
	}

	// GET
	req, _ = http.NewRequest("GET", "/foo", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("GET returned status %d, expected 200", w.Code)
		return
	}
	got := w.Body.String()
	if got != data {
		t.Errorf("GET returned %q, expected %q", got, data)
		return
	}

	// DELETE
	req, _ = http.NewRequest("DELETE", "/foo", b)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("DELETE returned status %d, expected 204", w.Code)
		return
	}

	// GET
	req, _ = http.NewRequest("GET", "/foo", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("GET returned status %d, expected 404", w.Code)
		return
	}
}

func TestAuth(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	env := Env{
		Read:      true,
		Write:     true,
		Delete:    true,
		Passwd:    "testdata/htpasswd",
		Directory: dir,
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	env.setupRouter(r)
	data := "different test data\n"
	b := strings.NewReader(data)

	// Unauth PUT
	req, _ := http.NewRequest("PUT", "/foo", b)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("PUT returned status %d, expected 401", w.Code)
		return
	}

	// Unauth GET
	req, _ = http.NewRequest("GET", "/foo", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("GET returned status %d, expected 401", w.Code)
		return
	}

	// Unauth DELETE
	req, _ = http.NewRequest("DELETE", "/foo", b)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("DELETE returned status %d, expected 401", w.Code)
		return
	}

	// Auth PUT
	req, _ = http.NewRequest("PUT", "/foo", b)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("PUT returned status %d, expected 204", w.Code)
		return
	}

	// Auth GET
	req, _ = http.NewRequest("GET", "/foo", nil)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("GET returned status %d, expected 200", w.Code)
		return
	}
	got := w.Body.String()
	if got != data {
		t.Errorf("GET returned %q, expected %q", got, data)
		return
	}

	// Auth DELETE
	req, _ = http.NewRequest("DELETE", "/foo", b)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("DELETE returned status %d, expected 204", w.Code)
		return
	}

	// Auth GET
	req, _ = http.NewRequest("GET", "/foo", nil)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("GET returned status %d, expected 404", w.Code)
		return
	}
}

func TestReadonly(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	env := Env{
		Read:      true,
		Write:     false,
		Delete:    false,
		Passwd:    "testdata/htpasswd",
		Directory: dir,
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	env.setupRouter(r)
	data := "different test data\n"
	b := strings.NewReader(data)

	// Auth GET - ok, expect a real 404
	req, _ := http.NewRequest("GET", "/foo", nil)
	req.SetBasicAuth(testuser, testpass)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("GET returned status %d, expected 404", w.Code)
		return
	}

	// Auth PUT - not ok, expect a no-route 404
	req, _ = http.NewRequest("PUT", "/foo", b)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("PUT returned status %d, expected 404", w.Code)
		return
	}

	// Auth DELETE - not ok, expect a no-route 404
	req, _ = http.NewRequest("DELETE", "/foo", b)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("DELETE returned status %d, expected 404", w.Code)
		return
	}

	// Write a test file to dir manually
	err := os.WriteFile(dir+"/foo", []byte(data), 0644)
	if err != nil {
		t.Errorf("error writing test file manually: %s", err.Error())
		return
	}

	// Auth GET - ok, expect content
	req, _ = http.NewRequest("GET", "/foo", nil)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("GET returned status %d, expected 200", w.Code)
		return
	}
	got := w.Body.String()
	if got != data {
		t.Errorf("GET returned %q, expected %q", got, data)
		return
	}
}

func TestWriteonly(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	env := Env{
		Read:      false,
		Write:     true,
		Delete:    true,
		Passwd:    "testdata/htpasswd",
		Directory: dir,
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	env.setupRouter(r)
	data := "different test data\n"
	b := strings.NewReader(data)

	// Auth PUT - ok
	req, _ := http.NewRequest("PUT", "/foo", b)
	req.SetBasicAuth(testuser, testpass)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("PUT returned status %d, expected 204", w.Code)
		return
	}

	// Auth GET - not ok, expect a no-route 404 (even though PUT was successful)
	req, _ = http.NewRequest("GET", "/foo", nil)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("GET returned status %d, expected 404", w.Code)
		return
	}

	// Auth DELETE - ok
	req, _ = http.NewRequest("DELETE", "/foo", b)
	req.SetBasicAuth(testuser, testpass)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("DELETE returned status %d, expected 204", w.Code)
		return
	}
}
