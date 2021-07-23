# go-servicebus-publisher

This is a library to publish a message to an azure service bus endpoint (topic or queue).
Uses Shared Access Signatures in order to access service bus
https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-sas

### Usage

Import go-servicebus-publisher

```import ("github.com/nsterg/go-servicebus-publisher/publisher")```

or

```go get github.com/nsterg/go-servicebus-publisher/publisher```

From your service initialize a publisher with all the dependant services

```go
	httpClient := &http.Client{Timeout: time.Duration(1000) * time.Millisecond}

	// these values will be derived from your azure service bus configuration
	config := publisher.ServiceBusConfig{
		Namespace:           "my-namespace",
		Endpoint:            "my-queue",
		SharedKeyName:       "my-shared-keyname",
		SigningKey:          "my-signing-key",
		EndpointURL:         "https://my-name-space.servicebus.windows.net"
		SasTokenExpiresMS: 3000,
	}
	p := publisher.NewPublisher(httpClient, config)
```

 and publish a message

```go
	p.Publish(map[string]string{"I am here": "to poke you"})
```

### Integration Testing

In your service import test package

```testpublisher "github.com/nsterg/go-servicebus-publisher/test"```

Mock out the publisher by initializing an httptest.Server like this
```go
		publisherMockServer := testpublisher.MockPublisher(t, "test-queue", "7000", 201)
```

This will create an httptest.Server in the specified port and will mock out the call to the azure servicebus endpoint by returning a 201 status code response

Don't forget to always close the created mock server using
```go
		defer publisherMockServer.Close()
```