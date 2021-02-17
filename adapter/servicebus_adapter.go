package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const messageCreationFailure = "Failed to send message to service bus due to statusCode %d"

type ServiceBusAdapter struct {
	client HTTPClient
}

func NewServiceBusAdapter(client HTTPClient) ServiceBusAdapter {
	return ServiceBusAdapter{
		client: client,
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SendMessage sends a message to an event bus queue using a POST http request
// serviceNamespace is the namespace of the azure service bus
// endpoint is the name of the endpoint (topic or queue)
// message is the actual message
func (a ServiceBusAdapter) SendMessage(baseURL string, sasToken string, message interface{}) error {
	url := fmt.Sprintf("%s/messages", baseURL)
	requestByte, _ := json.Marshal(message)
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestByte))

	if err != nil {
		log.Printf("Failed to create http request. Error was %s", err.Error())
		return err
	}

	r.Header.Add("Authorization", sasToken)

	resp, err := a.client.Do(r)
	if err != nil {
		log.Printf("Failed to send http request. Error was %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		err = fmt.Errorf(messageCreationFailure, resp.StatusCode)
		log.Print(err.Error())
		return err
	}
	log.Printf("Successfully sent message to %s", url)

	return nil
}
