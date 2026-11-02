package contract4contactus

import "errors"

// ContactAuthorityTransitionKind is neutral producer-declared provenance for a
// ContactAuthorityEvent. It records why the caller emitted an event, never what
// authority the subject has. It is not permission evidence, not an access grant
// and not a substitute for rereading current Contactus membership/authority.
//
// The empty value is the legacy V1/unknown kind: events written before this
// provenance field existed stay valid and must be treated as unknown, not as a
// native transition. ContactAuthorityTransitionNative is persisted as the
// exact string "authority_transition"; the public constant name stays neutral
// so nothing reads it as a role or permission decision.
//
// Consumers may use a non-empty kind to REFUSE a replayed/metadata event; they
// must not grant project roles from it. For a suppressed/blocked row, only a
// qualified ContactAuthorityTransitionNative event PLUS a membership-generation
// change relative to the CURRENT suppression anchor may request fresh
// admission. Legacy/unknown, metadata_repair, root and project changes must not
// lift suppression. Initial eligibility remains a separate full current
// authority/plan/cap check, and creation events do not indiscriminately
// override existing denials.
type ContactAuthorityTransitionKind string

const (
	// ContactAuthorityTransitionUnknown is the empty legacy V1 value. It means
	// the producer did not declare a cause; it is not proof of any transition.
	ContactAuthorityTransitionUnknown ContactAuthorityTransitionKind = ""
	// ContactAuthorityTransitionMetadataRepair is an epoch-only write emitted
	// while the scalar membership facts stayed unchanged. It must not be used
	// as evidence of a genuine leave/rejoin.
	ContactAuthorityTransitionMetadataRepair ContactAuthorityTransitionKind = "metadata_repair"
	// ContactAuthorityTransitionNative is a genuine server-owned authority
	// transition: a UID, status or member fact actually changed. It is persisted
	// as "authority_transition".
	ContactAuthorityTransitionNative ContactAuthorityTransitionKind = "authority_transition"
	// ContactAuthorityTransitionCreation is the first server-owned membership
	// fact snapshot, when there were no prior scalar facts.
	ContactAuthorityTransitionCreation ContactAuthorityTransitionKind = "creation"
)

func (k ContactAuthorityTransitionKind) validate() error {
	switch k {
	case ContactAuthorityTransitionUnknown,
		ContactAuthorityTransitionMetadataRepair,
		ContactAuthorityTransitionNative,
		ContactAuthorityTransitionCreation:
		return nil
	default:
		return errors.New("invalid contact authority transition kind")
	}
}

// ContactAuthorityState is the minimal scalar authority snapshot a producer
// captures before mutating membership. It deliberately carries no epoch,
// timestamp, MembershipChanged or inviteRepairsMembership-style boolean: those
// are write causes or observation metadata, not proof that membership facts
// changed.
type ContactAuthorityState struct {
	UserID string
	Status string
	Member bool
}

// ClassifyContactAuthorityTransition is the shared neutral classifier producers
// use so they do not duplicate provenance logic. A nil before snapshot means no
// prior membership facts existed, so the result is creation. Any actual UserID,
// Status or Member fact change is a native authority transition. Unchanged
// facts, when the caller is already emitting an event for a real server-owned
// epoch write, are a metadata repair.
//
// Callers must emit a ContactAuthorityEvent only after a real server-owned epoch
// write, must capture the before snapshot's scalar facts before mutation, and
// must not alias caller slices or records into it. The classifier never infers a
// native transition from an epoch difference or from any boolean such as
// MembershipChanged/inviteRepairsMembership (including index repair), and it
// does not treat observation time as a cause. It is not permission or evidence
// by itself and does not authorize anything.
func ClassifyContactAuthorityTransition(before *ContactAuthorityState, after ContactAuthorityState) ContactAuthorityTransitionKind {
	if before == nil {
		return ContactAuthorityTransitionCreation
	}
	if before.UserID != after.UserID || before.Status != after.Status || before.Member != after.Member {
		return ContactAuthorityTransitionNative
	}
	return ContactAuthorityTransitionMetadataRepair
}
