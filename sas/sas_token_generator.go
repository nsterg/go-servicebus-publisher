package sas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"text/template"
	"time"
)

type SasGenerator struct {
	sasToken sasToken
	clock    clock
}

// NewSasGenerator configures a SasGenerator with a default clock
func NewSasGenerator(clock clock) SasGenerator {
	return SasGenerator{sasToken: sasToken{}, clock: clock}
}

type sasToken struct {
	expiry time.Time
	sas    string
}

// Generate generates a sas token using the provided params
// It checks if the existing token is valid using the expiry time
func (s SasGenerator) Generate(resourceURI string, signingKey string, expiresInMins int, policyName string) string {
	if s.isSasInvalid() {
		uri := template.URLQueryEscaper(resourceURI)

		durationFromNow := time.Now().Add(time.Duration(expiresInMins) * time.Minute)
		expiry := strconv.FormatInt(durationFromNow.Unix(), 10)
		signed := uri + "\n" + expiry

		val := computeHmac256(signed, signingKey)
		encodedVal := template.URLQueryEscaper(val)

		s.sasToken = sasToken{expiry: durationFromNow, sas: fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%s&skn=%s", uri, encodedVal, expiry, policyName)}
	}

	return s.sasToken.sas
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s SasGenerator) isSasInvalid() bool {
	return s.sasToken.expiry.Before(s.clock.Now())
}

type clock interface {
	Now() time.Time
}
