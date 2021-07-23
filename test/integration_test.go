package test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	publisher "github.com/nsterg/go-servicebus-publisher"
	"github.com/stretchr/testify/assert"
)

func TestShouldPublishMessage(t *testing.T) {
	publisherMockServer := MockPublisher(t, "test-queue", "7000", http.StatusCreated)
	defer publisherMockServer.Close()

	c := publisher.ServiceBusConfig{
		BaseURL:       "http://localhost:7000",
		Endpoint:      "test-queue",
		SharedKeyName: "my-shared-key",
		SigningKey:    "my-signing-key",
	}
	p := publisher.NewPublisher(&http.Client{}, c)

	err := p.Publish("some message")
	assert.NoError(t, err)
}

func TestReturnErrorWhenPublishMessageFails(t *testing.T) {
	publisherMockServer := MockPublisher(t, "test-queue", "7000", http.StatusInternalServerError)
	defer publisherMockServer.Close()

	c := publisher.ServiceBusConfig{
		BaseURL:       "http://localhost:7000",
		Endpoint:      "test-queue",
		SharedKeyName: "my-shared-key",
		SigningKey:    "my-signing-key",
	}
	p := publisher.NewPublisher(&http.Client{}, c)

	err := p.Publish("some message")
	assert.Error(t, err)
}

func MockPublisher(t *testing.T, queue string, port string, statusCode int) *httptest.Server {
	publisherServeMux := http.NewServeMux()
	publisherServeMux.HandleFunc(fmt.Sprintf("/%s/messages", queue), func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodPost, req.Method)
		res.WriteHeader(statusCode)
	})

	mockServer := httptest.NewUnstartedServer(publisherServeMux)
	mockServer.Listener.Close()
	mockServer.Listener = createListener(t, port)
	mockServer.Start()

	return mockServer
}

func createListener(t *testing.T, port string) net.Listener {
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		t.Fatal(err)
	}
	return l
}
