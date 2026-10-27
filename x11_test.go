package webmoney

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const testWmid = "251045363054"

// Ответы паспорта. Сняты с боевого сервиса (см. mail/newmail#72, #110).
const (
	// Живой аттестат.
	xmlOk = `<response retval="0"><certinfo wmid="251045363054">` +
		`<attestat><row tid="130" recalled="0" status="110"/></attestat></certinfo></response>`

	// Отзыв аттестата.
	xmlRecalled = `<response retval="0"><certinfo wmid="251045363054">` +
		`<attestat><row tid="130" recalled="1" status="110"/></attestat></certinfo></response>`

	// WMID, про который паспорт ничего не знает: выдуманный и удалённый
	// отвечают побайтово одинаково. certinfo.wmid — эхо запроса.
	xmlUnknownWmid = `<?xml version="1.0" encoding="windows-1251"?><response retval="0">` +
		`<fullaccess>0</fullaccess><certinfo wmid="251045363054"></certinfo></response>`

	// Паспорт знает WMID, аттестата нет.
	xmlNoAttestat = `<response retval="0"><certinfo wmid="251045363054">` +
		`<wmids><row wmid="251045363054" datereg="2013-05-16T09:25:18"/></wmids></certinfo></response>`

	// Аварийные ответы.
	xmlNoRetval    = `<response><certinfo wmid="251045363054"></certinfo></response>`
	xmlEmptyRetval = `<response retval=""><certinfo wmid="251045363054"></certinfo></response>`
	xmlBadRetval   = `<response retval="ok"><certinfo wmid="251045363054"></certinfo></response>`
	xmlErrRetval   = `<response retval="-3" retdesc="bad wmid"></response>`
	// Снято живьём 2026-10-28: passportwmid="abc", "12345" и "" дают ровно это.
	xmlInvalidWmid  = `<?xml version="1.0" encoding="windows-1251"?><response retval="2" retdesc="неверно указан проверяемый WMID"></response>`
	xmlSpaceRetval  = `<response retval=" 0 "><certinfo wmid="251045363054"></certinfo></response>`
	xmlNoWmidAttr   = `<response retval="0"><certinfo></certinfo></response>`
	xmlNoCertinfo   = `<response retval="0"></response>`
	xmlOtherWmid    = `<response retval="0"><certinfo wmid="000000000000"></certinfo></response>`
	htmlBody        = `<html><body>Service Unavailable</body></html>`
	xmlCp1251Prolog = `<?xml version="1.0" encoding="windows-1251"?>` + xmlOk
)

// rewriteTransport направляет запросы к apipassport.web.money на тестовый сервер.
type rewriteTransport struct {
	base   http.RoundTripper
	target *url.URL
}

func (t rewriteTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r2 := r.Clone(r.Context())
	r2.URL.Scheme = t.target.Scheme
	r2.URL.Host = t.target.Host
	return t.base.RoundTrip(r2)
}

// passportServer поднимает фейковый паспорт и подменяет package-level клиент.
// handler получает тело запроса и возвращает статус и тело ответа.
func passportServer(t *testing.T, status int, body string) (*http.Client, *int32) {
	t.Helper()
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		if r.Method != http.MethodPost {
			t.Errorf("метод %s, ждали POST", r.Method)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	client := &http.Client{
		Transport: rewriteTransport{base: http.DefaultTransport, target: u},
		Timeout:   5 * time.Second,
	}
	SetPassportHTTPClient(client)
	t.Cleanup(func() { SetPassportHTTPClient(nil) })
	return client, &calls
}

func TestCheckWmid(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		body    string
		want    WmidStatus
		wantErr error // nil — ошибки быть не должно
	}{
		{"живой аттестат — существует", 200, xmlOk, WmidExists, nil},
		{"отзыв — WMID всё равно существует", 200, xmlRecalled, WmidExists, nil},
		{"есть wmids, нет аттестата — существует", 200, xmlNoAttestat, WmidExists, nil},
		{"пустой certinfo — не найден", 200, xmlUnknownWmid, WmidNotFound, nil},
		{"пролог windows-1251 разбирается", 200, xmlCp1251Prolog, WmidExists, nil},
		{"HTTP 500 — неизвестно", 500, xmlOk, WmidUnknown, ErrPassportUnavailable},
		{"HTTP 302 без Location — неизвестно", 302, "", WmidUnknown, ErrPassportUnavailable},
		{"retval=2 (невалидный WMID) — не найден", 200, xmlInvalidWmid, WmidNotFound, nil},
		{"retval с пробелами — число", 200, xmlSpaceRetval, WmidNotFound, nil},
		{"certinfo без атрибута wmid — неизвестно", 200, xmlNoWmidAttr, WmidUnknown, ErrPassportUnavailable},
		{"HTML вместо XML — неизвестно", 200, htmlBody, WmidUnknown, ErrPassportUnavailable},
		{"пустое тело — неизвестно", 200, "", WmidUnknown, ErrPassportUnavailable},
		{"нет атрибута retval — неизвестно", 200, xmlNoRetval, WmidUnknown, ErrPassportUnavailable},
		{"пустой retval — неизвестно", 200, xmlEmptyRetval, WmidUnknown, ErrPassportUnavailable},
		{"retval не число — неизвестно", 200, xmlBadRetval, WmidUnknown, ErrPassportUnavailable},
		{"retval=0 без certinfo — неизвестно", 200, xmlNoCertinfo, WmidUnknown, ErrPassportUnavailable},
		{"certinfo про другой WMID — неизвестно", 200, xmlOtherWmid, WmidUnknown, ErrPassportUnavailable},
		{"retval != 0 — неизвестно, ErrPassportRetval", 200, xmlErrRetval, WmidUnknown, ErrPassportRetval},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			passportServer(t, c.status, c.body)
			got, err := CheckWmid(testWmid)
			if got != c.want {
				t.Errorf("статус %v, ждали %v (err=%v)", got, c.want, err)
			}
			if c.wantErr == nil && err != nil {
				t.Errorf("неожиданная ошибка: %v", err)
			}
			if c.wantErr != nil && !errors.Is(err, c.wantErr) {
				t.Errorf("ошибка %v, ждали %v", err, c.wantErr)
			}
		})
	}
}

func TestCheckWmid_RetvalErrorCarriesDetails(t *testing.T) {
	passportServer(t, 200, xmlErrRetval)
	_, err := CheckWmid(testWmid)
	if err == nil || !strings.Contains(err.Error(), "retval=-3") || !strings.Contains(err.Error(), "bad wmid") {
		t.Errorf("в ошибке нет retval/retdesc: %v", err)
	}
	// GetInfoWmid при retval != 0 отдаёт сам ответ: потребители читают Retval/Retdesc.
	v, err := GetInfoWmid(testWmid)
	if !errors.Is(err, ErrPassportRetval) || v.Retval != -3 || v.Retdesc != "bad wmid" {
		t.Errorf("v=%+v err=%v", v, err)
	}
}

func TestRetvalErrorTruncatesRetdesc(t *testing.T) {
	long := strings.Repeat("x", 5000)
	passportServer(t, 200, `<response retval="-3" retdesc="`+long+`"></response>`)
	_, err := CheckWmid(testWmid)
	if err == nil || len(err.Error()) > 400 {
		t.Errorf("retdesc не обрезан: len=%d", len(err.Error()))
	}
}

func TestPassportRequestShape(t *testing.T) {
	var gotPath, gotCT, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotCT = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(xmlOk))
	}))
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	SetPassportHTTPClient(&http.Client{Transport: rewriteTransport{base: http.DefaultTransport, target: u}})
	t.Cleanup(func() { SetPassportHTTPClient(nil) })

	if _, err := CheckWmid(testWmid); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/asp/XMLGetWMPassport.asp" {
		t.Errorf("path %q", gotPath)
	}
	if gotCT != "text/xml" {
		t.Errorf("content-type %q", gotCT)
	}
	if !strings.Contains(gotBody, "<passportwmid>"+testWmid+"</passportwmid>") {
		t.Errorf("в теле нет passportwmid: %q", gotBody)
	}
}

func TestPassportDomainEnv(t *testing.T) {
	t.Setenv("API_PASSPORT_DOMAIN", "passport.example.test")
	if got := getPassportURL(); got != "https://passport.example.test/asp/XMLGetWMPassport.asp" {
		t.Errorf("url %q", got)
	}
	t.Setenv("API_PASSPORT_DOMAIN", "evil.test/path?x=1")
	if got := getPassportURL(); got != "https://apipassport.web.money/asp/XMLGetWMPassport.asp" {
		t.Errorf("невалидный домен должен игнорироваться: %q", got)
	}
}

func TestPassportRedirectNotFollowed(t *testing.T) {
	var hits int32
	mux := http.NewServeMux()
	mux.HandleFunc("/asp/XMLGetWMPassport.asp", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		http.Redirect(w, r, "/elsewhere", http.StatusFound)
	})
	mux.HandleFunc("/elsewhere", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 100)
		_, _ = w.Write([]byte(xmlOk))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)

	// Клиент по умолчанию: редирект не выполняется, 302 — сбой.
	SetPassportHTTPClient(nil)
	SetPassportTimeout(0)
	def := getPassportClient()
	def.Transport = rewriteTransport{base: http.DefaultTransport, target: u}
	t.Cleanup(func() { SetPassportHTTPClient(nil) })

	st, err := CheckWmid(testWmid)
	if st != WmidUnknown || !errors.Is(err, ErrPassportUnavailable) || !strings.Contains(err.Error(), "302") {
		t.Errorf("st=%v err=%v", st, err)
	}
	if atomic.LoadInt32(&hits) != 1 {
		t.Errorf("редирект был выполнен: hits=%d", hits)
	}
}

func TestPassportOversizedBody(t *testing.T) {
	big := `<response retval="0"><certinfo wmid="` + testWmid + `"><attestat><row tid="130" recalled="0" ` +
		`notary="` + strings.Repeat("a", maxPassportResponseSize) + `"/></attestat></certinfo></response>`
	passportServer(t, 200, big)
	st, err := CheckWmid(testWmid)
	if st != WmidUnknown || !errors.Is(err, ErrPassportUnavailable) {
		t.Errorf("усечённый ответ должен быть сбоем: st=%v err=%v", st, err)
	}
}

func TestSetPassportTimeoutDoesNotTouchCustomClient(t *testing.T) {
	custom := &http.Client{Timeout: 42 * time.Second}
	SetPassportHTTPClient(custom)
	t.Cleanup(func() { SetPassportHTTPClient(nil); SetPassportTimeout(0) })

	SetPassportTimeout(3 * time.Second)
	if getPassportClient() != custom || custom.Timeout != 42*time.Second {
		t.Error("SetPassportTimeout не должен трогать пользовательский клиент")
	}
	SetPassportHTTPClient(nil)
	if got := getPassportClient().Timeout; got != 3*time.Second {
		t.Errorf("после сброса клиента должен применяться сохранённый таймаут, got %v", got)
	}
}

func TestPassportSettingsConcurrency(t *testing.T) {
	passportServer(t, 200, xmlOk)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = CheckWmid(testWmid)
		}()
		go func(i int) {
			defer wg.Done()
			SetPassportTimeout(time.Duration(i+1) * time.Second)
		}(i)
	}
	wg.Wait()
	SetPassportTimeout(0)
}

func TestIssetWmid(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		body    string
		want    bool
		wantErr bool
	}{
		{"существует", 200, xmlOk, true, false},
		{"не найден — false без ошибки", 200, xmlUnknownWmid, false, false},
		{"невалидный WMID — false без ошибки", 200, xmlInvalidWmid, false, false},
		{"сбой — false с ошибкой", 500, "", false, true},
		{"retval != 0 — false с ошибкой", 200, xmlErrRetval, false, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			passportServer(t, c.status, c.body)
			got, err := IssetWmid(testWmid)
			if got != c.want {
				t.Errorf("получили %v, ждали %v", got, c.want)
			}
			if (err != nil) != c.wantErr {
				t.Errorf("err=%v, ждали ошибку: %v", err, c.wantErr)
			}
		})
	}
}

func TestGetInfoWmid(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		body    string
		wantErr error
		wantTID string
	}{
		{"живой аттестат", 200, xmlOk, nil, "130"},
		{"отзыв", 200, xmlRecalled, ErrAttestatRecalled, "130"},
		{"есть wmids, нет аттестата", 200, xmlNoAttestat, ErrAttestatNotFound, ""},
		{"WMID не найден — по-прежнему attestat not found", 200, xmlUnknownWmid, ErrAttestatNotFound, ""},
		{"HTTP 500", 500, xmlOk, ErrPassportUnavailable, ""},
		{"нет retval", 200, xmlNoRetval, ErrPassportUnavailable, ""},
		{"retval != 0 — не «аттестат не найден»", 200, xmlErrRetval, ErrPassportRetval, ""},
		{"retval=2 — по-прежнему attestat not found (api сверяет текст)", 200, xmlInvalidWmid, ErrAttestatNotFound, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			passportServer(t, c.status, c.body)
			got, err := GetInfoWmid(testWmid)
			if c.wantErr == nil && err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if c.wantErr != nil && !errors.Is(err, c.wantErr) {
				t.Fatalf("ошибка %v, ждали %v", err, c.wantErr)
			}
			if got.CertInfo.Attestat.TID != c.wantTID {
				t.Errorf("TID %q, ждали %q", got.CertInfo.Attestat.TID, c.wantTID)
			}
		})
	}
}

// Тексты сентинелов сверяет api (services/user.go: isWmidRejected) — менять нельзя.
func TestAttestatErrorTexts(t *testing.T) {
	if ErrAttestatNotFound.Error() != "attestat not found" {
		t.Errorf("ErrAttestatNotFound = %q", ErrAttestatNotFound.Error())
	}
	if ErrAttestatRecalled.Error() != "attestat recalled" {
		t.Errorf("ErrAttestatRecalled = %q", ErrAttestatRecalled.Error())
	}
	// Для WMID, которого нет, GetInfoWmid возвращает ровно сентинел без обёртки.
	passportServer(t, 200, xmlUnknownWmid)
	_, err := GetInfoWmid(testWmid)
	if err == nil || err.Error() != "attestat not found" {
		t.Errorf("err = %v, ждали ровно «attestat not found»", err)
	}
}

func TestWMPassportByWmid(t *testing.T) {
	t.Run("retval != 0 возвращается как есть", func(t *testing.T) {
		passportServer(t, 200, xmlErrRetval)
		got, err := WMPassportByWmid(testWmid)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if got.Retval != -3 || got.Retdesc != "bad wmid" {
			t.Errorf("Retval=%d Retdesc=%q", got.Retval, got.Retdesc)
		}
	})
	t.Run("живой ответ", func(t *testing.T) {
		passportServer(t, 200, xmlOk)
		got, err := WMPassportByWmid(testWmid)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if got.Retval != 0 || got.CertInfo.Wmid != testWmid || got.CertInfo.Attestat.TID != "130" {
			t.Errorf("разобрано неверно: %+v", got)
		}
	})
	t.Run("нет retval — ошибка, не нулевой ответ", func(t *testing.T) {
		passportServer(t, 200, xmlNoRetval)
		got, err := WMPassportByWmid(testWmid)
		if !errors.Is(err, ErrPassportUnavailable) || got != nil {
			t.Errorf("got=%v err=%v", got, err)
		}
	})
	t.Run("HTTP 503", func(t *testing.T) {
		passportServer(t, 503, xmlOk)
		_, err := WMPassportByWmid(testWmid)
		if !errors.Is(err, ErrPassportUnavailable) || !strings.Contains(err.Error(), "503") {
			t.Errorf("err=%v", err)
		}
	})
}

func TestPassportTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	}))
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	SetPassportHTTPClient(&http.Client{
		Transport: rewriteTransport{base: http.DefaultTransport, target: u},
		Timeout:   100 * time.Millisecond,
	})
	t.Cleanup(func() { SetPassportHTTPClient(nil) })

	start := time.Now()
	st, err := CheckWmid(testWmid)
	if st != WmidUnknown || !errors.Is(err, ErrPassportUnavailable) {
		t.Errorf("st=%v err=%v", st, err)
	}
	// Причина доступна через Unwrap: вызывающий может отличить таймаут.
	var ne net.Error
	if !errors.As(err, &ne) || !ne.Timeout() {
		t.Errorf("ошибка не разворачивается до net.Error с Timeout(): %v", err)
	}
	if time.Since(start) > time.Second {
		t.Errorf("таймаут не сработал: %v", time.Since(start))
	}
}

func TestDefaultPassportClientHasTimeout(t *testing.T) {
	SetPassportHTTPClient(nil)
	SetPassportTimeout(0)
	t.Cleanup(func() { SetPassportTimeout(0) })

	t.Run("по умолчанию 30s", func(t *testing.T) {
		unsetenvForTest(t, "API_PASSPORT_TIMEOUT")
		SetPassportTimeout(0)
		if got := getPassportClient().Timeout; got != defaultHTTPTimeout {
			t.Errorf("Timeout=%v, ждали %v", got, defaultHTTPTimeout)
		}
	})
	t.Run("из API_PASSPORT_TIMEOUT", func(t *testing.T) {
		t.Setenv("API_PASSPORT_TIMEOUT", "7s")
		SetPassportTimeout(0)
		if got := getPassportClient().Timeout; got != 7*time.Second {
			t.Errorf("Timeout=%v, ждали 7s", got)
		}
	})
	t.Run("кривой API_PASSPORT_TIMEOUT — по умолчанию", func(t *testing.T) {
		t.Setenv("API_PASSPORT_TIMEOUT", "fast")
		SetPassportTimeout(0)
		if got := getPassportClient().Timeout; got != defaultHTTPTimeout {
			t.Errorf("Timeout=%v, ждали %v", got, defaultHTTPTimeout)
		}
	})
	t.Run("SetPassportTimeout перекрывает env", func(t *testing.T) {
		t.Setenv("API_PASSPORT_TIMEOUT", "7s")
		SetPassportTimeout(3 * time.Second)
		if got := getPassportClient().Timeout; got != 3*time.Second {
			t.Errorf("Timeout=%v, ждали 3s", got)
		}
	})
	t.Run("клиент кэшируется", func(t *testing.T) {
		SetPassportTimeout(0)
		if getPassportClient() != getPassportClient() {
			t.Error("getPassportClient должен возвращать один экземпляр")
		}
	})
}

func TestWmClientX11(t *testing.T) {
	client, calls := passportServer(t, 200, xmlOk)
	SetPassportHTTPClient(nil) // WmClient не должен зависеть от package-level клиента
	w := &WmClient{Wmid: "000000000000", httpClient: client}

	st, err := w.CheckWmid(testWmid)
	if st != WmidExists || err != nil {
		t.Errorf("CheckWmid: %v %v", st, err)
	}
	ok, err := w.IssetWmid(testWmid)
	if !ok || err != nil {
		t.Errorf("IssetWmid: %v %v", ok, err)
	}
	info, err := w.GetInfoWmid(testWmid)
	if err != nil || info.CertInfo.Attestat.TID != "130" {
		t.Errorf("GetInfoWmid: %+v %v", info, err)
	}
	raw, err := w.WMPassportByWmid(testWmid)
	if err != nil || raw.Retval != 0 {
		t.Errorf("WMPassportByWmid: %+v %v", raw, err)
	}
	if atomic.LoadInt32(calls) != 4 {
		t.Errorf("запросов к серверу %d, ждали 4", *calls)
	}
}

func TestClassifyWmid(t *testing.T) {
	var v ResponseX11
	if classifyWmid(v) != WmidNotFound {
		t.Error("пустой ответ с retval 0 — not found")
	}
	v.CertInfo.Attestat.Datecrt = "2020-01-01"
	if classifyWmid(v) != WmidExists {
		t.Error("datecrt без tid — существует")
	}
	v = ResponseX11{Retval: 1}
	if classifyWmid(v) != WmidUnknown {
		t.Error("retval != 0 — unknown")
	}
	if WmidExists.String() != "exists" || WmidNotFound.String() != "not found" || WmidUnknown.String() != "unknown" {
		t.Error("String()")
	}
}

// unsetenvForTest снимает переменную окружения и возвращает её по завершении теста.
func unsetenvForTest(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "") // регистрирует восстановление прежнего значения
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}
}
