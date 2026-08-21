package handler

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUploadNilStoreNoPanicP1001(t *testing.T) {
	h := NewUploadHandler(nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "a.pdf")
	if err != nil { t.Fatalf("create form file: %v", err) }
	if _, err := part.Write([]byte("hello")); err != nil { t.Fatalf("write part: %v", err) }
	if err := writer.Close(); err != nil { t.Fatalf("close writer: %v", err) }
	req := httptest.NewRequest(http.MethodPost, "/uploads", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	h.Upload(c)
}


func TestUploadGetNilStoreNoPanicP1003(t *testing.T) {
	h := NewUploadHandler(nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Params = gin.Params{{Key: "key", Value: "some-key"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/files/some-key", nil)
	h.Get(c)
}


type badStore struct{ marker int }

func (s *badStore) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_ = s.marker
	panic("store should not be called when nil")
}

func (s *badStore) Get(ctx context.Context, objectName string) (io.ReadCloser, int64, string, error) {
	_ = s.marker
	panic("store should not be called when nil")
}

func TestUploadNilStoreGuardPreventsCallP1005(t *testing.T) {
	h := &UploadHandler{store: (*badStore)(nil), logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	gin.SetMode(gin.TestMode)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "a.pdf")
	if err != nil { t.Fatalf("create form file: %v", err) }
	if _, err := part.Write([]byte("hello")); err != nil { t.Fatalf("write part: %v", err) }
	if err := writer.Close(); err != nil { t.Fatalf("close writer: %v", err) }
	req := httptest.NewRequest(http.MethodPost, "/uploads", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	h.Upload(c)
}
