# Axon 🔌

![GoVersion](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)
![Zero Dependencies](https://img.shields.io/badge/Zero-Dependencies-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/kolosys/axon.svg)](https://pkg.go.dev/github.com/kolosys/axon)
[![Go Report Card](https://goreportcard.com/badge/github.com/kolosys/axon)](https://goreportcard.com/report/github.com/kolosys/axon)

A high-performance, zero-allocation WebSocket library for Go, designed for enterprise-grade applications with zero dependencies.

## Features

- **Zero-allocation hot paths** - Buffer pooling and efficient frame parsing
- **Context-aware** - Native `context.Context` support for cancellation and timeouts
- **Type-safe** - Generic `Conn[T]` for type-safe message handling
- **RFC 6455 compliant** - Full WebSocket protocol implementation
- **Zero dependencies** - Uses only the Go standard library
- **Built-in metrics** - Atomic counters for connection and message statistics
- **Thread-safe** - Safe for concurrent use

## Installation

```bash
go get github.com/kolosys/axon
```

## Quick Start

```go
package main

import (
    "net/http"
    "github.com/kolosys/axon"
)

type Message struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
}

func main() {
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn, err := axon.Upgrade[Message](w, r, nil)
        if err != nil {
            return
        }
        defer conn.Close(1000, "done")

        // Read messages
        msg, err := conn.Read(r.Context())
        if err != nil {
            return
        }

        // Write messages
        response := Message{ID: msg.ID, Text: "echo: " + msg.Text}
        conn.Write(r.Context(), response)
    })

    http.ListenAndServe(":8080", nil)
}
```

## Configuration

Configure the upgrader using `UpgradeOptions`:

```go
conn, err := axon.Upgrade[Message](w, r, &axon.UpgradeOptions{
    ReadBufferSize:  8192,
    WriteBufferSize: 8192,
    MaxFrameSize:    8192,
    MaxMessageSize:  2097152, // 2MB
    ReadDeadline:    30 * time.Second,
    WriteDeadline:   30 * time.Second,
    PingInterval:    30 * time.Second,
    PongTimeout:     5 * time.Second,
    Subprotocols:    []string{"chat", "json"},
    Compression:     true,
})
```

For default settings, pass `nil`:

```go
conn, err := axon.Upgrade[Message](w, r, nil)
```

## Performance

Axon is designed for high-performance scenarios:

- Zero allocations in frame parsing hot paths
- Buffer pooling for read/write operations
- Efficient masking and unmasking
- Atomic metrics collection

## Comparison

| Feature                   | kolosys/axon | gorilla/websocket | nhooyr.io/websocket | gobwas/ws |
| ------------------------- | ------------ | ----------------- | ------------------- | --------- |
| Zero-allocation hot paths | ✅           | ❌                | ❌                  | ✅        |
| Context-aware API         | ✅           | ❌                | ✅                  | ❌        |
| Type-safe generics        | ✅           | ❌                | ❌                  | ❌        |
| Dependencies              | 0            | 1                 | 1                   | 0         |
| Developer experience      | ⭐⭐⭐⭐⭐   | ⭐⭐⭐⭐          | ⭐⭐⭐⭐⭐          | ⭐⭐      |
