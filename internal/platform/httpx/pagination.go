package httpx

import (
	"errors"
	"math"
)

const (
	maxPage     = 1_000_000
	maxPageSize = 100
)

// ErrInvalidPage means page or page_size is outside the OpenAPI range, or the offset does not fit in int32.
var ErrInvalidPage = errors.New("invalid page")

// ParsePage converts 1-based page query params to sqlc limit/offset.
// Out-of-range values are rejected so a large page cannot wrap into a negative OFFSET.
func ParsePage(page, pageSize int) (limit, offset int32, err error) {
	if page < 1 || page > maxPage || pageSize < 1 || pageSize > maxPageSize {
		return 0, 0, ErrInvalidPage
	}
	offset64 := int64(page-1) * int64(pageSize)
	if offset64 > math.MaxInt32 {
		return 0, 0, ErrInvalidPage
	}
	return int32(pageSize), int32(offset64), nil
}
