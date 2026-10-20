package chessraiders_test

import (
	"testing"

	"github.com/sneat-co/sneat-ext-contracts/chessraiders"
)

func TestVersion(t *testing.T) {
	if chessraiders.Version == "" {
		t.Error("Version is empty")
	}
}
