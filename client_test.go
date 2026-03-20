package webmoney

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestSendRequestDoesNotCorruptGlobalClient(t *testing.T) {
	// Сохраняем исходное состояние глобального транспорта
	originalTransport := http.DefaultClient.Transport

	// Включаем использование сертификатов для активации логики создания транспорта
	os.Setenv("USE_W3S_CERT", "1")
	defer os.Unsetenv("USE_W3S_CERT")

	// Имитируем настройки Classic (WinPro) клиента
	w := &WmClient{
		Wmid: "000000000000",
		Key:  "dGVzdF9rZXk=", // валидный base64 для прохождения проверок
		Pass: "test_pass",
	}

	// Выполняем запрос к невалидному адресу. 
	// Нам не важен результат самого запроса, нам важно поведение внутреннего клиента.
	_, _ = w.sendRequest("http://127.0.0.1:0", "<xml/>")

	// Проверяем, что глобальный транспорт не изменился
	if http.DefaultClient.Transport != originalTransport {
		t.Errorf("Security Breach: Global http.DefaultClient.Transport was corrupted! Got %v, expected %v",
			http.DefaultClient.Transport, originalTransport)
	}
}

func TestGetHTTPClientCachesClient(t *testing.T) {
	os.Unsetenv("USE_W3S_CERT")

	w := &WmClient{
		Wmid: "000000000000",
		Key:  "dGVzdF9rZXk=",
		Pass: "test_pass",
	}

	client1, err := w.getHTTPClient()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	client2, err := w.getHTTPClient()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if client1 != client2 {
		t.Error("getHTTPClient() should return the same cached instance on subsequent calls")
	}
}

func TestGetHTTPClientHasTimeout(t *testing.T) {
	os.Unsetenv("USE_W3S_CERT")

	w := &WmClient{
		Wmid: "000000000000",
		Key:  "dGVzdF9rZXk=",
		Pass: "test_pass",
	}

	client, err := w.getHTTPClient()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if client.Timeout == 0 {
		t.Error("HTTP client should have a non-zero timeout")
	}

	if client.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.Timeout)
	}
}
