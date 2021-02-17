package publisher

import (
	"fmt"
	"log"
	"time"

	"github.com/nsterg/go-servicebus-publisher/adapter"
	"github.com/nsterg/go-servicebus-publisher/sas"
)

type Publisher struct {
	serviceBusAdapter ServiceBusAdapter
	generator         SasTokenGenerator
	config            ServiceBusConfig
}

func NewPublisher(client adapter.HTTPClient, config ServiceBusConfig) Publisher {
	return Publisher{
		serviceBusAdapter: adapter.NewServiceBusAdapter(client),
		generator:         sas.NewSasGenerator(realClock{}),
		config:            config,
	}
}

func (p Publisher) Publish(message interface{}) error {
	var (
		baseURL       = p.config.EndpointBaseURL
		namespace     = p.config.Namespace
		endpoint      = p.config.Endpoint
		signingKey    = p.config.SigningKey
		expiry        = p.config.SigningKeyExpiresMS
		sharedKeyName = p.config.SharedKeyName
	)

	sasToken := p.generator.Generate(fmt.Sprintf("%s.servicebus.windows.net/%s", namespace, endpoint), signingKey, expiry, sharedKeyName)

	err := p.serviceBusAdapter.SendMessage(baseURL, sasToken, message)
	if err != nil {
		log.Printf("Failed to publish message to endpoint %s due to error %s", endpoint, err.Error())
		return err
	}
	log.Printf("Failed to publish message to endpoint %s", endpoint)
	return nil
}

type ServiceBusAdapter interface {
	SendMessage(url string, sasToken string, message interface{}) error
}

type SasTokenGenerator interface {
	Generate(resourceUri string, signingKey string, expiresInMins int, policyName string) string
}

type ServiceBusConfig struct {
	EndpointBaseURL     string
	Namespace           string
	Endpoint            string
	SharedKeyName       string
	SigningKey          string
	SigningKeyExpiresMS int
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }
