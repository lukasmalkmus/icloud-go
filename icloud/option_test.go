package icloud

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_SetHTTPClient(t *testing.T) {
	exp := &http.Client{
		Timeout: 0,
	}
	opt := SetHTTPClient(exp)

	evaluateOption(t, opt, func(client *Client) {
		assert.Equal(t, exp, client.httpClient)
	})
}

func TestOption_SetUserAgent(t *testing.T) {
	exp := "icloud-go/1.0.0"
	opt := SetUserAgent(exp)

	evaluateOption(t, opt, func(client *Client) {
		assert.Equal(t, exp, client.userAgent)
	})
}

func evaluateOption(t *testing.T, opt Option, f func(client *Client)) {
	client, _ := NewClient(container, keyID, nil, environment)

	err := client.Options(opt)
	assert.NoError(t, err)

	f(client)
}
