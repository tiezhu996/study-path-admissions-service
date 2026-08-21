package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/handler"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

type sa010BadStore struct{ marker int }

func (s *sa010BadStore) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_ = s.marker
	panic("store should not be called when nil")
}

func (s *sa010BadStore) Get(ctx context.Context, objectName string) (io.ReadCloser, int64, string, error) {
	_ = s.marker
	panic("store should not be called when nil")
}

func sa010UploadRequest(t *testing.T) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "a.pdf")
	if err != nil { t.Fatalf("create form file: %v", err) }
	if _, err := part.Write([]byte("hello")); err != nil { t.Fatalf("write part: %v", err) }
	if err := writer.Close(); err != nil { t.Fatalf("close writer: %v", err) }
	req := httptest.NewRequest(http.MethodPost, "/uploads", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestSA010UploadNil(t *testing.T) {
	h := handler.NewUploadHandler(nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = sa010UploadRequest(t)
	h.Upload(c)
}

func TestSA010UploadGuard(t *testing.T) {
	h := handler.NewUploadHandlerWithStore((*sa010BadStore)(nil), slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = sa010UploadRequest(t)
	h.Upload(c)
}

func TestSA010UploadGet(t *testing.T) {
	h := handler.NewUploadHandler(nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{{Key: "key", Value: "some-key"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/files/some-key", nil)
	h.Get(c)
}

func TestSA010MinioNil(t *testing.T) {
	var m *util.MinIOClient
	if _, err := m.Upload(context.Background(), "a.pdf", strings.NewReader("x"), 1, "application/pdf"); err == nil {
		t.Fatalf("expected error from nil MinIO client")
	}
}

func TestSA010MinioGet(t *testing.T) {
	var m *util.MinIOClient
	if _, _, _, err := m.Get(context.Background(), "a.pdf"); err == nil {
		t.Fatalf("expected error from nil MinIO client")
	}
}
