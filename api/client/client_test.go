package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupAPIClientTest() (*APIClient, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	client := &APIClient{
		HTTP:   server.Client(),
		APIURL: server.URL,
		ShopID: "your-shop-id",
		Secret: "your-secret",
	}
	return client, server
}

func teardownAPIClienTest(server *httptest.Server) {
	server.Close()
}

func TestAPIClient_get(t *testing.T) {
	client, server := setupAPIClientTest()
	defer teardownAPIClienTest(server)

	ctx := context.Background()
	uri := "test-endpoint"

	response, err := client.get(ctx, uri)

	if err != nil {
		t.Fatalf("Error calling get: %v", err)
	}

	if response.Request.Method != http.MethodGet {
		t.Errorf("Expected status method %s, but got %s", http.MethodGet, response.Request.Method)
	}

	if response.Request.Header.Get("User-Agent") != userAgent {
		t.Errorf("Expected User-Agent: %s, but got: %s", userAgent, response.Request.Header.Get("User-Agent"))
	}

	username, password, ok := response.Request.BasicAuth()
	if !ok || username != "your-shop-id" || password != "your-secret" {
		t.Error("Invalid or missing Basic Auth headers")
	}
}

func TestAPIClient_post(t *testing.T) {
	client, server := setupAPIClientTest()
	defer teardownAPIClienTest(server)

	ctx := context.Background()
	uri := "test-endpoint"
	idempotencyKey := "test-idempotency-key"
	body := []byte(`{"key": "value"}`)

	response, err := client.post(ctx, uri, idempotencyKey, body)
	if err != nil {
		t.Fatalf("Error calling post: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}

	if response.Request.Header.Get("User-Agent") != userAgent {
		t.Errorf("Expected User-Agent: %s, but got: %s", userAgent, response.Request.Header.Get("User-Agent"))
	}

	if response.Request.Header.Get("Idempotence-Key") != idempotencyKey {
		t.Errorf("Expected Idempotence-Key: %s, but got: %s", idempotencyKey, response.Request.Header.Get("Idempotence-Key"))
	}

	username, password, ok := response.Request.BasicAuth()
	if !ok || username != "your-shop-id" || password != "your-secret" {
		t.Error("Invalid or missing Basic Auth headers")
	}
}
