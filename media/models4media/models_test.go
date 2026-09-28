// Copyright 2026 Sneat.app

package models4media

import "testing"

func TestCropValidate(t *testing.T) {
	for _, tc := range []struct {
		crop    Crop
		wantErr bool
	}{
		{crop: Crop{CenterX: 0.5, CenterY: 0.5, Zoom: 1}, wantErr: false},
		{crop: Crop{CenterX: -0.1, CenterY: 0.5, Zoom: 1}, wantErr: true},
		{crop: Crop{CenterX: 1.1, CenterY: 0.5, Zoom: 1}, wantErr: true},
		{crop: Crop{CenterX: 0.5, CenterY: -0.1, Zoom: 1}, wantErr: true},
		{crop: Crop{CenterX: 0.5, CenterY: 1.1, Zoom: 1}, wantErr: true},
		{crop: Crop{CenterX: 0.5, CenterY: 0.5, Zoom: 0.9}, wantErr: true},
		{crop: Crop{CenterX: 0.5, CenterY: 0.5, Zoom: 8.1}, wantErr: true},
		{crop: Crop{CenterX: 0.5, CenterY: 0.5, Zoom: 8}, wantErr: false},
	} {
		err := tc.crop.Validate()
		if (err != nil) != tc.wantErr {
			t.Fatalf("Crop%+v.Validate() error = %v, wantErr %v", tc.crop, err, tc.wantErr)
		}
	}
}

func TestRefEqual(t *testing.T) {
	c1 := &Crop{CenterX: 0.5, CenterY: 0.5, Zoom: 1}
	c2 := &Crop{CenterX: 0.5, CenterY: 0.5, Zoom: 1}
	c3 := &Crop{CenterX: 0.6, CenterY: 0.5, Zoom: 1}

	var nilRef *Ref

	if !nilRef.Equal(nil) {
		t.Fatal("expected nilRef.Equal(nil) == true")
	}
	r1 := &Ref{MediaID: "m1"}
	if nilRef.Equal(r1) {
		t.Fatal("expected nilRef.Equal(r1) == false")
	}
	if r1.Equal(nilRef) {
		t.Fatal("expected r1.Equal(nilRef) == false")
	}
	r2 := &Ref{MediaID: "m2"}
	if r1.Equal(r2) {
		t.Fatal("expected r1.Equal(r2) == false for different media IDs")
	}
	r3 := &Ref{MediaID: "m1"}
	if !r1.Equal(r3) {
		t.Fatal("expected r1.Equal(r3) == true for same media ID with nil crops")
	}
	r4 := &Ref{MediaID: "m1", Crop: c1}
	if r1.Equal(r4) {
		t.Fatal("expected r1.Equal(r4) == false when one has crop and other does not")
	}
	if r4.Equal(r1) {
		t.Fatal("expected r4.Equal(r1) == false when one has crop and other does not")
	}
	r5 := &Ref{MediaID: "m1", Crop: c2}
	if !r4.Equal(r5) {
		t.Fatal("expected r4.Equal(r5) == true for equal crop values")
	}
	r6 := &Ref{MediaID: "m1", Crop: c3}
	if r4.Equal(r6) {
		t.Fatal("expected r4.Equal(r6) == false for different crop values")
	}
}

func TestRefValidate(t *testing.T) {
	for _, tc := range []struct {
		ref     Ref
		wantErr bool
	}{
		{ref: Ref{MediaID: ""}, wantErr: true},
		{ref: Ref{MediaID: "   "}, wantErr: true},
		{ref: Ref{MediaID: "m1"}, wantErr: false},
		{ref: Ref{MediaID: "m1", Crop: &Crop{CenterX: 0.5, CenterY: 0.5, Zoom: 1}}, wantErr: false},
		{ref: Ref{MediaID: "m1", Crop: &Crop{CenterX: 2.0, CenterY: 0.5, Zoom: 1}}, wantErr: true},
	} {
		err := tc.ref.Validate()
		if (err != nil) != tc.wantErr {
			t.Fatalf("Ref%+v.Validate() error = %v, wantErr %v", tc.ref, err, tc.wantErr)
		}
	}
}

func TestTargetValidate(t *testing.T) {
	for _, tc := range []struct {
		target  Target
		wantErr bool
	}{
		{target: Target{Type: "", ID: "id1", Scope: TargetScopeRoot}, wantErr: true},
		{target: Target{Type: "post", ID: "", Scope: TargetScopeRoot}, wantErr: true},
		{target: Target{Type: "post", ID: "id1", Scope: TargetScopeRoot, SpaceID: "s1"}, wantErr: true},
		{target: Target{Type: "post", ID: "id1", Scope: TargetScopeRoot, SpaceID: ""}, wantErr: false},
		{target: Target{Type: "post", ID: "id1", Scope: TargetScopeSpace, SpaceID: ""}, wantErr: true},
		{target: Target{Type: "post", ID: "id1", Scope: TargetScopeSpace, SpaceID: "   "}, wantErr: true},
		{target: Target{Type: "post", ID: "id1", Scope: TargetScopeSpace, SpaceID: "s1"}, wantErr: false},
		{target: Target{Type: "post", ID: "id1", Scope: "invalid_scope"}, wantErr: true},
	} {
		err := tc.target.Validate()
		if (err != nil) != tc.wantErr {
			t.Fatalf("Target%+v.Validate() error = %v, wantErr %v", tc.target, err, tc.wantErr)
		}
	}
}
