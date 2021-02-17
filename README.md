# go-servicebus-publisher

This is a library to publish a message to an azure service endpoint (topic or queue).

### Usage

Import go-servicebus-publisher

```import ("github.com/nsterg/go-servicebus-publisher/publisher")```

or

```go get github.com/nsterg/go-servicebus-publisher/publisher```

From your service initialize a publisher with all the dependant services

```go
	serviceBusClient := &http.Client{Timeout: time.Duration(1000) * time.Millisecond}

	// these values will be derived from your azure service bus configuration
	config := publisher.ServiceBusConfig{
		Namespace:           "my-name-space",
		Endpoint:            "my-queue",
		SharedKeyName:       "my-shared-keyname",
		SigningKey:          "my-signing-key",
		SigningKeyExpiresMS: 1234,
		EndpointURL:         "https://my-name-space.servicebus.windows.net/my-queue"
	}
	p := publisher.NewPublisher(serviceBusClient, config)
```

 and publish a message

```go
	p.Publish(map[string]string{"I am here": "to poke you"})
```
