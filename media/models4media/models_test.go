// Copyright 2026 Sneat.app

package models4media

import "testing"

func TestCropValidate(t *testing.T) {
	if err := (Crop{CenterX: .5, CenterY: .5, Zoom: 1}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Crop{CenterX: 2, CenterY: .5, Zoom: 1}).Validate(); err == nil {
		t.Fatal("expected normalized crop validation error")
	}
}
