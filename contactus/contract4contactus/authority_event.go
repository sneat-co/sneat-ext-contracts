package contract4contactus

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	core "github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

const AuthorityEventsCollection = "authorityEvents"

// ContactAuthorityEvent is immutable native transition evidence, not an access
// grant or a membership directory. Products must reread current authority when
// processing it. Store in the SAME transaction as the new contact epochs.
// Delivery is durable: consumers acknowledge only after their own idempotent
// command commits, never before attempting delivery.
//
// TransitionKind is producer-declared cause, not authority. It never changes ID:
// IDs stay bound to the exact published (V, SpaceID, ContactID, UserID,
// MembershipGeneration, BindingGeneration, Status, Member) tuple, so existing
// point reads and retries keep working. Event writes must remain immutable
// transaction inserts; if the same immutable facts already have an event, a
// different TransitionKind (or ObservedAt) must not be used to overwrite that
// existing event. Consumers may use kind to REFUSE replayed/metadata events but
// must not grant project roles from it. See ContactAuthorityTransitionKind for
// the suppression and current-authority rules.
type ContactAuthorityEvent struct {
	V                    int                            `json:"v" firestore:"v"`
	SpaceID              string                         `json:"spaceID" firestore:"spaceID"`
	ContactID            string                         `json:"contactID" firestore:"contactID"`
	UserID               string                         `json:"userID,omitempty" firestore:"userID,omitempty"`
	MembershipGeneration string                         `json:"membershipGeneration" firestore:"membershipGeneration"`
	BindingGeneration    string                         `json:"bindingGeneration" firestore:"bindingGeneration"`
	Status               string                         `json:"status" firestore:"status"`
	Member               bool                           `json:"member" firestore:"member"`
	TransitionKind       ContactAuthorityTransitionKind `json:"transitionKind,omitempty" firestore:"transitionKind,omitempty"`
	ObservedAt           time.Time                      `json:"observedAt" firestore:"observedAt"`
}

func (e ContactAuthorityEvent) Validate() error {
	if e.V != 1 || coretypes.ValidateSpaceID(coretypes.SpaceID(e.SpaceID)) != nil || !core.IsAlphanumericOrUnderscore(e.ContactID) || e.ObservedAt.IsZero() {
		return errors.New("invalid native contact authority event")
	}
	for _, epoch := range []string{e.MembershipGeneration, e.BindingGeneration} {
		raw, err := hex.DecodeString(epoch)
		if err != nil || len(raw) != 16 || epoch != strings.ToLower(epoch) {
			return errors.New("invalid native authority epoch")
		}
	}
	switch e.Status {
	case "active", "archived", "deleted":
	default:
		return errors.New("invalid native contact authority status")
	}
	return e.TransitionKind.validate()
}

// ID binds the complete native authority tuple. Observation time and
// TransitionKind are excluded so a retry of the same transition has the same
// durable event identity regardless of declared cause or observation time.
func (e ContactAuthorityEvent) ID() string {
	raw, _ := json.Marshal([]any{e.V, e.SpaceID, e.ContactID, e.UserID, e.MembershipGeneration, e.BindingGeneration, e.Status, e.Member})
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
