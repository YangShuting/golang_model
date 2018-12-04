package lib

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(in ...string) string {
	buf := new(bytes.Buffer)

	for idx := range in {
		buf.WriteString(in[idx])
	}
	tmp := sha1.Sum(buf.Bytes())
	return hex.EncodeToString(tmp[:])
}
