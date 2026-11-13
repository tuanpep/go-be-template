package migrations

import "embed"

// FS holds versioned SQL migrations for goose.
//
//go:embed *.sql
var FS embed.FS
