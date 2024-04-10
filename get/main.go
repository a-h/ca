package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

var flagCert = flag.String("crt", "cert.pem", "certificate file")
var flagURL = flag.String("url", "https://localhost:8443", "URL to fetch")

func main() {
	flag.Parse()

	// Load x509 certs into pool.
	pool := x509.NewCertPool()

	// Read the certificate file.
	cert, err := os.ReadFile(*flagCert)
	if err != nil {
		log.Fatal("failed to read certificate file:", err)
	}
	pool.AppendCertsFromPEM(cert)

	// Add the certificate to the pool.
	// Create a new HTTP client with the custom transport.
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
		},
	}

	// Make the HTTP GET request.
	resp, err := client.Get(*flagURL)
	if err != nil {
		log.Fatal("failed to make GET request:", err)
	}

	// Print the response.
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatal("failed to read response body:", err)
	}
}
