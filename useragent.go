package tesla

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

func userAgent() (string, error) {
	const prefixBytes = 6
	var buf [prefixBytes + 8]byte
	if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
		return "", err
	}

	var b strings.Builder
	e := base32.NewEncoder(base32.StdEncoding.WithPadding(base32.NoPadding), &b)

	_, err := e.Write(buf[:prefixBytes])
	if err == nil {
		err = e.Close()
	}

	if err == nil {
		_, err = b.WriteRune('/')
	}
	if err == nil {
		_, err = b.Write(strconv.AppendUint(nil, binary.BigEndian.Uint64(buf[prefixBytes:]), 10))
	}

	return b.String(), err
}
