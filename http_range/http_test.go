package http_range

import (
	"io"
	"net/http"
	"testing"
)

func TestClient_GetWithByteRange(t *testing.T) {
	url := "https://jsonplaceholder.typicode.com/todos/1"
	client := Client{
		&http.Client{},
	}
	resp, err := client.GetWithByteRange(url, 5)
	if err != nil {
		t.Error(err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if len(bytes) != int(resp.ContentLength) {
		t.Errorf("Expected %d bytes, got %d", resp.ContentLength, len(bytes))
	}
}
