package wmsigner

import (
	"encoding/base64"
	"testing"
)

func TestSignWithInvalidKeyDoesNotPanic(t *testing.T) {
	// Создаем "ключ" нужной длины (164 байта), но с невалидным содержимым
	invalidKey := make([]byte, 164)
	keyStr := base64.StdEncoding.EncodeToString(invalidKey)

	s := NewSigner("000000000000", "wrong_pass", keyStr)

	// Это не должно паниковать — должна вернуться ошибка
	_, err := s.Sign("some data")

	if err == nil {
		t.Error("Expected error for invalid key data, but got nil")
	}
}

func TestSignRejectsInvalidWmid(t *testing.T) {
	invalidKey := make([]byte, 164)
	keyStr := base64.StdEncoding.EncodeToString(invalidKey)

	tests := []struct {
		name string
		wmid string
	}{
		{"too short", "12345"},
		{"too long", "1234567890123"},
		{"letters", "abcdefghijkl"},
		{"with prefix", "abc000000000000"},
		{"empty", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSigner(tc.wmid, "pass", keyStr)
			_, err := s.Sign("data")
			if err == nil {
				t.Errorf("Expected error for WMID %q, got nil", tc.wmid)
			}
		})
	}
}

func TestSignRejectsEmptyPassword(t *testing.T) {
	invalidKey := make([]byte, 164)
	keyStr := base64.StdEncoding.EncodeToString(invalidKey)

	s := NewSigner("000000000000", "", keyStr)
	_, err := s.Sign("data")
	if err == nil {
		t.Error("Expected error for empty password, got nil")
	}
}

func TestSignRejectsInvalidBase64Key(t *testing.T) {
	s := NewSigner("000000000000", "pass", "not-valid-base64!!!")
	_, err := s.Sign("data")
	if err == nil {
		t.Error("Expected error for invalid base64 key, got nil")
	}
}

func TestSignRejectsWrongKeyLength(t *testing.T) {
	shortKey := make([]byte, 100) // не 164 байт
	keyStr := base64.StdEncoding.EncodeToString(shortKey)

	s := NewSigner("000000000000", "pass", keyStr)
	_, err := s.Sign("data")
	if err == nil {
		t.Error("Expected error for wrong key length, got nil")
	}
}
