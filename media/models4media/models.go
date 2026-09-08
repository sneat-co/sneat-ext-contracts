// Copyright 2026 Sneat.app

package models4media

import (
	"fmt"
	"strings"

	"github.com/strongo/validation"
)

type Access string
type AssetStatus string
type LinkStatus string
type LinkRetention string
type TargetScope string

const (
	AccessPublic  Access = "public"
	AccessPrivate Access = "private"

	AssetStatusUploading AssetStatus = "uploading"
	AssetStatusReady     AssetStatus = "ready"
	AssetStatusOrphaned  AssetStatus = "orphaned"
	AssetStatusPurging   AssetStatus = "purging"
	AssetStatusDeleted   AssetStatus = "deleted"
	AssetStatusFailed    AssetStatus = "failed"

	LinkStatusActive  LinkStatus = "active"
	LinkStatusDeleted LinkStatus = "deleted"

	LinkRetentionFollowsSource LinkRetention = "follows_source"
	LinkRetentionRetained      LinkRetention = "retained"

	TargetScopeRoot  TargetScope = "root"
	TargetScopeSpace TargetScope = "space"
)

type Crop struct {
	CenterX float64 `json:"centerX" firestore:"centerX"`
	CenterY float64 `json:"centerY" firestore:"centerY"`
	Zoom    float64 `json:"zoom" firestore:"zoom"`
}

func (v Crop) Validate() error {
	if v.CenterX < 0 || v.CenterX > 1 || v.CenterY < 0 || v.CenterY > 1 {
		return validation.NewErrBadRecordFieldValue("center", "must be normalized to 0..1")
	}
	if v.Zoom < 1 || v.Zoom > 8 {
		return validation.NewErrBadRecordFieldValue("zoom", "must be between 1 and 8")
	}
	return nil
}

type Ref struct {
	MediaID string `json:"mediaID" firestore:"mediaID"`
	Crop    *Crop  `json:"crop,omitempty" firestore:"crop,omitempty"`
}

func (v *Ref) Equal(v2 *Ref) bool {
	if v == nil || v2 == nil {
		return v == v2
	}
	if v.MediaID != v2.MediaID || v.Crop == nil || v2.Crop == nil {
		return v.MediaID == v2.MediaID && v.Crop == v2.Crop
	}
	return *v.Crop == *v2.Crop
}

func (v Ref) Validate() error {
	if strings.TrimSpace(v.MediaID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("mediaID")
	}
	if v.Crop != nil {
		return v.Crop.Validate()
	}
	return nil
}

type Target struct {
	Scope    TargetScope `json:"scope" firestore:"scope"`
	SpaceID  string      `json:"spaceID,omitempty" firestore:"spaceID,omitempty"`
	Type     string      `json:"type" firestore:"type"`
	ID       string      `json:"id" firestore:"id"`
	ParentID string      `json:"parentID,omitempty" firestore:"parentID,omitempty"`
}

func (v Target) Validate() error {
	if strings.TrimSpace(v.Type) == "" || strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("target.type|target.id")
	}
	switch v.Scope {
	case TargetScopeRoot:
		if v.SpaceID != "" {
			return validation.NewErrBadRecordFieldValue("spaceID", "must be empty for root targets")
		}
	case TargetScopeSpace:
		if strings.TrimSpace(v.SpaceID) == "" {
			return validation.NewErrRecordIsMissingRequiredField("spaceID")
		}
	default:
		return fmt.Errorf("unknown media target scope %q", v.Scope)
	}
	return nil
}
