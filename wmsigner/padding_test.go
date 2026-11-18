package wmsigner

import (
	"encoding/binary"
	"math/big"
	"regexp"
	"strconv"
	"testing"
)

// decodeSig разворачивает формат вывода Sign (конкатенация %04x
// little-endian uint16) обратно в число.
func decodeSig(t *testing.T, sig string) *big.Int {
	t.Helper()
	if len(sig)%4 != 0 {
		t.Fatalf("signature length %d is not a multiple of 4", len(sig))
	}
	le := make([]byte, 0, len(sig)/2)
	for i := 0; i < len(sig); i += 4 {
		v, err := strconv.ParseUint(sig[i:i+4], 16, 16)
		if err != nil {
			t.Fatalf("bad hex chunk %q: %v", sig[i:i+4], err)
		}
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v))
		le = append(le, b...)
	}
	be := make([]byte, len(le))
	for j := range le {
		be[len(le)-1-j] = le[j]
	}
	return new(big.Int).SetBytes(be)
}

// Регрессия: padding-цикл писал modlen+1 байт, а упаковка в uint16 молча
// отбрасывала лишнее. Теперь длина ровно modlen (с явным дополнением для
// нечётных модулей); формат вывода для реальных (чётных) ключей не меняется.
func TestSignPaddingEvenModulus(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	modules := big.NewInt(0x01000001) // 4 байта
	key := buildTestKey(wmid, pass, big.NewInt(3), modules)

	s := NewSigner(wmid, pass, key)
	sig, err := s.Sign("some data")
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^[0-9a-f]{8}$`).MatchString(sig) {
		t.Fatalf("signature = %q, want 8 hex chars for 4-byte modulus", sig)
	}
	if v := decodeSig(t, sig); v.Cmp(modules) >= 0 {
		t.Errorf("signature value %s >= modulus %s", v, modules)
	}
}

func TestSignPaddingOddModulus(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	modules := big.NewInt(65537) // 3 байта — нечётная длина
	key := buildTestKey(wmid, pass, big.NewInt(3), modules)

	s := NewSigner(wmid, pass, key)
	sig, err := s.Sign("some data")
	if err != nil {
		t.Fatal(err)
	}
	// 3 байта дополняются до 4 для упаковки uint16 → 8 hex (как и раньше),
	// но теперь дополнение явное, а не результат молчаливой обрезки.
	if !regexp.MustCompile(`^[0-9a-f]{8}$`).MatchString(sig) {
		t.Fatalf("signature = %q, want 8 hex chars for padded 3-byte modulus", sig)
	}
	if v := decodeSig(t, sig); v.Cmp(modules) >= 0 {
		t.Errorf("signature value %s >= modulus %s", v, modules)
	}
}