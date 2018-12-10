package lib

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/Yangshuting/golang_model/storage"
)

func GenRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenRandomString(s int) (string, error) {
	b, err := GenRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func Gen3rdSession(s int) (string, error) {
	session, err := GenRandomString(s)
	return session, err
}

func SetSession(cc *Cusctx, uid string) (string, error) {
	session, err := Gen3rdSession(32)
	storage.SetRedis(session, uid)
	return session, err
}
