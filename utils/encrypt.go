package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

type Aes struct {
	//16、24或者32位的[]byte，分别对应AES-128, AES-192或AES-256算法
	Cip cipher.Block `inject:"aes.cipher"`
}

func (p *Aes) Encrypt(pn []byte) ([]byte, error) {

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(p.Cip, iv)
	ct := make([]byte, len(pn))
	cfb.XORKeyStream(ct, pn)

	return append(ct, iv...), nil

}

func (p *Aes) Decrypt(sr []byte) ([]byte, error) {
	bln := len(sr)
	cln := bln - aes.BlockSize
	ct := sr[0:cln]
	iv := sr[cln:bln]

	cfb := cipher.NewCFBDecrypter(p.Cip, iv)
	pt := make([]byte, cln)
	cfb.XORKeyStream(pt, ct)
	return pt, nil
}

//==============================================================================

type Hmac struct {
	Key []byte           `inject:"hmac.key"` //32 bits
	Fn  func() hash.Hash `inject:"hmac.fn"`
}

func (p *Hmac) Sum(src []byte) []byte {
	mac := hmac.New(p.Fn, p.Key)
	mac.Write(src)
	return mac.Sum(nil)
}

func (p *Hmac) Equal(src, dst []byte) bool {
	return hmac.Equal(src, dst)
}

//==============================================================================

func Md5(p []byte) string {
	buf := md5.Sum([]byte(p))
	return hex.EncodeToString(buf[:])
}

func Sha512(p []byte) string {
	buf := sha512.Sum512(p)
	return hex.EncodeToString(buf[:])
}

func Ssha512(p []byte, l int) (string, error) {
	salt := make([]byte, l)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return ssha512(p, salt), nil
}

func ssha512(d, s []byte) string {
	buf := sha512.Sum512(append(d, s...))
	return base64.StdEncoding.EncodeToString(append(buf[:], s...))
}

func Csha512(d string, p []byte) (bool, error) {
	buf, err := base64.StdEncoding.DecodeString(d)
	if err == nil {
		salt := buf[sha512.Size:]
		return ssha512(p, salt) == d, nil
	}
	return false, err
}
