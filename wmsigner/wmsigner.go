package wmsigner

/* example for use

s := signer.Newsigner(wmid,pass,key)

signature := s.Sign(String)

*/

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/crypto/md4"
	"math/big"
	"regexp"
)

// wmidRegex is pre-compiled once at package init for performance.
var wmidRegex = regexp.MustCompile(`^[0-9]{12}$`)

// NewSigner return wmsigner object with wmid, pass and key
func NewSigner(wmid string, pass string, key string) *Wmsigner {
	return &Wmsigner{wmid: wmid, pass: pass, key: key}
}

/*
function hashing  butes slice to md4
*/
func md4Hash(buff []byte) []byte {
	hash := md4.New()
	hash.Write(buff)
	return hash.Sum(nil)
}

/*
xor two byte slice
*/
func xorWm(buff []byte, hash []byte) []byte {
	b := bytes.NewBuffer([]byte{})
	for i := 0; i < len(buff); i++ {
		b.WriteByte(buff[i] ^ hash[i%len(hash)])
	}
	return b.Bytes()
}

/*
reverse byte in bytes slice
*/

func reversebyte(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

/*create bignum object from bytes.Buffer*/
func readBn(b *bytes.Buffer) (*big.Int, error) {
	var count uint16
	if err := binary.Read(b, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("malformed key data: %w", err)
	}
	data := b.Next(int(count))
	if len(data) != int(count) {
		return nil, errors.New("malformed key data: truncated big number")
	}
	i := new(big.Int)
	i.SetBytes(reversebyte(data))
	return i, nil
}

type Wmsigner struct {
	wmid string
	pass string
	key  string
}

// Sign signs data with a classic RSA or v2 Ed25519 WebMoney key.
func (w *Wmsigner) Sign(data string) (string, error) {

	//check wmid
	if !wmidRegex.MatchString(w.wmid) {
		return "", errors.New("Wmid is incorrect: must be exactly 12 digits")
	}

	//check key
	decode_key, err := base64.StdEncoding.DecodeString(w.key)
	if err != nil {

		return "", err
	}
	if len(decode_key) != 164 && len(decode_key) != kwmV2Size {
		return "", errors.New("Illegal size for keydata")
	}

	//checkpassword
	if len(w.pass) == 0 {
		return "", errors.New("Password can not to be empty!")
	}
	if len(decode_key) == kwmV2Size {
		return w.signV2(data, decode_key)
	}
	powers, modules, err := w.initkey()
	if err != nil {
		return "", err
	}
	hash := md4Hash([]byte(data))
	buf := bytes.NewBuffer(hash)
	var ini uint32
	//add rendom 40 bytes
	for i := 1; i <= 10; i++ {
		if err := binary.Read(rand.Reader, binary.LittleEndian, &ini); err != nil {
			return "", fmt.Errorf("crypto/rand failed: %w", err)
		}
		binary.Write(buf, binary.LittleEndian, ini)
	}
	hash = buf.Bytes()
	lenn := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenn, uint16(len(hash)))
	b := bytes.NewBuffer(lenn)
	b.Write(hash)

	for b.Len() < len(modules.Bytes()) {
		b.WriteByte(byte(0))
	}

	data_res := b.Bytes()
	data_res = reversebyte(data_res)
	i := new(big.Int)
	i.SetBytes(data_res)
	signature := i.Exp(i, powers, modules)

	modlen := len(modules.Bytes())
	sig := bytes.NewBuffer([]byte{})
	for i := len(signature.Bytes()); i < modlen; i++ {
		sig.WriteByte(0)
	}
	sig.Write(signature.Bytes())
	sigbyte := reversebyte(sig.Bytes())
	if len(sigbyte)%2 != 0 {
		// Odd modulus length: pad the most significant byte so the
		// uint16 packing below does not silently drop it. Real WM keys
		// have even-length moduli; their output is unchanged.
		sigbyte = append(sigbyte, 0)
	}
	sig = bytes.NewBuffer(sigbyte)
	var s = make([]uint16, len(sigbyte)/2)
	err = binary.Read(sig, binary.LittleEndian, s)
	if err != nil {
		return "", err
	}
	result := ""
	for _, val := range s {
		result += fmt.Sprintf("%04x", val)
	}

	return result, nil

}

func (w *Wmsigner) initkey() (*big.Int, *big.Int, error) {
	// Совместимость с WM Keeper Classic: ключи, созданные старыми версиями
	// Keeper, зашифрованы первой половиной пароля. Официальный WMSigner
	// (Signer::SecureKeyByIDPWHalf) и эталонный PHP-signer пробуют полный
	// пароль и его первую половину. Сам пароль никогда не мутирует.
	powers, modules, err := w.initkeyWithPass(w.pass)
	if err == nil {
		return powers, modules, nil
	}
	if half := len(w.pass) / 2; half < len(w.pass) {
		if p2, m2, err2 := w.initkeyWithPass(w.pass[:half]); err2 == nil {
			return p2, m2, nil
		}
	}
	// Обе попытки не прошли — возвращаем ошибку полного пароля.
	return nil, nil, err
}

func (w *Wmsigner) initkeyWithPass(pass string) (*big.Int, *big.Int, error) {
	data, err := base64.StdEncoding.DecodeString(w.key)
	if err != nil {
		return nil, nil, err
	}

	r := bytes.NewBuffer(data)
	header := r.Next(30)
	body := r.Bytes()
	body = w.encryptKey(body, pass)
	r = bytes.NewBuffer([]byte{})

	r.Write(header)
	r.Write(body)

	reserved1 := r.Next(2)
	r.Next(2)
	crc := r.Next(16)
	length := r.Next(4)

	body = r.Bytes()
	r = bytes.NewBuffer([]byte{})
	reserved1 = append(reserved1, byte(0), byte(0))
	r.Write(reserved1)
	r.Write(bytes.Repeat([]byte{0}, 16))
	r.Write(length)
	r.Write(body)

	calculated_crc := md4Hash(r.Bytes())
	if !bytes.Equal(crc, calculated_crc) {
		// Несовпадение CRC для данной попытки пароля. Решение о повторе
		// с половиной пароля принимает initkey; состояние не мутирует.
		return nil, nil, errors.New("key CRC mismatch: wrong password or corrupted key data")
	}
	r = bytes.NewBuffer(body)
	r.Next(4)
	powers, err := readBn(r)
	if err != nil {
		return nil, nil, err
	}
	modules, err := readBn(r)
	if err != nil {
		return nil, nil, err
	}
	if powers.Sign() <= 0 || modules.Sign() <= 0 {
		return nil, nil, errors.New("invalid key data: zero exponent or modulus")
	}
	return powers, modules, nil
}

func (w *Wmsigner) encryptKey(buff []byte, pass string) []byte {
	hashResult := md4Hash([]byte(w.wmid + pass))
	return xorWm(buff, hashResult)
}
