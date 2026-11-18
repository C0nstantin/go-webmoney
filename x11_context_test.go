package webmoney

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchPassportResponseBoundaries(t *testing.T) {
	const valid = `<response retval="0"><certinfo wmid="123456789012"/></response>`
	for _, tt := range []struct {
		name, body string
		limit      int64
		wantErr    bool
	}{
		{"exact limit", valid, int64(len(valid)), false},
		{"valid prefix oversized", valid + strings.Repeat(" ", 100), int64(len(valid)), true},
		{"second root", valid + `<response retval="2"/>`, 1024, true},
		{"trailing text", valid + `garbage`, 1024, true},
		{"broken trailing comment", valid + `<!--`, 1024, true},
		{"legal trailing markup", valid + ` <!--ok--><?test ok?> `, 1024, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(tt.body)) }))
			defer srv.Close()
			_, err := FetchPassport(context.Background(), "123456789012", PassportRequestOptions{HTTPClient: srv.Client(), EndpointURL: srv.URL, MaxResponseBytes: tt.limit})
			if tt.wantErr != errors.Is(err, ErrPassportUnavailable) {
				t.Fatalf("error = %v, want unavailable %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchPassportCancellation(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(started); <-release }))
	defer srv.Close()
	defer close(release)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := FetchPassport(ctx, "123456789012", PassportRequestOptions{HTTPClient: srv.Client(), EndpointURL: srv.URL})
		done <- err
	}()
	<-started
	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) || !errors.Is(err, ErrPassportUnavailable) {
		t.Fatalf("error = %v", err)
	}
}

func TestFetchPassportCP1251AndCheckLock(t *testing.T) {
	// Actual Windows-1251 bytes for «Иван»; UTF-8 under this declaration must not pass unnoticed.
	body := `<?xml version="1.0" encoding="windows-1251"?><response retval="0"><fullaccess>0</fullaccess><certinfo wmid="123456789012"><attestat><row tid="130" recalled="1"/></attestat><userinfo><value><row fname="` + "\xc8\xe2\xe0\xed" + `"/></value><check-lock><row fname="10" nickname="00"/></check-lock></userinfo></certinfo></response>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(body)) }))
	defer srv.Close()
	got, err := FetchPassport(context.Background(), "123456789012", PassportRequestOptions{HTTPClient: srv.Client(), EndpointURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	if got.FullAccess != "0" || got.CertInfo.UserInfo.Fname != "Иван" || got.CertInfo.CheckLock.Fname != "10" || got.CertInfo.CheckLock.Nickname != "00" || got.CertInfo.Attestat.Recalled != "1" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestX11ContextAPIs(t *testing.T) {
	client, calls := passportServer(t, http.StatusOK, xmlOk)
	w := &WmClient{httpClient: client}
	ctx := context.Background()
	for _, run := range []func() error{
		func() error { _, e := CheckWmidContext(ctx, testWmid); return e },
		func() error { _, e := IssetWmidContext(ctx, testWmid); return e },
		func() error { _, e := GetInfoWmidContext(ctx, testWmid); return e },
		func() error { _, e := WMPassportByWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.CheckWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.IssetWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.GetInfoWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.WMPassportByWmidContext(ctx, testWmid); return e },
	} {
		if err := run(); err != nil {
			t.Fatal(err)
		}
	}
	if *calls != 8 {
		t.Fatalf("calls=%d, want 8", *calls)
	}
}

func TestX11ContextAPIsPropagateCancellation(t *testing.T) {
	client, calls := passportServer(t, http.StatusOK, xmlOk)
	w := &WmClient{httpClient: client}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, run := range []func() error{
		func() error { _, e := CheckWmidContext(ctx, testWmid); return e },
		func() error { _, e := IssetWmidContext(ctx, testWmid); return e },
		func() error { _, e := GetInfoWmidContext(ctx, testWmid); return e },
		func() error { _, e := WMPassportByWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.CheckWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.IssetWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.GetInfoWmidContext(ctx, testWmid); return e },
		func() error { _, e := w.WMPassportByWmidContext(ctx, testWmid); return e },
	} {
		if err := run(); !errors.Is(err, context.Canceled) || !errors.Is(err, ErrPassportUnavailable) {
			t.Fatalf("error = %v", err)
		}
	}
	if *calls != 0 {
		t.Fatalf("cancelled requests reached server: %d", *calls)
	}
}
