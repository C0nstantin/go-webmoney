package webmoney

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Регрессия: sendRequest парсил тело при любом HTTP-статусе, маскируя
// ошибки сервера/прокси под «некорректный XML».

func TestSendRequestRejectsNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<error><secret_key>topsecret` + strings.Repeat("x", 220) + `</secret_key></error>`))
	}))
	defer srv.Close()

	c := &WmClient{}
	_, err := c.sendRequest(srv.URL, "<request/>")
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error %q does not mention status 500", err)
	}
	if strings.Contains(err.Error(), "topsecret") {
		t.Errorf("error leaks secret from response body: %q", err)
	}
}

func TestSendRequestAccepts200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<response retval="0"/>`))
	}))
	defer srv.Close()

	c := &WmClient{}
	got, err := c.sendRequest(srv.URL, "<request/>")
	if err != nil {
		t.Fatal(err)
	}
	if got != `<response retval="0"/>` {
		t.Errorf("body = %q", got)
	}
}
