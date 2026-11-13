---
name: go-be-template-structure
description: Enforces feature folders and import direction in go-be-template. Use when creating a package, moving files, refactoring layers, or fixing an import cycle.
paths:
  - "internal/**/*.go"
  - "cmd/**/*.go"
---

# Structure

1. Read [structure.md](../go-be-template/docs/structure.md).
2. Put new behavior in `internal/<feature>/`. Keep `platform` free of feature imports.
3. Run `.cursor/skills/go-be-template/scripts/verify-structure.sh`.
