package sas

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	layout = "2006-01-02T15:04:05.000Z"
	now    = "2020-09-13T11:12:13.000Z"
	future = "2020-09-13T11:12:14.000Z"
	past   = "2020-09-13T11:12:12.000Z"
)

func TestShouldCreateSasTokenWhenExistingIsEmpty(t *testing.T) {
	mockClock := new(FakeClock)
	mockClock.On("Now").Return(toTime(now))
	generator := NewSasGenerator(mockClock)
	token := generator.Generate("my-uri", "my-key", 123456, "my-policy")

	assert.NotNil(t, token)
	assert.True(t, strings.HasPrefix(token, "SharedAccessSignature "))
}

func TestShouldCreateSasTokenWhenExistingHasExpired(t *testing.T) {
	mockClock := new(FakeClock)
	mockClock.On("Now").Return(toTime(now))
	generator := NewSasGenerator(mockClock)
	generator.sasToken = sasToken{sas: "SharedAccessSignature existing-key", expiry: toTime(past)}

	token := generator.Generate("my-uri", "my-key", 123456, "my-policy")

	assert.NotNil(t, token)
	assert.True(t, strings.HasPrefix(token, "SharedAccessSignature "))
	assert.NotEqual(t, token, "SharedAccessSignature existing-key")
}

func TestShouldNotCreateSasTokenWhenExistingStillValid(t *testing.T) {
	mockClock := new(FakeClock)
	mockClock.On("Now").Return(toTime(now))
	generator := NewSasGenerator(mockClock)
	generator.sasToken = sasToken{sas: "SharedAccessSignature existing-key", expiry: toTime(future)}
	token := generator.Generate("my-uri", "my-key", 123456, "my-policy")

	assert.NotNil(t, token)
	assert.True(t, strings.HasPrefix(token, "SharedAccessSignature "))
	assert.Equal(t, token, "SharedAccessSignature existing-key")
}

func toTime(d string) time.Time {
	datetime, _ := time.Parse(layout, d)
	return datetime
}

type FakeClock struct {
	mock.Mock
}

func (c *FakeClock) Now() time.Time {
	args := c.Called()
	return args.Get(0).(time.Time)
}
