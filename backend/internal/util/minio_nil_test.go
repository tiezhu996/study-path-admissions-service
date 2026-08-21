package util

import (
	"context"
	"strings"
	"testing"
)

func TestMinioNilReceiverNoPanicP1002(t *testing.T) {
	var m *MinIOClient
	_, err := m.Upload(context.Background(), "a.pdf", strings.NewReader("x"), 1, "application/pdf")
	if err == nil {
		t.Fatalf("expected error from nil MinIO client")
	}
}


func TestMinioGetNilReceiverNoPanicP1004(t *testing.T) {
	var m *MinIOClient
	_, _, _, err := m.Get(context.Background(), "a.pdf")
	if err == nil {
		t.Fatalf("expected error from nil MinIO client")
	}
}
