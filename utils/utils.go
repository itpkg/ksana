package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pborman/uuid"
)

func FromToml(f string, v interface{}) error {
	_, err := toml.DecodeFile(f, v)
	return err
}

func ToToml(f string, v interface{}) error {
	fi, err := os.Create(f)
	defer fi.Close()

	if err == nil {
		end := toml.NewEncoder(fi)
		err = end.Encode(v)
	}
	return err
}

func PkgRoot(o interface{}) string {
	return fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), reflect.TypeOf(o).Elem().PkgPath())
}

func Mkdirs(d string, m os.FileMode) error {
	fi, err := os.Stat(d)
	if err == nil {
		if fi.IsDir() {
			return nil
		}
		return errors.New(fmt.Sprintf("%s is a file", d))
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(d, m)
	}
	return err

}

func Uuid() string {
	return uuid.New()
}

func FromHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func ToHex(bs []byte) string {
	return hex.EncodeToString(bs)
}

func ToBits(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func FromBits(data []byte, obj interface{}) error {
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(data)
	err := dec.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

func ToJson(o interface{}) ([]byte, error) {
	return json.Marshal(o)

}

func FromJson(j []byte, o interface{}) error {
	return json.Unmarshal(j, o)
}

func ToBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func FromBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func Shell(cmd string, args ...string) error {
	bin, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}
	return syscall.Exec(bin, append([]string{cmd}, args...), os.Environ())
}

func Equal(src []byte, dst []byte) bool {
	if src == nil && dst == nil {
		return true
	}
	if len(src) == len(dst) {
		for i, b := range src {
			if b != dst[i] {
				return false
			}
		}
	}
	return false
}

func AppendSalt(src, salt []byte) []byte {
	return append(src, salt...)
}

func ParseSalt(src []byte, length int) ([]byte, []byte) {
	size := len(src)
	return src[0 : size-length], src[size-length : size]
}

func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func init() {
	gob.Register(time.Time{})
}
