package httpx

import (
	"errors"
	"testing"
)

func TestParsePage(t *testing.T) {
	t.Parallel()
	limit, offset, err := ParsePage(1, 20)
	if err != nil || limit != 20 || offset != 0 {
		t.Fatalf("page 1: limit=%d offset=%d err=%v", limit, offset, err)
	}
	limit, offset, err = ParsePage(3, 10)
	if err != nil || limit != 10 || offset != 20 {
		t.Fatalf("page 3: limit=%d offset=%d err=%v", limit, offset, err)
	}
	if _, _, err = ParsePage(0, 500); !errors.Is(err, ErrInvalidPage) {
		t.Fatalf("out of range: %v", err)
	}
	if _, _, err = ParsePage(2147483647, 2); !errors.Is(err, ErrInvalidPage) {
		t.Fatalf("overflow: %v", err)
	}
}
