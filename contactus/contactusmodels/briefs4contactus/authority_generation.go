package briefs4contactus

import (
	"crypto/rand"
	"encoding/hex"
)

// RotateAuthorityGenerations fences inherited access across membership/binding
// lifecycle changes, including delete/recreate ABA. Persist BOTH returned fields
// in the same transaction as the authority mutation. Empty legacy epochs are
// deliberately not usable evidence for inherited access.
func (v *ContactBase) RotateAuthorityGenerations() error {
	var tokens [32]byte
	if _, err := rand.Read(tokens[:]); err != nil {
		return err
	}
	v.MembershipGeneration = hex.EncodeToString(tokens[:16])
	v.BindingGeneration = hex.EncodeToString(tokens[16:])
	return nil
}
