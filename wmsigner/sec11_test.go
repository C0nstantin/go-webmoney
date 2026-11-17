package wmsigner

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"math/big"
	"regexp"
	"testing"
)

// encodeBn кодирует big.Int в формат WM-ключа: uint16 LE длина + байты
// в обратном порядке (readBn делает reversebyte перед SetBytes).
func encodeBn(i *big.Int) []byte {
	be := i.Bytes()
	le := make([]byte, len(be))
	for j, b := range be {
		le[len(be)-1-j] = b
	}
	out := make([]byte, 2+len(le))
	binary.LittleEndian.PutUint16(out, uint16(len(le)))
	copy(out[2:], le)
	return out
}

// sealTestKey проставляет CRC и шифрует регион [30:164] ключа WM
// (XOR md4(wmid+pass)); CRC считается по незашифрованной раскладке.
func sealTestKey(plain []byte, wmid, pass string) string {
	crcBuf := bytes.NewBuffer(nil)
	crcBuf.Write(plain[0:2])
	crcBuf.Write([]byte{0, 0})
	crcBuf.Write(make([]byte, 16))
	crcBuf.Write(plain[20:24])
	crcBuf.Write(plain[24:])
	copy(plain[4:20], md4Hash(crcBuf.Bytes()))

	hash := md4Hash([]byte(wmid + pass))
	for i := 30; i < len(plain); i++ {
		plain[i] ^= hash[(i-30)%len(hash)]
	}
	return base64.StdEncoding.EncodeToString(plain)
}

// buildTestKey собирает валидный 164-байтовый ключ с заданными
// экспонентой и модулем (body с offset 24: 4 reserved + bn powers + bn modules).
func buildTestKey(wmid, pass string, powers, modules *big.Int) string {
	plain := make([]byte, 164)
	buf := bytes.NewBuffer(nil)
	buf.Write(make([]byte, 4))
	buf.Write(encodeBn(powers))
	buf.Write(encodeBn(modules))
	copy(plain[24:], buf.Bytes())
	return sealTestKey(plain, wmid, pass)
}

// Регрессия SEC-11: при несовпадении CRC signer молча усекал пароль вдвое
// и рекурсивно повторял, портя состояние объекта. Теперь — честная ошибка
// без мутации, объект остаётся пригодным для повторных вызовов.
func TestSignWrongPasswordReturnsErrorWithoutMutation(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	key := buildTestKey(wmid, pass, big.NewInt(3), big.NewInt(0x01000001))

	s := NewSigner(wmid, "wrong-password", key)
	for i := 0; i < 3; i++ {
		_, err := s.Sign("some data")
		if err == nil {
			t.Fatalf("call %d: expected CRC error for wrong password, got nil", i)
		}
		if s.pass != "wrong-password" {
			t.Fatalf("call %d: signer password was mutated to %q", i, s.pass)
		}
	}
}

func TestSignValidKeySucceeds(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	key := buildTestKey(wmid, pass, big.NewInt(3), big.NewInt(0x01000001))

	s := NewSigner(wmid, pass, key)
	sig, err := s.Sign("some data")
	if err != nil {
		t.Fatal(err)
	}
	// модуль 4 байта → подпись 4 байта → 8 hex-символов
	if !regexp.MustCompile(`^[0-9a-f]{8}$`).MatchString(sig) {
		t.Errorf("signature = %q, want 8 lowercase hex chars", sig)
	}
}

func TestSignZeroModulusRejected(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	key := buildTestKey(wmid, pass, big.NewInt(3), big.NewInt(0))

	s := NewSigner(wmid, pass, key)
	if _, err := s.Sign("some data"); err == nil {
		t.Error("expected error for zero modulus, got nil")
	}
}

func TestSignTruncatedBnRejected(t *testing.T) {
	const wmid, pass = "123456789012", "correct-password"
	plain := make([]byte, 164)
	// count экспоненты 0xFFFF при пустом буфере — усечённые данные
	binary.LittleEndian.PutUint16(plain[28:30], 0xFFFF)
	key := sealTestKey(plain, wmid, pass)

	s := NewSigner(wmid, pass, key)
	if _, err := s.Sign("some data"); err == nil {
		t.Error("expected error for truncated big number, got nil")
	}
}