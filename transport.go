package tesla

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"
)

type Transport struct {
	http.RoundTripper
	userAgent     string
	userAgentTime time.Time
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if now := time.Now(); t.userAgent == "" || now.Sub(t.userAgentTime) > 6*time.Hour {
		var err error
		if t.userAgent, err = userAgent(); err != nil {
			return nil, err
		}
		t.userAgentTime = now
	}
	for _, h := range []struct{ k, v string }{
		{"Accept", "*/*"},
		{"Accept-Encoding", "gzip, deflate, br"},
		{"User-Agent", t.userAgent},
	} {
		if _, ok := req.Header[h.k]; ok {
			continue
		}
		req.Header.Set(h.k, h.v)
	}

	res, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if strings.EqualFold(res.Header.Get("Content-Encoding"), "gzip") {
		res.Body = &gzipReader{body: res.Body}
	}

	return res, err
}

type gzipReader struct {
	body io.ReadCloser
	zr   *gzip.Reader
	zerr error
}

func (gz *gzipReader) Read(p []byte) (n int, err error) {
	if gz.zr == nil {
		if gz.zerr == nil {
			gz.zr, gz.zerr = gzip.NewReader(gz.body)
		}
		if gz.zerr != nil {
			return 0, gz.zerr
		}
	}
	return gz.zr.Read(p)
}

func (gz *gzipReader) Close() error {
	return gz.body.Close()
}
