package tesla

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/oauth2"
)

const (
	mfaVerifyURL = "https://auth.tesla.com/oauth2/v3/authorize/mfa/verify"
)

// Device is the multi-factor device returned by the /authorize/mfa/factors endpoint
type Device struct {
	DispatchRequired bool      `json:"dispatchRequired"`
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	FactorType       string    `json:"factorType"`
	FactorProvider   string    `json:"factorProvider"`
	SecurityLevel    int       `json:"securityLevel"`
	Activated        time.Time `json:"activatedAt"`
	Updated          time.Time `json:"updatedAt"`
}

type auth struct {
	Client       *http.Client
	AuthURL      string
	SelectDevice func(ctx context.Context, devices []Device) (d Device, passcode string, err error)
	SolveCaptcha func(ctx context.Context, captcha io.Reader) (res string, err error)
}

func (a *auth) initClient(ctx context.Context) {
	if client, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		a.Client = client
		return
	}

	a.Client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func (a *auth) Do(ctx context.Context, username, password string) (code string, err error) {
	if a.Client == nil {
		a.initClient(ctx)
	}

	if a.Client.Jar == nil {
		var err error
		a.Client.Jar, err = cookiejar.New(nil)
		if err != nil {
			return "", fmt.Errorf("new cookie jar: %w", err)
		}
	}

	cr := a.Client.CheckRedirect
	a.Client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() { a.Client.CheckRedirect = cr }()

	res, v, err := a.login(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("login: %w", err)
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusFound:
		return codeFromResponse(res)
	default:
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	transactionID := v.Get("transaction_id")

	devices, err := a.listDevices(ctx, transactionID)
	if err != nil {
		return "", fmt.Errorf("list devices: %w", err)
	}

	if len(devices) == 0 {
		return "", errors.New("no devices")
	}

	d, passcode, err := a.SelectDevice(ctx, devices)
	if err != nil {
		return "", fmt.Errorf("select device: %w", err)
	}

	csrf := v.Get("_csrf")
	if err := a.verify(ctx, csrf, transactionID, d, passcode); err != nil {
		return "", fmt.Errorf("verify: %w", err)
	}
	return a.commit(ctx, transactionID)
}

func (a *auth) login(ctx context.Context, username, password string) (*http.Response, url.Values, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.AuthURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new request: %w", err)
	}

	res, err := a.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	v := url.Values{
		"identity":   {username},
		"credential": {password},
	}

	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("new document: %w", err)
	}

	d.Find("input[type=hidden],input[name=captcha]").Each(func(_ int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if !ok {
			return
		}
		value, ok := s.Attr("value")
		if name == "captcha" {
			ok = true
		}
		if !ok {
			return
		}
		v.Set(name, value)
	})

	if _, required := v["captcha"]; required {
		u, _ := url.Parse(a.AuthURL)
		u.Path = "/captcha"

		res, err := a.Client.Get(u.String())
		if err != nil {
			return nil, nil, fmt.Errorf("access captcha: %w", err)
		}
		defer res.Body.Close()

		solution, err := a.SolveCaptcha(ctx, res.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("solve captcha: %w", err)
		}

		v["captcha"] = []string{solution}
	}

	v.Set("_phase", "authenticate")
	v.Set("_process", "1")

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, a.AuthURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for _, cookie := range res.Cookies() {
		cookie := &http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
		req.AddCookie(cookie)
	}

	res, err = a.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do request: %w", err)
	}

	return res, v, err
}

func (a *auth) listDevices(ctx context.Context, transactionID string) ([]Device, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, (&url.URL{
		Scheme:   "https",
		Host:     "auth.tesla.com",
		Path:     "oauth2/v3/authorize/mfa/factors",
		RawQuery: url.Values{"transaction_id": {transactionID}}.Encode(),
	}).String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	res, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	var out struct {
		Data []Device `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	return out.Data, nil
}

func (a *auth) verify(ctx context.Context, csrf string, transactionID string, d Device, passcode string) error {

	if csrf == "" {
		return errors.New("csrf token is missing for verifing MFA")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(map[string]string{
		"transaction_id": transactionID,
		"factor_id":      d.ID,
		"passcode":       passcode,
		"_csrf":          csrf,
	}); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mfaVerifyURL, &buf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	var out MFAVerify
	if err := json.Unmarshal(b, &out); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	if !out.Data.Approved {
		return errors.New("not approved")
	}
	return nil
}

func (a *auth) commit(ctx context.Context, transactionID string) (code string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.AuthURL, strings.NewReader(url.Values{
		"transaction_id": {transactionID},
	}.Encode()))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := a.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusFound {
		return "", fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return codeFromResponse(res)
}

func codeFromResponse(res *http.Response) (code string, err error) {
	u, err := res.Location()
	if err != nil {
		return "", fmt.Errorf("response location: %w", err)
	}
	return u.Query().Get("code"), nil
}
