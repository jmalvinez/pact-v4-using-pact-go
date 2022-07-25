package pact_v4

import (
	"bytes"
	"fmt"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestConsumer(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "TestConsumer",
		Provider: "TestProvider",
	})
	assert.NoError(t, err)

	err = mockProvider.
		AddInteraction().
		UponReceiving("A request to do a foo").
		WithRequest("POST", "/foobar", func(b *consumer.V4RequestBuilder) {
			b.
				Header("Content-Type", matchers.S("application/json"))
		}).
		WillRespondWith(http.StatusOK).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := newClient(config.Host, config.Port)

			err = client.SendProduct(`{"id": "1"}`)

			assert.NoError(t, err)

			return err
		})
}

type productAPIClient struct {
	port int
	host string
}

func newClient(host string, port int) *productAPIClient {
	return &productAPIClient{
		host: host,
		port: port,
	}
}

func (u *productAPIClient) SendProduct(productJson string) error {
	_, err := http.Post(
		fmt.Sprintf("http://%s:%d/foobar", u.host, u.port),
		"application/json",
		bytes.NewBuffer([]byte(productJson)),
	)

	return err
}
