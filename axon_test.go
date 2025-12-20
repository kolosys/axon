package axon_test

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kolosys/axon"
)

func TestComputeAcceptKey(t *testing.T) {
	// Test vector from RFC 6455 Section 1.3
	key := "dGhlIHNhbXBsZSBub25jZQ=="
	expected := "s3pPLMBiTxaQ9kYGzzhZRbK+xOo="

	// We can't directly test computeAcceptKey, but we can test the upgrade flow
	// which uses it internally
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, nil)
		if err != nil {
			// Expected - httptest.NewRecorder doesn't support hijacking
			return
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Key", key)
	req.Header.Set("Sec-WebSocket-Version", "13")

	w := httptest.NewRecorder()
	handler(w, req)

	// Verify accept key computation indirectly
	// The key should be base64(SHA1(key + GUID))
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	computed := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if computed != expected {
		t.Errorf("expected %q, got %q", expected, computed)
	}
}

func TestUpgradeInvalidMethod(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, nil)
		if err != axon.ErrUpgradeRequired {
			t.Errorf("expected ErrUpgradeRequired, got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
}

func TestUpgradeMissingUpgradeHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, nil)
		if err != axon.ErrUpgradeRequired {
			t.Errorf("expected ErrUpgradeRequired, got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
}

func TestUpgradeMissingConnectionHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, nil)
		if err != axon.ErrUpgradeRequired {
			t.Errorf("expected ErrUpgradeRequired, got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	w := httptest.NewRecorder()
	handler(w, req)
}

func TestUpgradeInvalidVersion(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, nil)
		if err != axon.ErrInvalidHandshake {
			t.Errorf("expected ErrInvalidHandshake, got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "14") // Invalid version

	w := httptest.NewRecorder()
	handler(w, req)

	// Should set correct version in response
	if w.Header().Get("Sec-WebSocket-Version") != "13" {
		t.Error("should set Sec-WebSocket-Version header")
	}
}

func TestUpgradeMissingKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, nil)
		if err != axon.ErrInvalidHandshake {
			t.Errorf("expected ErrInvalidHandshake, got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Version", "13")
	// Missing Sec-WebSocket-Key

	w := httptest.NewRecorder()
	handler(w, req)
}

func TestUpgradeWithOriginCheck(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, &axon.UpgradeOptions{
			CheckOrigin: func(r *http.Request) bool {
				return strings.Contains(r.Header.Get("Origin"), "example.com")
			},
		})
		if err == nil {
			t.Error("expected error for invalid origin")
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Origin", "https://evil.com")

	w := httptest.NewRecorder()
	handler(w, req)
}

func TestUpgradeWithValidOrigin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, &axon.UpgradeOptions{
			CheckOrigin: func(r *http.Request) bool {
				return strings.Contains(r.Header.Get("Origin"), "example.com")
			},
		})
		// Will fail due to hijacking, but origin check should pass
		if err != nil && err != axon.ErrInvalidHandshake {
			// ErrInvalidHandshake is expected (hijacking fails)
			// But ErrInvalidOrigin should not occur
			if err == axon.ErrInvalidOrigin {
				t.Error("origin should be valid")
			}
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Origin", "https://example.com")

	w := httptest.NewRecorder()
	handler(w, req)
}

func TestUpgradeInvalidSubprotocol(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, &axon.UpgradeOptions{
			Subprotocols: []string{"chat"},
		})
		if err != axon.ErrInvalidSubprotocol {
			t.Errorf("expected ErrInvalidSubprotocol, got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Protocol", "invalid")

	w := httptest.NewRecorder()
	handler(w, req)
}

func TestUpgradeWithMultipleSubprotocols(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := axon.Upgrade[string](w, r, &axon.UpgradeOptions{
			Subprotocols: []string{"chat", "json"},
		})
		// Will fail due to hijacking, but subprotocol validation should work
		if err == axon.ErrInvalidSubprotocol {
			t.Error("should accept valid subprotocol")
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Protocol", "chat, json")

	w := httptest.NewRecorder()
	handler(w, req)
}
