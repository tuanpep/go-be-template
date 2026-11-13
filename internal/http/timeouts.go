package httpserver

import "time"

// Timeouts align Chi request deadline with net/http Server fields (write must exceed Chi timeout).
const (
	RequestTimeout    = 55 * time.Second
	ReadHeaderTimeout = 5 * time.Second
	ReadTimeout       = 15 * time.Second
	WriteTimeout      = 60 * time.Second
	IdleTimeout       = 60 * time.Second
	ReadinessDrain    = 2 * time.Second
	ShutdownTimeout   = 10 * time.Second
)

// MaxRequestBodyBytes limits JSON request bodies on API routes.
const MaxRequestBodyBytes = 1 << 20 // 1 MiB
