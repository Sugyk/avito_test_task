package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Service endpoints
	serviceAPIHost  = "http://localhost:8080"
	serviceUtilHost = "http://localhost:8081"

	// Timeouts
	defaultTimeout = 30 * time.Second
	healthTimeout  = 60 * time.Second
)

func connectTestDB(t *testing.T, ctx context.Context) (*sqlx.DB, error) {
	t.Helper()
	connStr := "postgres://postgres:postgres@postgres-db:5432/postgres?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), "pgx")
	return db, nil
}

func NewTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := connectTestDB(t, context.Background())
	require.NoError(t, err)

	require.NoError(t, db.Ping())

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

// WaitFor waits for a condition to be true or times out
func WaitFor(t *testing.T, timeout time.Duration, callback func() bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			assert.FailNow(t, "waiting for condition failed: timeout")
			return
		case <-ticker.C:
			if callback() {
				return
			}
		}
	}
}

// WaitForServiceReady waits for the service to be ready
func WaitForServiceReady(t *testing.T) {
	t.Helper()
	t.Log("Waiting for service to be ready...")

	WaitFor(t, healthTimeout, func() bool {
		resp, err := http.Get(serviceUtilHost + "/ready")
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	})

	t.Log("Service is ready")
}

// DoRequest performs an HTTP request and returns the response
func DoRequest(t *testing.T, method, path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err, "failed to marshal request body")
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	url := serviceAPIHost + path
	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(t, err, "failed to create request")

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to execute request")

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")
	resp.Body.Close()

	return resp, respBody
}

// DoGET performs a GET request
func DoGET(t *testing.T, path string, headers map[string]string) (*http.Response, []byte) {
	t.Helper()
	return DoRequest(t, http.MethodGet, path, nil, headers)
}

// DoPOST performs a POST request
func DoPOST(t *testing.T, path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	t.Helper()
	return DoRequest(t, http.MethodPost, path, body, headers)
}

// DoPUT performs a PUT request
func DoPUT(t *testing.T, path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
	t.Helper()
	return DoRequest(t, http.MethodPut, path, body, headers)
}

// DoDELETE performs a DELETE request
func DoDELETE(t *testing.T, path string, headers map[string]string) (*http.Response, []byte) {
	t.Helper()
	return DoRequest(t, http.MethodDelete, path, nil, headers)
}

// UnmarshalJSON unmarshals JSON response body
func UnmarshalJSON(t *testing.T, data []byte, v interface{}) {
	t.Helper()
	err := json.Unmarshal(data, v)
	require.NoError(t, err, "failed to unmarshal JSON: %s", string(data))
}

// AssertStatusCode asserts HTTP status code
func AssertStatusCode(t *testing.T, resp *http.Response, expectedStatus int) {
	t.Helper()
	assert.Equal(t, expectedStatus, resp.StatusCode,
		fmt.Sprintf("expected status %d, got %d", expectedStatus, resp.StatusCode))
}
