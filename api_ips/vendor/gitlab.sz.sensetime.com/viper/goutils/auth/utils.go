package auth

import (
	"crypto/md5" // #nosec
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net"
	"os"
)

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func GenerateUnqiueClientID(entropy bool) string {
	macs, _ := getMacAddr()
	md5 := md5.New() // #nosec
	for _, v := range macs {
		_, _ = md5.Write([]byte(v))
	}
	sum := md5.Sum(nil)
	host, _ := os.Hostname()
	if host == "" {
		host = "unknownhost"
	}
	var e [8]byte
	if entropy {
		_, _ = rand.Read(e[:])
	}
	return host + "-" + hex.EncodeToString(sum) + "-" + hex.EncodeToString(e[:])
}

func checkEncryptionKey(key interface{}) ([]byte, error) {
	encryptionKey, ok := key.(string)
	if ok && len(encryptionKey) == 0 {
		return nil, nil
	} else if ok && len(encryptionKey) == 16 {
		return []byte(encryptionKey), nil
	}
	return nil, errors.New("bad key format")
}

// nolint
func checksumLic(b []byte) string {
	h := sha256.New()
	h.Write(b) // nolint
	return hex.EncodeToString(h.Sum(nil))
}
