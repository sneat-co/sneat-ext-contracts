package contract4contactus

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	core "github.com/sneat-co/sneat-go-core"
)

const AuthorityEventsCollection = "authorityEvents"

// ContactAuthorityEvent is immutable native transition evidence, not an access
// grant or a membership directory. Products must reread current authority when
// processing it. Store in the SAME transaction as the new contact epochs.
// Delivery is durable: consumers acknowledge only after their own idempotent
// command commits, never before attempting delivery.
type ContactAuthorityEvent struct {
	V                    int       `json:"v" firestore:"v"`
	SpaceID              string    `json:"spaceID" firestore:"spaceID"`
	ContactID            string    `json:"contactID" firestore:"contactID"`
	UserID               string    `json:"userID,omitempty" firestore:"userID,omitempty"`
	MembershipGeneration string    `json:"membershipGeneration" firestore:"membershipGeneration"`
	BindingGeneration    string    `json:"bindingGeneration" firestore:"bindingGeneration"`
	Status               string    `json:"status" firestore:"status"`
	Member               bool      `json:"member" firestore:"member"`
	ObservedAt           time.Time `json:"observedAt" firestore:"observedAt"`
}

func (e ContactAuthorityEvent) Validate() error {
	if e.V != 1 || !core.IsAlphanumericOrUnderscore(e.SpaceID) || !core.IsAlphanumericOrUnderscore(e.ContactID) || e.ObservedAt.IsZero() {
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
		return nil
	default:
		return errors.New("invalid native contact authority status")
	}
}

// ID binds the complete native authority tuple. Observation time is excluded
// so a retry of the same transition has the same durable event identity.
func (e ContactAuthorityEvent) ID() string {
	raw, _ := json.Marshal([]any{e.V, e.SpaceID, e.ContactID, e.UserID, e.MembershipGeneration, e.BindingGeneration, e.Status, e.Member})
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
