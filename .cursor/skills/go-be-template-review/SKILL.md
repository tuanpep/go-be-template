---
name: go-be-template-review
description: Reviews pull requests and local diffs against go-be-template architecture, OpenAPI, sqlc, and CI. Use when the user asks for a code review, PR review, or audit.
---

# Review

1. Read [code-review.md](../go-be-template/docs/code-review.md).
2. Label findings Critical, Major, or Minor.
3. Run `.cursor/skills/go-be-template/scripts/verify-review.sh`.

This skill has no `paths`. It loads when the user asks for a review, not because a file is open.
