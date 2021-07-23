package publisher

import (
	"fmt"
	"log"
	"time"

	"github.com/nsterg/go-servicebus-publisher/adapter"
	"github.com/nsterg/go-servicebus-publisher/sas"
)

// NewPublisher initiates a Publisher with an httpClient and some configuration
func NewPublisher(client adapter.HTTPClient, config ServiceBusConfig) publisher {
	return publisher{
		serviceBusAdapter: adapter.NewServiceBusAdapter(client),
		generator:         sas.NewSasGenerator(realClock{}),
		config:            config,
	}
}

// Publish publishes the given message, using publisher's already provided configuration
func (p publisher) Publish(message interface{}) error {
	baseURL := p.config.BaseURL
	namespace := p.config.Namespace
	endpoint := p.config.Endpoint
	signingKey := p.config.SigningKey
	expiry := p.config.SasTokenExpiresMS
	sharedKeyName := p.config.SharedKeyName

	sasToken := p.generator.Generate(getResourceUri(namespace, endpoint), signingKey, expiry, sharedKeyName)

	err := p.serviceBusAdapter.SendMessage(getMessagesUri(baseURL, endpoint), sasToken, message)
	if err != nil {
		log.Printf("Failed to publish message to endpoint %s due to error %s", endpoint, err.Error())
		return err
	}
	log.Printf("Successfully published message to endpoint %s", endpoint)
	return nil
}

// ServiceBusConfig is a struct holding all the configuration necessary to execute the call to azure servicebus messages endpoint
// BaseURL follows this format http{s}://{serviceNamespace}.servicebus.windows.net/
// Namespace is the serviceNamespace
// Endpoint is the queue or topic name
// SharedKeyName is the rule name
// SigningKey is the sharedAccessKey
// SasTokenExpiresMS is the time in milliseconds that the sas token generated will be valid
type ServiceBusConfig struct {
	BaseURL           string
	Namespace         string
	Endpoint          string
	SharedKeyName     string
	SigningKey        string
	SasTokenExpiresMS int
}

type ServiceBusAdapter interface {
	SendMessage(url string, sasToken string, message interface{}) error
}

type SasTokenGenerator interface {
	Generate(resourceUri string, signingKey string, expiresInMins int, policyName string) string
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

type publisher struct {
	serviceBusAdapter ServiceBusAdapter
	generator         SasTokenGenerator
	config            ServiceBusConfig
}

func getMessagesUri(baseURL, endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURL, endpoint)
}

func getResourceUri(namespace, endpoint string) string {
	return fmt.Sprintf("%s.servicebus.windows.net/%s", namespace, endpoint)
}
