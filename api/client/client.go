package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const userAgent = "yandex-kassa-go-api/1.0"

// APIClient определяет транспортный уровень коммуникаций с API.
type APIClient struct {
	HTTP           *http.Client
	APIURL         string
	ShopID, Secret string
}

func (c *APIClient) get(ctx context.Context, uri string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.APIURL, uri), nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(c.ShopID, c.Secret)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("User-Agent", userAgent)

	return c.HTTP.Do(request.WithContext(ctx))
}

func (c *APIClient) post(ctx context.Context, uri, idempotencyKey string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.APIURL, uri), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Idempotence-Key", idempotencyKey)
	request.SetBasicAuth(c.ShopID, c.Secret)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("User-Agent", userAgent)

	return c.HTTP.Do(request.WithContext(ctx))
}

func (c *APIClient) Create(ctx context.Context, idempotencyKey string, body *[]byte) (io.ReadCloser, error) {
	response, err := c.post(ctx, "payments", idempotencyKey, *body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		return nil, c.errorWrap(response.Body)
	}
	return response.Body, nil
}

func (c *APIClient) Find(ctx context.Context, paymentID string) (io.ReadCloser, error) {
	response, err := c.get(ctx, fmt.Sprintf("payments/%s", paymentID))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		return nil, c.errorWrap(response.Body)
	}
	return response.Body, nil
}

func (c *APIClient) Cancel(ctx context.Context, idempotencyKey, paymentID string) (io.ReadCloser, error) {
	response, err := c.post(ctx, fmt.Sprintf("payments/%s/cancel", paymentID), idempotencyKey, []byte("{}"))
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (c *APIClient) Refund(ctx context.Context, idempotencyKey string, body *[]byte) (io.ReadCloser, error) {
	response, err := c.post(ctx, "refunds", idempotencyKey, *body)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (c *APIClient) Capture(ctx context.Context, idempotencyKey, paymentID string, body *[]byte) (io.ReadCloser, error) {
	response, err := c.post(ctx, fmt.Sprintf("payments/%s/capture", paymentID), idempotencyKey, *body)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}
