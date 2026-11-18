package wmsigner

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestSignKWMV2(t *testing.T) {
	const wmid = "012345678901"
	const password = "test-😀-password"

	// RFC 8032 test seed. The container is generated solely for this test.
	seed, err := hex.DecodeString("9d61b19deffd5a60ba844af492ec2cc44449c5697b326919703bac031cae7f60")
	if err != nil {
		t.Fatal(err)
	}
	header := make([]byte, 52)
	binary.LittleEndian.PutUint16(header[0:2], 0x81)
	binary.LittleEndian.PutUint32(header[20:24], 112)
	binary.LittleEndian.PutUint32(header[24:28], 2)

	// .NET ASCII encoding replaces each UTF-16 surrogate of 😀 with a question mark.
	key := argon2.IDKey([]byte("test-??-password"), header[4:20], 2, 65536, 1, 32)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		t.Fatal(err)
	}
	plain := make([]byte, 68)
	binary.LittleEndian.PutUint16(plain[0:2], 32)
	copy(plain[2:34], seed)
	binary.LittleEndian.PutUint16(plain[34:36], 32)
	file := aead.Seal(header, header[28:52], plain, append(append([]byte(nil), header...), wmid...))
	encoded := base64.StdEncoding.EncodeToString(file)

	signer := NewSigner(wmid, password, encoded)
	got, err := signer.Sign("Cafe\u0301")
	if err != nil {
		t.Fatal(err)
	}
	want := hex.EncodeToString(ed25519.Sign(ed25519.NewKeyFromSeed(seed), []byte("Café")))
	if got != want {
		t.Fatalf("signature = %q, want %q", got, want)
	}

	file[len(file)-1] ^= 1
	tampered := base64.StdEncoding.EncodeToString(file)
	file[len(file)-1] ^= 1
	file[24] = 3
	invalidHeader := base64.StdEncoding.EncodeToString(file)
	for _, tc := range []struct {
		name, wmid, password, key, message, wantError string
	}{
		{"wrong password", wmid, "wrong", encoded, "message", "authentication"},
		{"wrong WMID", "112345678901", password, encoded, "message", "authentication"},
		{"tampered tag", wmid, password, tampered, "message", "authentication"},
		{"invalid header", wmid, password, invalidHeader, "message", "header"},
		{"carriage return", wmid, password, encoded, "test\r\n", "carriage return"},
		{"invalid UTF-8", wmid, password, encoded, "test\xff", "invalid UTF-8"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewSigner(tc.wmid, tc.password, tc.key).Sign(tc.message)
			if got != "" || err == nil || !strings.Contains(err.Error(), tc.wantError) {
				t.Fatalf("Sign = (%q, %v), want error containing %q", got, err, tc.wantError)
			}
		})
	}
}
