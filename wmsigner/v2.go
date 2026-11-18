package wmsigner

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/text/unicode/norm"
)

const (
	kwmV2Size         = 136
	kwmV2HeaderSize   = 52
	kwmV2BodySize     = 112
	kwmV2PlainSize    = 68
	kwmV2SaltStart    = 4
	kwmV2SaltEnd      = 20
	kwmV2NonceStart   = 28
	kwmV2ArgonMemory  = 64 * 1024 // KiB
	kwmV2ArgonPasses  = 2
	kwmV2ArgonThreads = 1
)

// The v2 container format is implemented by V2KeyFile in WebMoney Key Extractor.
func (w *Wmsigner) signV2(data string, file []byte) (string, error) {
	if binary.LittleEndian.Uint16(file[:2]) != 0x81 ||
		binary.LittleEndian.Uint32(file[20:24]) != kwmV2BodySize ||
		binary.LittleEndian.Uint32(file[24:28]) != 2 {
		return "", errors.New("invalid WebMoney v2 key header")
	}
	if !utf8.ValidString(data) {
		return "", errors.New("WebMoney v2 message contains invalid UTF-8")
	}
	if strings.ContainsRune(data, '\r') {
		return "", errors.New("WebMoney v2 message contains a carriage return")
	}

	// Match .NET ASCII encoding, which replaces each non-ASCII UTF-16 unit.
	password := make([]byte, 0, len(w.pass))
	for _, r := range w.pass {
		if r <= 0x7f {
			password = append(password, byte(r))
			continue
		}
		password = append(password, '?')
		if r > 0xffff {
			password = append(password, '?')
		}
	}
	fileKey := argon2.IDKey(password, file[kwmV2SaltStart:kwmV2SaltEnd],
		kwmV2ArgonPasses, kwmV2ArgonMemory, kwmV2ArgonThreads, chacha20poly1305.KeySize)
	wipeBytes(password)
	aead, err := chacha20poly1305.NewX(fileKey)
	wipeBytes(fileKey)
	if err != nil {
		return "", err
	}

	associatedData := make([]byte, 0, kwmV2HeaderSize+len(w.wmid))
	associatedData = append(associatedData, file[:kwmV2HeaderSize]...)
	associatedData = append(associatedData, w.wmid...)
	plain, err := aead.Open(nil, file[kwmV2NonceStart:kwmV2HeaderSize], file[kwmV2HeaderSize:], associatedData)
	if err != nil {
		return "", fmt.Errorf("WebMoney v2 key authentication failed: %w", err)
	}
	defer wipeBytes(plain)

	const seedEnd = 2 + ed25519.SeedSize
	if len(plain) != kwmV2PlainSize ||
		binary.LittleEndian.Uint16(plain[:2]) != ed25519.SeedSize ||
		binary.LittleEndian.Uint16(plain[seedEnd:seedEnd+2]) != 32 {
		return "", errors.New("invalid WebMoney v2 key body")
	}

	privateKey := ed25519.NewKeyFromSeed(plain[2:seedEnd])
	defer wipeBytes(privateKey)
	signature := ed25519.Sign(privateKey, []byte(norm.NFC.String(data)))
	return hex.EncodeToString(signature), nil
}

func wipeBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
