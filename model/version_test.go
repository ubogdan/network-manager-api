package model_test

import (
	"testing"

	"github.com/ubogdan/network-manager-api/model"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	if got, want := model.Version().String(), "0.1.1"; got != want {
		t.Errorf("Want version %s, got %s", want, got)
	}
}
