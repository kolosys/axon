package axon_test

import (
	"testing"
	"time"

	"github.com/kolosys/axon"
)

func TestOptions(t *testing.T) {
	u := axon.NewUpgrader(&axon.UpgradeOptions{
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
		MaxFrameSize:    8192,
		MaxMessageSize:  2097152,
		ReadDeadline:    10 * time.Second,
		WriteDeadline:   10 * time.Second,
		PingInterval:    30 * time.Second,
		PongTimeout:     5 * time.Second,
		Subprotocols:    []string{"chat", "json"},
		Compression:     true,
	})

	if u == nil {
		t.Fatal("upgrader should not be nil")
	}
}

func TestDefaultOptions(t *testing.T) {
	u := axon.NewUpgrader(nil)

	if u == nil {
		t.Fatal("upgrader should not be nil")
	}
}
