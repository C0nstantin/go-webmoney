package wmsigner

import (
	"math/big"
	"testing"
)

// Совместимость с WM Keeper Classic: ключи старых версий Keeper зашифрованы
// первой половиной пароля (официальный C++ WMSigner SecureKeyByIDPWHalf;
// та же фича в PHP-реализации baibaratsky/php-wmsigner из списка WM wiki).
// Signer обязан принимать такие ключи с полным паролем, не мутируя сам пароль.

func TestSignHalfPasswordKeyCompatibility(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password" // 16 символов → половина 8
	key := buildTestKey(wmid, pass[:len(pass)/2], big.NewInt(3), big.NewInt(0x01000001))

	s := NewSigner(wmid, pass, key)
	sig, err := s.Sign("some data")
	if err != nil {
		t.Fatalf("key sealed with half password must open with full password: %v", err)
	}
	if len(sig) != 8 {
		t.Errorf("signature length = %d, want 8", len(sig))
	}
	if s.pass != pass {
		t.Errorf("signer password was mutated to %q", s.pass)
	}
	// Повторные вызовы работают так же (старый код портил состояние).
	for i := 0; i < 3; i++ {
		if _, err := s.Sign("some data"); err != nil {
			t.Fatalf("repeat call %d: %v", i, err)
		}
	}
}

func TestSignHalfPasswordOddLength(t *testing.T) {
	const wmid, pass = "123456789012", "odd-pass9" // 9 символов → половина 4
	key := buildTestKey(wmid, pass[:len(pass)/2], big.NewInt(3), big.NewInt(0x01000001))

	s := NewSigner(wmid, pass, key)
	if _, err := s.Sign("some data"); err != nil {
		t.Fatalf("odd-length half password: %v", err)
	}
}

func TestSignFullPasswordPreferredOverHalf(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	// Ключ зашифрован ПОЛНЫМ паролем — должен открываться с первой попытки.
	key := buildTestKey(wmid, pass, big.NewInt(3), big.NewInt(0x01000001))

	s := NewSigner(wmid, pass, key)
	if _, err := s.Sign("some data"); err != nil {
		t.Fatal(err)
	}
}

func TestSignWrongPasswordStillFailsWithHalfFallback(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	key := buildTestKey(wmid, pass, big.NewInt(3), big.NewInt(0x01000001))

	// "wrong-password" → половина "wrong-p": ни одна не подходит.
	s := NewSigner(wmid, "wrong-password", key)
	for i := 0; i < 3; i++ {
		if _, err := s.Sign("some data"); err == nil {
			t.Fatalf("call %d: expected error for wrong password", i)
		}
		if s.pass != "wrong-password" {
			t.Fatalf("call %d: password mutated to %q", i, s.pass)
		}
	}
	// После неудач signer остаётся пригодным: верный пароль открывает ключ.
	s.pass = pass
	if _, err := s.Sign("some data"); err != nil {
		t.Fatalf("signer unusable after failed attempts: %v", err)
	}
}