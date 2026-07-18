package db

import (
	"context"
	"testing"
)

func TestConnectWithContext_InvalidURL(t *testing.T) {
	_, err := ConnectWithContext(context.Background(), "invalid-url")
	if err == nil {
		t.Fatal("expected an error for an invalid database URL")
	}
}
