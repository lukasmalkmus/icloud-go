package icloud

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	container   = "iCloud.com.lukasmalkmus.Example-App"
	keyID       = "6459b5a4f4ce9c2dbaaf09ae996235cfa8a0c96bd1b595fb9fc229d35beb4c9a"
	environment = Development
)

var privateKey *ecdsa.PrivateKey

func init() {
	var err error
	privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
}

// SetStrictDecoding is a special testing-only client option that failes JSON
// response decoding if fields not present in the destination struct are
// encountered.
func SetStrictDecoding() Option {
	return func(c *Client) error {
		c.strictDecoding = true
		return nil
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient(container, keyID, privateKey, environment)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Are endpoints/resources present?
	assert.NotNil(t, client.Records)

	// Is default configuration present?
	expURL := "https://api.apple-cloudkit.com/database/1/iCloud.com.lukasmalkmus.Example-App/development"
	assert.Equal(t, expURL, client.baseURL.String())
	assert.NotEmpty(t, client.userAgent)
	assert.False(t, client.strictDecoding)
	assert.NotNil(t, client.httpClient)
}

// func TestClient_newRequest_BadURL(t *testing.T) {
// 	client, _ := NewClient(container, keyID, privateKey, environment)

// 	_, err := client.newRequest(context.Background(), http.MethodGet, ":", nil)
// 	assert.Error(t, err)

// 	if assert.IsType(t, new(url.Error), err) {
// 		urlErr := err.(*url.Error)
// 		assert.Equal(t, urlErr.Op, "parse")
// 	}
// }

// If a nil body is passed to NewRequest, make sure that nil is also passed to
// http.NewRequest. In most cases, passing an io.Reader that returns no content
// is fine, since there is no difference between an HTTP request body that is an
// empty string versus one that is not set at all. However in certain cases,
// intermediate systems may treat these differently resulting in subtle errors.
func TestClient_newRequest_EmptyBody(t *testing.T) {
	client, _ := NewClient(container, keyID, privateKey, environment)

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	assert.Empty(t, req.Body)
}

func TestClient_do(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, `{"A":"a"}`)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	type foo struct {
		A string
	}

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	var body foo
	err = client.do(req, &body)
	require.NoError(t, err)

	assert.Equal(t, foo{"a"}, body)
}

func TestClient_do_ioWriter(t *testing.T) {
	content := `{"A":"a"}`

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = fmt.Fprint(w, content)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = client.do(req, &buf)
	require.NoError(t, err)

	assert.Equal(t, content, buf.String())
}

func TestClient_do_HTTPError(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusBadRequest

		httpErr := Error{
			Reason: http.StatusText(code),
		}

		err := json.NewEncoder(w).Encode(httpErr)
		assert.NoError(t, err)

		w.WriteHeader(code)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.NoError(t, err)
}

func TestClient_do_RedirectLoop(t *testing.T) {
	hf := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	client, teardown := setup(t, "/", hf)
	defer teardown()

	req, err := client.newRequest(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	err = client.do(req, nil)
	require.Error(t, err)

	assert.IsType(t, err, new(url.Error))
}

// setup sets up a test HTTP server along with a client that is configured to
// talk to that test server. Tests should pass a handler function which provides
// the response for the API method being tested.
func setup(t *testing.T, path string, handler http.HandlerFunc) (*Client, func()) { //nolint:unparam // ...
	t.Helper()

	r := http.NewServeMux()
	r.HandleFunc(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("accept"), "application/json", "bad accept header present on the request")
		assert.Equal(t, r.Header.Get("user-agent"), "icloud-go", "bad user-agent header present on the request")

		assert.Equal(t, r.Header.Get("x-apple-cloudkit-request-keyid"), keyID, "bad x-apple-cloudkit-request-keyid header present on the request")
		assert.NotEmpty(t, r.Header.Get("x-apple-cloudkit-request-iso8601date"), "missing x-apple-cloudkit-request-iso8601date header")
		assert.NotEmpty(t, r.Header.Get("x-apple-cloudkit-request-signaturev1"), "missing x-apple-cloudkit-request-signaturev1 header")

		if r.ContentLength > 0 {
			assert.NotEmpty(t, r.Header.Get("Content-type"), "no Content-type header present on the request")
		}

		handler.ServeHTTP(w, r)
	}))
	srv := httptest.NewServer(r)

	urlStr := fmt.Sprintf("%s/database/%d/%s/%s", srv.URL, apiVersion, container, environment)
	baseURL, err := url.ParseRequestURI(urlStr)
	require.NoError(t, err)

	client, err := NewClient(container, keyID, privateKey, environment, SetHTTPClient(srv.Client()), SetStrictDecoding())
	require.NoError(t, err)

	client.baseURL = baseURL

	return client, func() { srv.Close() }
}

// func mustTimeParse(t *testing.T, layout, value string) time.Time {
// 	ts, err := time.Parse(layout, value)
// 	require.NoError(t, err)
// 	return ts
// }
